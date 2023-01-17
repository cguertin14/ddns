# ddns

DDNS tool written in Go for Cloudflare

## Usage

This tool can be used as a container image, simply use the image `ghcr.io/cguertin14/ddns`.

## Configuration

The following environment variables **must** be set in order for the bot to work:

* `CLOUDFLARE_TOKEN`: your Cloudflare token with edit permissions on your zone;
* `RECORD_NAME`: the name of the A record you would like to update;
* `ZONE_NAME`: the name of the Cloudflare zone to edit.

The following environment variables are **optional**:

* `LOG_LEVEL`: the log level of the app, defaults to 'info';
* `UPDATE_GITHUB_TERRAFORM`: either true or false, wether or not to update terraform file on Github;
* `GITHUB_TOKEN`: the Github token to use with write access to repository;
* `GITHUB_FILE_PATH`: the file path of the Terraform file to update in the Github repository;
* `GITHUB_REPO_OWNER`: the user who is the owner of the Github repository;
* `GITHUB_REPO_NAME`: the name of the Github repository to update;
* `GITHUB_BASE_BRANCH`: the name of the default branch to use for the Github repository.