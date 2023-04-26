package ddns

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/cguertin14/ddns/pkg/config"
	legacy "github.com/cguertin14/ddns/pkg/github"
	"github.com/cguertin14/logger"
	legacy_cf "github.com/cloudflare/cloudflare-go"
	"github.com/google/go-github/v52/github"
)

var (
	//go:embed report.tpl
	prReportTPL string
)

type RunReport struct {
	DnsChanged bool
	NewIP      string
}

type PRReport struct {
	NewIP, OldIP         string
	ZoneName, RecordName string
}

// returns wether or not DNS changed
func (c Client) Run(ctx context.Context, cfg config.Config) (RunReport, error) {
	// fetch zone ID from name
	zoneID, err := c.cloudflare.ZoneIDByName(cfg.ZoneName)
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch zone ID: %s", err)
	}
	identifier := legacy_cf.ZoneIdentifier(zoneID)

	// list dns records and find the one needed
	records, _, err := c.cloudflare.ListDNSRecords(ctx, identifier, legacy_cf.ListDNSRecordsParams{
		Name: fmt.Sprintf("%s.%s", cfg.RecordName, cfg.ZoneName),
	})
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch dns records: %s", err)
	}
	if len(records) == 0 {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to find dns record on cloudflare")
	}

	// fetch public IP
	newIP, err := getPublicIP()
	if err != nil {
		return RunReport{DnsChanged: false}, fmt.Errorf("failed to fetch public IP: %s", err)
	}
	logs := logger.NewFromContextOrDefault(ctx)
	logs.Infof("Current IP is %s", newIP)

	// check if IP changed and act if it did
	record := records[0]
	if record.Content != newIP {
		// step 1: update dns record
		if _, err := c.cloudflare.UpdateDNSRecord(ctx, identifier, legacy_cf.UpdateDNSRecordParams{
			Name:    cfg.RecordName,
			ID:      record.ID,
			Content: newIP,
			Type:    "A",
		}); err != nil {
			return RunReport{DnsChanged: false}, fmt.Errorf("failed to update dns record: %s", err)
		}

		// step 2: open PR on given repo
		oldIP := record.Content
		if cfg.UpdateGithubTerraform {
			if err := c.createPR(ctx, cfg, oldIP, newIP); err != nil {
				return RunReport{DnsChanged: true, NewIP: newIP}, err
			}
		}

		return RunReport{
			DnsChanged: true,
			NewIP:      newIP,
		}, nil
	}

	return RunReport{DnsChanged: false}, nil
}

func getPublicIP() (string, error) {
	res, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c Client) createPR(ctx context.Context, cfg config.Config, oldIP, newIP string) error {
	// start by fetching ref branch
	mainBranch, _, err := c.github.GetBranch(ctx, legacy.GetBranchRequest{
		Owner:      cfg.GithubRepoOwner,
		Repo:       cfg.GithubRepoName,
		BranchName: fmt.Sprintf("refs/heads/%s", cfg.GithubBaseBranch),
	})
	if err != nil {
		return fmt.Errorf("error when fetching branch %q: %s", cfg.GithubBaseBranch, err)
	}

	// create a new branch from ref branch
	branchName := fmt.Sprintf("minor/update-ip-%s", newIP)
	_, _, err = c.github.CreateBranch(ctx, legacy.CreateBranchRequest{
		Owner: cfg.GithubRepoOwner,
		Repo:  cfg.GithubRepoName,
		Reference: &github.Reference{
			Ref:    github.String(fmt.Sprintf("refs/heads/%s", branchName)),
			Object: mainBranch.Object,
		},
	})
	if err != nil {
		return fmt.Errorf("error when creating branch %q: %s", branchName, err)
	}

	// update file in new branch
	if err := c.updateFile(ctx, cfg, branchName, oldIP, newIP); err != nil {
		return fmt.Errorf("failed to update file %q: %s", cfg.GithubFilePath, err)
	}

	// create PR from new branch to ref branch.
	tpl, err := template.New("pr_report").Parse(prReportTPL)
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, PRReport{
		NewIP:      newIP,
		OldIP:      oldIP,
		ZoneName:   cfg.ZoneName,
		RecordName: cfg.RecordName,
	}); err != nil {
		return fmt.Errorf("failed to fill template: %s", err)
	}

	if _, _, err := c.github.CreatePullRequest(ctx, legacy.CreatePRRequest{
		Owner: cfg.GithubRepoOwner,
		Repo:  cfg.GithubRepoName,
		NewPullRequest: &github.NewPullRequest{
			Base: github.String(cfg.GithubBaseBranch),
			Head: github.String(branchName),
			Body: github.String(buffer.String()),
			Title: github.String(
				fmt.Sprintf("DNS update: public IP changed to %s", newIP),
			),
		},
	}); err != nil {
		return fmt.Errorf("error when opening pull request on repository: %s", err)
	}

	return nil
}

func (c Client) updateFile(ctx context.Context, cfg config.Config, branch, oldIP, newIP string) error {
	// fetch repository
	repoContent, _, _, err := c.github.GetRepositoryContents(ctx, legacy.GetRepositoryContentsRequest{
		Owner:  cfg.GithubRepoOwner,
		Repo:   cfg.GithubRepoName,
		Path:   cfg.GithubFilePath,
		Branch: branch,
	})
	if err != nil {
		return fmt.Errorf("error when fetching repo: %s", err)
	}

	// fetch current file content
	decoded, err := base64.StdEncoding.DecodeString(*repoContent.Content)
	if err != nil {
		return fmt.Errorf("error when decoding %q: %s", cfg.GithubFilePath, err)
	}
	fileContent := string(decoded)

	newFileContent := strings.ReplaceAll(fileContent, oldIP, newIP)
	_, _, err = c.github.UpdateFile(ctx, legacy.UpdateFileRequest{
		Owner:    cfg.GithubRepoOwner,
		Repo:     cfg.GithubRepoName,
		FilePath: cfg.GithubFilePath,
		RepositoryContentFileOptions: &github.RepositoryContentFileOptions{
			Content: []byte(newFileContent),
			Branch:  github.String(branch),
			Committer: &github.CommitAuthor{
				Name:  github.String("ddns-bot"),
				Email: github.String("ddns@cloudflare.com"),
				Date:  &github.Timestamp{Time: time.Now()},
			},
			Message: github.String(
				fmt.Sprintf("Updated public IP from %q to %q.", oldIP, newIP),
			),
			SHA: repoContent.SHA,
		},
	})
	if err != nil {
		return fmt.Errorf("error when updating file %q: %s", cfg.GithubFilePath, err)
	}

	// at this point, file is updated and commited to
	// corresponding repository.
	return nil
}
