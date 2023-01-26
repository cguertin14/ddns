# Kubernetes Manifests

Here reside kubernetes manifests to use ddns in your cluster. 
Since `ddns` is not web service/microservice but rather a job in itself, it is defined here as  a cronjob that runs `hourly`.

## Tokens Secret

You may notice that inside the `ddns-cron` cronjob, it reads a secret called `tokens`, but we don't define it.
That's because you'll have to define it and create it in your cluster before creating the cronjob, like so (make sure to replace the actual token value by yours):

Make sure your tokens have the appropriate permissions:
* Github token: read & write to your repository of choice
* Cloudflare token: read & write DNS records in your chosen zone

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tokens
  namespace: ddns
type: Opaque
data:
  GITHUB_TOKEN: <YOUR_GITHUB_TOKEN_HERE>
  CLOUDFLARE_TOKEN: <YOUR_CLOUDFLARE_TOKEN_HERE>
```
