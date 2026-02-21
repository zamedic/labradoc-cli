# Labradoc CLI

Command-line client for Labradoc authentication and API operations (tasks and files).

## Requirements

- Go 1.26+

## Install

Build locally:

```bash
go build -o labradoc .
```

Or install:

```bash
go install ./...
```

## Configuration

The CLI reads configuration from the current working directory using Viper:

- Base config: `labrador.yaml`
- Environment override: `labrador.<env>.yaml` where `<env>` comes from `ENVIRONMENT` (default: `prod`)
- Environment variables override everything (dots become underscores)

Example `labrador.yaml` (only set values you want to override):

```yaml
api_url: https://labradoc.eu
api_token: your-api-key
keycloak:
  url: https://auth.labradoc.eu
  realm: labradoc
log:
  debug: false
```

Common environment variables (override config defaults only when needed):

- `API_URL`
- `API_TOKEN`
- `KEYCLOAK_URL`
- `KEYCLOAK_REALM`
- `LOG_DEBUG`
- `ENVIRONMENT`

Tokens and PKCE state are stored under the user config directory:

- `~/.config/labradoc/cli/token.json`
- `~/.config/labradoc/cli/pkce.json`

## Authentication

Default auth settings (only override when required):

- Auth URL: `https://auth.labradoc.eu`
- Realm: `labradoc`
- Client ID: `labradoc-openclaw`
- Scope: `openid profile email offline_access`

Login using a local callback:

```bash
labradoc auth login --api-url https://api.labradoc.eu
```

Generate a PKCE authorization URL (for manual flow):

```bash
labradoc auth url --redirect-uri http://127.0.0.1:18080/callback
```

Exchange a code for a token (uses saved PKCE state if present):

```bash
labradoc auth exchange --code <authorization-code>
```

Other auth commands:

```bash
labradoc auth token
labradoc auth refresh
labradoc auth status --api-url https://labradoc.eu
labradoc auth logout
```

API commands use API tokens by default. API token auth is the preferred method. OAuth is available if you prefer it — use `labradoc auth login` and pass `--use-auth-token` (or provide a bearer token explicitly).

## API Usage

The API commands accept either (defaults apply unless you override them):

- `--api-token` (sent as `X-API-Key`; default from `api_token`), or
- `--token` (Bearer token), or
- `--use-auth-token` to use the stored OAuth token from `labradoc auth login`

Raw request:

```bash
labradoc api request /api/tasks --method GET
labradoc api request /api/tasks --method POST --body '{"name":"Example"}'
```

Tasks:

```bash
labradoc api tasks list
labradoc api tasks close --id <task-id>
labradoc api tasks close --ids <task-id> --ids <task-id>
```

Files:

```bash
labradoc api files list --status pending --status processed --page-size 50
labradoc api files upload --file ./document.pdf
labradoc api files get --id <file-id>
labradoc api files content --id <file-id> --out content.txt
labradoc api files ocr --id <file-id>
labradoc api files download --id <file-id> --out original.pdf
labradoc api files question --id <file-id> --question "What is the due date?"
labradoc api files search --question "Find all invoices from Acme"
```

## Notes

- The binary name is `labradoc` (see `cmd/root.go`).
- `ENVIRONMENT=dev` loads `labrador.dev.yaml` in addition to `labrador.yaml`.
