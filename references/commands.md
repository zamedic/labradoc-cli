# Labradoc CLI — Full Command Reference

All commands use the `labradoc` binary. Top-level groups: `auth`, `api`.

---

## Global Flags (both `auth` and `api`)

```text
--api-url     API base URL (default https://labradoc.eu)
--timeout     HTTP timeout (default 30s)
```

---

## `labradoc auth` — OAuth PKCE Authentication

Persistent flags:

```text
--auth-url     Auth server (default https://auth.labradoc.eu)
--realm        Keycloak realm (default labradoc)
--client-id    OAuth client ID (default labradoc-openclaw)
--api-url      API base URL (default https://labradoc.eu)
--scope        OAuth scope (default "openid profile email offline_access")
```

### `labradoc auth login`

Starts a local callback listener on `127.0.0.1` and prints an auth URL to open in the browser. Saves the resulting token to `~/.config/labradoc/cli/token.json`.

```bash
labradoc auth login --api-url https://labradoc.eu
labradoc auth login --timeout 2m --json
```

Flags: `--timeout` (default `2m`), `--json` (prints JSON to stdout; auth URL goes to stderr).

### `labradoc auth url`

Generates a PKCE authorization URL and saves PKCE state to `pkce.json`. Use this for manual/custom flows.

```bash
labradoc auth url --redirect-uri http://127.0.0.1:18080/callback --json
```

Flags: `--redirect-uri` (default `http://127.0.0.1:18080/callback`), `--json`.

### `labradoc auth exchange`

Exchanges an authorization code for a token and saves it.

```bash
labradoc auth exchange --code <authorization-code>
```

Flags: `--code` (required), `--code-verifier`, `--redirect-uri`, `--state` (uses saved `pkce.json` if missing), `--json`.

### `labradoc auth token`

Prints the stored OAuth access token.

```bash
labradoc auth token
labradoc auth token --json
```

### `labradoc auth refresh`

Refreshes the stored OAuth token.

```bash
labradoc auth refresh
labradoc auth refresh --json
```

### `labradoc auth status`

Validates the stored token against `GET /api/validate`.

```bash
labradoc auth status --api-url https://labradoc.eu
```

### `labradoc auth logout`

Deletes `~/.config/labradoc/cli/token.json`.

```bash
labradoc auth logout
```

---

## `labradoc api` — Labradoc REST API

Persistent flags:

```text
--api-url        API base URL (default https://labradoc.eu)
--api-token      API key (sent as X-API-Key; recommended over --token)
--token          Bearer token (overridden by --api-token)
--use-auth-token Use stored OAuth token from labradoc auth login
--timeout        HTTP timeout (default 30s)
```

---

### `labradoc api request [path]`

Make a raw API request to any endpoint.

```bash
labradoc api request /api/tasks --method GET
labradoc api request /api/tasks --method POST --body '{"name":"Example"}'
labradoc api request /api/tasks --method POST --body-file ./payload.json
labradoc api request /api/tasks --method POST --body-file -   # stdin
labradoc api request /api/user/files/search --accept text/event-stream --out results.sse
```

Flags: `--method` (default `GET`), `--body`, `--body-file` (`-` for stdin), `--content-type`, `--accept`, `--out`, `--no-auth`.

---

### `labradoc api tasks` — Task Operations

#### `labradoc api tasks list`

Lists all tasks. Endpoint: `GET /api/tasks`.

```bash
labradoc api tasks list
```

#### `labradoc api tasks close`

Close tasks by ID. Use `--id` for a single task, `--ids` for multiple.

```bash
# Single task: POST /api/tasks/<id>/close
labradoc api tasks close --id <task-id>

# Multiple: POST /api/tasks/close with {"id":[...]}
labradoc api tasks close --ids <task-id> --ids <task-id> --ids <task-id>
labradoc api tasks close --ids <task-id> --ids <task-id> --out response.json
```

Flags: `--id` (single), `--ids` (repeatable), `--out`.

---

### `labradoc api files` — File Operations

#### `labradoc api files list`

Lists user files. Endpoint: `GET /api/user/files`.

```bash
labradoc api files list
labradoc api files list --status completed --status extraction --page-size 50 --page-number 1
```

Valid `--status` values (repeatable):

```
New | multipart | googleDocument | Check_Duplicate | detectFileType |
htmlToPdf | preview | ocr | process_image | embedding | name_predictor |
document_type | extraction | task | completed | ignored | error |
not_supported | on_hold | duplicated
```

Flags: `--status` (repeatable), `--page-size`, `--page-number`.

#### `labradoc api files upload`

Uploads a file. Endpoint: `PUT /api/user/files` (multipart).

```bash
labradoc api files upload --file ./document.pdf
```

Flags: `--file` (required).

#### `labradoc api files get`

Gets file metadata. Endpoint: `GET /api/user/files/<id>`.

```bash
labradoc api files get --id <file-id>
```

Flags: `--id` (required).

#### `labradoc api files content`

Gets extracted text content. Endpoint: `GET /api/user/files/<id>/content`.

```bash
labradoc api files content --id <file-id> --out content.txt
```

Flags: `--id` (required), `--out`.

#### `labradoc api files ocr`

Gets OCR text. Endpoint: `GET /api/user/files/<id>/ocr`.

```bash
labradoc api files ocr --id <file-id> --out ocr.txt
```

Flags: `--id` (required), `--out`.

#### `labradoc api files download`

Downloads the original PDF. Endpoint: `GET /api/user/files/<id>/download`.

```bash
labradoc api files download --id <file-id> --out original.pdf
# Default output filename: <id>.pdf
```

Flags: `--id` (required), `--out`.

#### `labradoc api files question`

Ask a question about a file. Endpoint: `POST /api/user/files/<id>/question`.

```bash
labradoc api files question --id <file-id> --question "What is the due date?"
labradoc api files question --id <file-id> --body '{"question":"What is the total amount?"}'
printf '{"question":"Summarize the document"}' | labradoc api files question --id <file-id> --body-file -
```

Flags: `--id` (required), `--question` (sets JSON `question` field), `--body`, `--body-file` (`-` for stdin), `--out`.

#### `labradoc api files search`

AI-powered file search. Returns a **Server-Sent Events (SSE)** stream. Endpoint: `POST /api/user/files`.

```bash
labradoc api files search --question "Find all invoices from Acme Corp"
labradoc api files search --body '{"question":"Find contracts signed in 2025"}' --out results.sse
```

Parse the SSE stream for real-time results. Flags: `--question`, `--body`, `--body-file` (`-` for stdin), `--out`.

#### `labradoc api files fields`

Gets extracted fields. Endpoint: `GET /api/user/files/<id>/fields`.

```bash
labradoc api files fields --id <file-id> --out fields.json
```

Flags: `--id` (required), `--out`.

#### `labradoc api files related`

Gets related documents. Endpoint: `GET /api/user/files/<id>/related`.

```bash
labradoc api files related --id <file-id> --out related.json
```

Flags: `--id` (required), `--out`.

#### `labradoc api files reprocess`

Triggers reprocessing. Endpoint: `GET /api/user/files/<id>/reprocess`.

```bash
labradoc api files reprocess --id <file-id>
```

Flags: `--id` (required), `--out`.

#### `labradoc api files tasks`

Gets tasks for a document. Endpoint: `GET /api/user/files/<id>/tasks`.

```bash
labradoc api files tasks --id <file-id> --out tasks.json
```

Flags: `--id` (required), `--out`.

#### `labradoc api files image`

Gets a page image. Endpoint: `GET /api/user/files/<id>/image/<pageNumber>`.

```bash
labradoc api files image --id <file-id> --page 1 --out page-1.png
```

Flags: `--id` (required), `--page` (required, 1-indexed), `--out`.

#### `labradoc api files preview`

Gets a smaller preview image. Endpoint: `GET /api/user/files/<id>/image/preview/<pageNumber>`.

```bash
labradoc api files preview --id <file-id> --page 1 --out page-1-preview.png
```

Flags: `--id` (required), `--page` (required, 1-indexed), `--out`.

#### `labradoc api files archive`

Archives files (excludes them from retrieval). Endpoint: `POST /api/user/files/archive`.

```bash
labradoc api files archive --id <file-id>
labradoc api files archive --ids <file-id> --ids <file-id> --ids <file-id>
labradoc api files archive --ids <file-id> --ids <file-id> --out archive-response.json
```

Flags: `--id` (single), `--ids` (repeatable), `--out`.

---

### `labradoc api apikeys` — API Key Management

#### `labradoc api apikeys list`

Lists all API keys. Endpoint: `GET /api/user/apikeys`.

```bash
labradoc api apikeys list
```

#### `labradoc api apikeys create`

Creates a new API key. Endpoint: `POST /api/user/apikeys`.

```bash
labradoc api apikeys create --name "CI token"
labradoc api apikeys create --name "Production" --expires-at 2026-06-01T00:00:00Z
```

Flags: `--name` (required), `--expires-at` (RFC 3339 format, optional).

#### `labradoc api apikeys revoke`

Revokes an API key. Endpoint: `DELETE /api/user/apikeys/<keyId>`.

```bash
labradoc api apikeys revoke --id <key-id>
```

Flags: `--id` (required).

---

### `labradoc api user` — User Operations

#### `labradoc api user credits`

Gets AI credit balance. Endpoint: `GET /api/user/ai/credits`.

```bash
labradoc api user credits
```

#### `labradoc api user stats`

Gets user statistics. Endpoint: `GET /api/user/stats`.

```bash
labradoc api user stats
```

#### `labradoc api user language get`

Gets user language preference. Endpoint: `GET /api/user/preference/language`.

```bash
labradoc api user language get
```

#### `labradoc api user language set`

Sets user language preference. Endpoint: `POST /api/user/preference/language`.

```bash
labradoc api user language set --language en
labradoc api user language set --language de
```

Flags: `--language` (required).

---

### `labradoc api email` — Email Operations

#### `labradoc api email addresses`

Lists registered email addresses. Endpoint: `GET /api/emailAddresses`.

```bash
labradoc api email addresses
```

#### `labradoc api email request`

Requests a new inbound email address. Endpoint: `POST /api/emailAddress`.

```bash
labradoc api email request --description "Inbound invoices"
```

Flags: `--description` (optional).

#### `labradoc api email list`

Lists received emails. Endpoint: `GET /api/emails`.

```bash
labradoc api email list
```

#### `labradoc api email body`

Gets an email body by ID and index. Endpoint: `GET /api/email/<id>/<index>`.

```bash
labradoc api email body --id <email-id> --index 1 --out body.eml
```

Flags: `--id` (required), `--index` (required, 1-indexed), `--out`.

---

### `labradoc api google` — Google Integrations

#### Google Drive

```bash
# Check if Drive is connected
labradoc api google drive status

# Get OAuth URL for Drive scope
labradoc api google drive token --scope "https://www.googleapis.com/auth/drive.readonly"

# Exchange OAuth code
labradoc api google drive code --code <oauth-code>

# Trigger refresh/sync
labradoc api google drive refresh

# Revoke Drive access
labradoc api google drive revoke
```

Endpoints: `GET /api/google/drive`, `GET /api/google/drive/token?scope=...`, `GET /api/google/drive/code?code=...`, `GET /api/google/drive/refresh`, `DELETE /api/google/drive/token`.

#### Gmail

```bash
# Check if Gmail is connected
labradoc api google gmail status

# Get OAuth URL for Gmail
labradoc api google gmail token

# Exchange OAuth code
labradoc api google gmail code --code <oauth-code>

# Revoke Gmail access
labradoc api google gmail revoke
```

Endpoints: `GET /api/google/gmail`, `GET /api/google/gmail/token`, `GET /api/google/gmail/code?code=...`, `GET /api/google/gmail/revoke`.

---

### `labradoc api microsoft` — Microsoft Integration

#### Outlook

```bash
# Get OAuth URL for Outlook
labradoc api microsoft outlook token

# Exchange OAuth code
labradoc api microsoft outlook code --code <oauth-code>
```

Endpoints: `GET /api/microsoft/outlook/token`, `GET /api/microsoft/outlook/code?code=...`.

---

### `labradoc api stripe` — Stripe Billing

#### `labradoc api stripe checkout`

Creates a Stripe checkout session for AI credit purchase. Endpoint: `POST /api/stripe/checkout`.

```bash
labradoc api stripe checkout
```

#### `labradoc api stripe pages-checkout`

Creates a Stripe checkout session for the unlimited pages subscription. Endpoint: `POST /api/stripe/pages/checkout`.

```bash
labradoc api stripe pages-checkout
```

#### `labradoc api stripe webhook`

Receives a Stripe webhook event. Endpoint: `POST /api/stripe/webhook`. Does not require API authentication.

```bash
labradoc api stripe webhook --body-file ./stripe-event.json
labradoc api stripe webhook --body '{"type":"checkout.session.completed",...}'
cat stripe-event.json | labradoc api stripe webhook --body-file -
```

Flags: `--body`, `--body-file` (`-` for stdin).
