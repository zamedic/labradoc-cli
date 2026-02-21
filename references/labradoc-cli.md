# Labradoc CLI (AI Skill Reference)

Purpose: describe the Labradoc CLI in a format suitable for an AI Skill that needs to run auth and API operations.

## Binary

- Command: `labradoc`
- Top-level groups: `auth`, `api`

## Configuration And Precedence

Config is loaded from the current working directory:

1. Base config: `labrador.yaml` (optional)
2. Environment override: `labrador.<env>.yaml` where `<env>` is `ENVIRONMENT` (default `prod`)
3. Environment variables override everything (dots become underscores)

Common config keys and env vars:

- `api_url` -> `API_URL`
- `api_token` -> `API_TOKEN`
- `keycloak.url` -> `KEYCLOAK_URL`
- `keycloak.realm` -> `KEYCLOAK_REALM`
- `log.debug` -> `LOG_DEBUG`
- `ENVIRONMENT` selects the env-specific config file

Auth state files are stored under the OS user config directory in `labradoc/cli`:

- `token.json` (OAuth token)
- `pkce.json` (PKCE state)

Linux example path: `~/.config/labradoc/cli/`.

## Auth Model

API commands require one of:

- API token (preferred): sent as `X-API-Key` via `--api-token` or `api_token` config
- Bearer token: `--token`
- Stored OAuth access token: `--use-auth-token` (requires `labradoc auth login`)

`--api-token` overrides `--token` when both are set.

The `auth` commands implement OAuth PKCE via Keycloak. Defaults:

- Auth URL: `https://auth.labradoc.eu`
- Realm: `labradoc`
- Client ID: `labradoc-openclaw`
- Scope: `openid profile email offline_access`
- API URL default: `https://labradoc.eu`

## Command Reference

### `auth` (OAuth PKCE)

Persistent flags:

- `--auth-url` (default `https://auth.labradoc.eu`)
- `--realm` (default `labradoc`)
- `--client-id` (default `labradoc-openclaw`)
- `--api-url` (default `https://labradoc.eu`)
- `--scope` (default `openid profile email offline_access`)

Commands:

- `labradoc auth login`
  - Starts a local callback listener on `127.0.0.1` and prints an auth URL.
  - Saves the resulting token to `token.json`.
  - Flags: `--timeout` (default `2m`), `--json` (prints JSON to stdout; auth URL goes to stderr).
- `labradoc auth url`
  - Generates a PKCE authorization URL and saves PKCE state to `pkce.json`.
  - Flags: `--redirect-uri` (default `http://127.0.0.1:18080/callback`), `--json`.
- `labradoc auth exchange`
  - Exchanges an auth code for a token and saves it.
  - Required: `--code`.
  - Optional: `--code-verifier`, `--redirect-uri`, `--state` (uses saved `pkce.json` if missing).
  - Flag: `--json`.
- `labradoc auth token`
  - Prints the stored access token. Flag: `--json`.
- `labradoc auth refresh`
  - Refreshes the stored token. Flag: `--json`.
- `labradoc auth status`
  - Validates the stored token against `GET /api/validate`.
- `labradoc auth logout`
  - Deletes `token.json`.

### `api` (Labradoc API)

Persistent flags:

- `--api-url` (default `https://labradoc.eu`)
- `--api-token` (API key; `X-API-Key`)
- `--token` (Bearer token, ignored if `--api-token` is set)
- `--use-auth-token` (use stored OAuth token)
- `--timeout` (default `30s`)

Commands:

- `labradoc api request [path]`
  - Makes a raw API request.
  - Flags: `--method`, `--body`, `--body-file` (`-` for stdin), `--content-type`, `--accept`, `--out`, `--no-auth`.
  - If body is provided and method is not GET, Content-Type defaults to `application/json`.
  - Requires auth unless `--no-auth` is set.

- `labradoc api tasks list`
  - GET `/api/tasks`.

- `labradoc api tasks close`
  - Close tasks by ID.
  - Flags: `--id` (single), `--ids` (repeatable), `--out`.
  - If `--id` is set: POST `/api/tasks/<id>/close`.
  - Else: POST `/api/tasks/close` with `{"id":[...]}`.

- `labradoc api files list`
  - GET `/api/user/files` with optional query params.
  - Flags: `--status` (repeatable), `--page-size`, `--page-number`.

- `labradoc api files upload`
  - PUT `/api/user/files` with multipart form.
  - Flag: `--file`.

- `labradoc api files get`
  - GET `/api/user/files/<id>`.
  - Flag: `--id`.

- `labradoc api files content`
  - GET `/api/user/files/<id>/content`.
  - Flags: `--id`, `--out`.

- `labradoc api files ocr`
  - GET `/api/user/files/<id>/ocr`.
  - Flags: `--id`, `--out`.

- `labradoc api files download`
  - GET `/api/user/files/<id>/download`.
  - Flags: `--id`, `--out` (default output is `<id>.pdf`).

- `labradoc api files question`
  - POST `/api/user/files/<id>/question` with `text/plain` body.
  - Flags: `--id`, `--question`, `--body`, `--body-file` (`-` for stdin), `--out`.

- `labradoc api files search`
  - POST `/api/user/files` with `text/plain` body.
  - Flags: `--question`, `--body`, `--body-file` (`-` for stdin), `--out`.

## Examples

API token usage:

```bash
labradoc api tasks list --api-token "$API_TOKEN"
```

OAuth flow (local callback):

```bash
labradoc auth login --api-url https://labradoc.eu
labradoc api files list --use-auth-token
```

Manual PKCE flow:

```bash
labradoc auth url --redirect-uri http://127.0.0.1:18080/callback
labradoc auth exchange --code <authorization-code>
```

Raw request with JSON body:

```bash
labradoc api request /api/tasks --method POST --body '{"name":"Example"}'
```

Ask a question from stdin:

```bash
printf "What is the due date?" | labradoc api files question --id <file-id> --body-file -
```
