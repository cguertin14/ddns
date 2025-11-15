# ddns

DDNS utility written in Go for Cloudflare.

## Usage

This tool can be used as a container image: via `ghcr.io/cguertin14/ddns`.

## Configuration

The following environment variables **must** be set in order for the bot to work:

* `CLOUDFLARE_TOKEN`: your Cloudflare token with edit permissions on your zone;
* `RECORD_NAME`: the name of the A record you would like to update;
* `ZONE_NAME`: the name of the Cloudflare zone to edit.

The following environment variables are **optional**:

* `LOG_LEVEL`: the log level of the app, defaults to 'info'.