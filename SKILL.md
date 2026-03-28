---
name: labradoc-cli
description: Use the Labradoc CLI (labradoc) or its MCP server to authenticate and call Labradoc API endpoints — tasks, files, users, API keys, email, Google/Microsoft integrations, and Stripe billing — from OpenClaw or any MCP-compatible client. Trigger when the user wants to manage Labradoc documents, run natural-language searches, handle tasks, or interact with the Labradoc API via command line.
compatibility: "CLI: Go 1.21+ for building from source; prebuilt binaries available at github.com/zamedic/labradoc-cli/releases. MCP: any MCP-compatible client (Claude Desktop, OpenAI Agents, OpenClaw, Cursor, etc.)."
metadata:
  version: "1.0"
  author: Labradoc / Zamedic
---

# Labradoc CLI

Labradoc is an AI document intelligence platform that unifies emails, documents, and photos into one searchable system. It provides natural-language search and contextual answers over your own data, supports Gmail and Google Drive integrations, email forwarding, and manual uploads, with GDPR-aligned hosting in Germany.

This skill covers the `labradoc-cli` (binary: `labradoc`) and its MCP server. The CLI exposes two command groups: `auth` (OAuth PKCE) and `api` (all REST endpoints).

For full command reference see [commands reference](references/commands.md).

## When to Use

- Managing Labradoc tasks, files, and documents
- Uploading, searching, or extracting data from files
- Managing API keys, user preferences, or email addresses
- Setting up or revoking Google Drive, Gmail, or Microsoft Outlook OAuth
- Checking user credits or initiating Stripe billing flows
- Calling any Labradoc REST endpoint directly via `api request`

## Install

Prebuilt binaries: https://github.com/zamedic/labradoc-cli/releases

```bash
# macOS/Linux (homebrew)
brew install zamedic/tap/labradoc-cli

# Or download the binary and place on PATH
```

Or build from source (Go 1.21+):

```bash
git clone https://github.com/zamedic/labradoc-cli.git
cd labradoc-cli && go build -o labradoc .
```

## Configuration

Token precedence (highest wins):

```text
--api-token flag
LABRADOC_API_TOKEN env var
labrador.yaml  (api_token field)
```

Base URL override:

```text
--api-url flag
LABRADOC_API_URL env var
labrador.yaml  (api_url field)
```

Config file: `labrador.yaml` in the current working directory, or `labrador.<ENVIRONMENT>.yaml` for environment-specific overrides.

## Authentication

**API token (preferred):** Get your token at https://labradoc.eu/profile

```bash
export LABRADOC_API_TOKEN="your-token-here"
labradoc api tasks list
```

**OAuth PKCE (browser-based):**

```bash
labradoc auth login --api-url https://labradoc.eu
# Opens browser; stores token to ~/.config/labradoc/cli/token.json
labradoc api files list --use-auth-token
```

Other `auth` commands: `url`, `exchange`, `token`, `refresh`, `status`, `logout`.

## Common Operations

```bash
# List tasks
labradoc api tasks list

# Close a task
labradoc api tasks close --id <task-id>

# List files
labradoc api files list --status completed --page-size 50

# Upload a file
labradoc api files upload --file ./document.pdf

# Ask a question about a file
labradoc api files question --id <file-id> --question "What is the total amount?"

# Search files (SSE stream)
labradoc api files search --body '{"question":"Find all invoices from Acme"}'

# Get file content as text
labradoc api files content --id <file-id> --out content.txt

# Download original PDF
labradoc api files download --id <file-id> --out original.pdf

# Get page image
labradoc api files image --id <file-id> --page 1 --out page-1.png

# List API keys
labradoc api apikeys list

# Create API key
labradoc api apikeys create --name "CI token" --expires-at 2026-06-01T00:00:00Z

# Check AI credits
labradoc api user credits

# List email addresses
labradoc api email addresses

# Raw API request
labradoc api request /api/tasks --method GET
```

## Wrapper Script

```bash
./scripts/run-labradoc.sh api tasks list
```

The wrapper ensures `labradoc` is on PATH before forwarding arguments.

## Output Behaviour

| Condition | stdout | stderr | Exit code |
|-----------|--------|--------|-----------|
| Success | Response body (JSON or binary) | — | 0 |
| HTTP 4xx / 5xx | — | Response body (error JSON) | 1 |
| CLI error (missing flag, etc.) | — | Error message | 1 |

**Parsing rules for agents:**
- Parse stdout only on exit code 0.
- On exit code 1, read stderr for the error detail (may include the server's JSON error body).
- `files image`, `files preview`, `files download`, `files content`, `files ocr` return binary data; always use `--out <file>`.
- `files search` returns a **Server-Sent Events (SSE)** stream; parse it as SSE, not plain JSON.

## MCP Server

The same binary runs as an MCP server when invoked with MCP-specific arguments. Configure in your MCP client's config file:

```json
{
  "mcpServers": {
    "labradoc": {
      "command": "labradoc",
      "args": ["mcp"],
      "env": {
        "LABRADOC_API_TOKEN": "your-token-here"
      }
    }
  }
}
```

## Troubleshooting

```text
Missing token:  provide --api-token, LABRADOC_API_TOKEN, or api_token in labrador.yaml
401/403:         confirm API token and --api-url
Connection:       check network access to https://labradoc.eu
OAuth expires:    run labradoc auth refresh or labradoc auth login
```

## Integrations

See [references/integrations.md](references/integrations.md) for Google Drive, Gmail, and Microsoft Outlook OAuth commands.
