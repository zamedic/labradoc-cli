---
name: labradoc-cli
description: Use the Labradoc CLI with API token auth to call Labradoc endpoints (tasks, files, users, API keys, email, integrations, billing) and run raw API requests.
metadata: {"owner":"zamedic","repo":"labradoc-cli"}
---

# Labradoc CLI

Labradoc is an AI document intelligence platform that unifies emails, documents, and photos into one searchable system. It provides natural-language search and contextual answers over your own data, supports Gmail and Google Drive integrations, email forwarding, and manual uploads, and emphasizes GDPR-aligned hosting in Germany with strong privacy controls.

Use this skill to operate the `labradoc` CLI with API token authentication. It covers configuration and every available command.

## Install

```bash
go build -o labradoc .
```

```bash
go install ./...
```

## API Token Auth

The CLI sends the API token as the `X-API-Key` header.

Set the token using one of the following (highest wins):

```text
--api-token flag
API_TOKEN env var
labrador.yaml (api_token)
```

Optional base URL override:

```text
--api-url flag
API_URL env var
labrador.yaml (api_url)
```

Config file precedence:

```text
labrador.yaml
labrador.<ENVIRONMENT>.yaml
ENV vars (dots become underscores)
```

## Global Flags

```text
--api-url     API base URL (default https://labradoc.eu)
--api-token   API token (X-API-Key)
--timeout     HTTP timeout (default 30s)
```

## Raw Request

```bash
labradoc api request /api/tasks --method GET
labradoc api request /api/tasks --method POST --body '{"name":"Example"}'
labradoc api request /api/tasks --method POST --body-file ./payload.json
```

## Tasks

```bash
labradoc api tasks list
labradoc api tasks close --id <task-id>
labradoc api tasks close --ids <task-id> --ids <task-id>
```

## Files

```bash
labradoc api files list --status pending --status processed --page-size 50
labradoc api files upload --file ./document.pdf
labradoc api files get --id <file-id>
labradoc api files content --id <file-id> --out content.txt
labradoc api files ocr --id <file-id> --out ocr.txt
labradoc api files download --id <file-id> --out original.pdf
labradoc api files fields --id <file-id>
labradoc api files related --id <file-id>
labradoc api files reprocess --id <file-id>
labradoc api files tasks --id <file-id>
labradoc api files image --id <file-id> --page 1 --out page-1.png
labradoc api files preview --id <file-id> --page 1 --out page-1-preview.png
labradoc api files archive --id <file-id>
labradoc api files archive --ids <file-id> --ids <file-id>
labradoc api files question --id <file-id> --body '{"question":"What is the due date?"}'
labradoc api files search --body '{"question":"Find all invoices from Acme"}'
```

Note: `files search` returns a Server-Sent Events (SSE) stream.

## API Keys

```bash
labradoc api apikeys list
labradoc api apikeys create --name "CI token" --expires-at 2026-06-01T00:00:00Z
labradoc api apikeys revoke --id <key-id>
```

## User

```bash
labradoc api user credits
labradoc api user stats
labradoc api user language get
labradoc api user language set --language en
```

## Email

```bash
labradoc api email addresses
labradoc api email request --description "Inbound invoices"
labradoc api email list
labradoc api email body --id <email-id> --index 1 --out body.eml
```

## Google Integrations

```bash
labradoc api google drive status
labradoc api google drive token --scope "https://www.googleapis.com/auth/drive.readonly"
labradoc api google drive code --code <code>
labradoc api google drive refresh
labradoc api google drive revoke
labradoc api google gmail status
labradoc api google gmail token
labradoc api google gmail code --code <code>
labradoc api google gmail revoke
```

## Microsoft Integrations

```bash
labradoc api microsoft outlook token
labradoc api microsoft outlook code --code <code>
```

## Billing (Stripe)

```bash
labradoc api stripe checkout
labradoc api stripe pages-checkout
labradoc api stripe webhook --body-file ./stripe-event.json
```

## Troubleshooting

```text
Missing token: provide --api-token, API_TOKEN, or api_token in labrador.yaml
401/403: confirm API token and --api-url
```
