# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Sub-muse** is a terminal-based music streaming application (TUI) written in Go that interfaces with Subsonic-compatible servers (Subsonic, Airsonic, Funkwhale, etc.). It uses [Bubbletea](https://github.com/charmbracelet/bubbletea) for the TUI framework.

**Status:** Early development. Core Subsonic API client is functional with tests; UI components are being implemented.

## High-Level Architecture

### Package Structure

```
internal/
├── config/        - Configuration loading from environment variables
├── subsonic/      - Subsonic API client and models
│   ├── client.go      - HTTP client with interface-based DI for testability
│   ├── api.go         - API endpoint methods (GetArtists, GetAlbums, GetSongs, etc.)
│   ├── models.go      - XML data models (Song, Album, Artist, Child, etc.)
│   ├── mock.go        - Mock HTTP client for unit tests
│   ├── client_test.go - Client and request building tests
│   └── api_test.go    - API method tests
└── ui/            - TUI components (stubs)
    └── navigation.go   - Navigation state and tab management
```

### Subsonic Client Design

The `subsonic.Client` is designed for testability using interface-based dependency injection:

- **HTTPClient interface** (`client.go`): Allows mocking `http.Client.Do()` in tests without real network calls
- **NewClient()**: Returns a `Client` with a real `*http.Client` (30s timeout) for production use
- **buildRequest()**: Constructs Subsonic API URLs with required query parameters (username, password, client version)
- **sendRequest()**: Low-level HTTP call + XML decoding; used by all API methods
- **API methods** (`api.go`): GetArtists, GetAlbums, GetSongs, GetAlbum, GetSongsByArtist, GetCoverArt, GetMusicFolders

### Configuration

Configuration is loaded via environment variables (with defaults) in `internal/config/config.go`:
- `SUBSONIC_URL` (default: `http://localhost:4040`)
- `SUBSONIC_USERNAME` (required)
- `SUBSONIC_PASSWORD` (required)
- `SUBSONIC_CLIENT_NAME` (default: `music-tui`)

For development/testing, set these as env vars or create `.env` files (git-ignored).

## Common Commands

### Build
```bash
go build -o sub-muse .
```

### Run
```bash
./sub-muse
```

Dev mode (currently shows placeholder message):
```bash
go run . dev
```

### Testing
Run all tests with verbose output:
```bash
go test -v ./internal/...
```

Run tests in a specific package:
```bash
go test -v ./internal/subsonic/...
```

Run a single test:
```bash
go test -v ./internal/subsonic/... -run TestGetArtists_Success
```

With coverage (generates `coverage.out`):
```bash
go test -v ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out  # View in browser
```

### Dependencies
Install/update dependencies:
```bash
go mod tidy
```

Current dependencies:
- `github.com/stretchr/testify` - Testing assertions and mocking

## Testing

### Test Setup

Tests use `stretchr/testify` for assertions. The Subsonic client tests use a **mock HTTP client** to avoid real network calls:

1. **Mock HTTP Client** (`internal/subsonic/mock.go`):
   - Implements the `HTTPClient` interface
   - Accepts a function (`doFunc`) to customize response behavior per test

2. **Test Credentials**:
   - Loaded from environment variables with defaults (see `client_test.go`)
   - Defaults allow tests to run without configuration, but can be overridden for E2E tests

### Running Tests

All tests are unit tests (mocked, no network):
```bash
go test -v ./internal/subsonic/...
```

Tests cover:
- Client initialization (`TestNewClient_Success`)
- Request building (`TestClient_buildRequest_*`)
- Send request success/failure (`TestClient_sendRequest_*`)
- API methods: GetArtists, GetAlbums, GetSongs (see `api_test.go`)

Test coverage expectation: >80% for `client.go`, >70% for `api.go`.

## Development Workflow

### Adding a New API Method

1. Add the Subsonic response model to `internal/subsonic/models.go` if needed
2. Implement the method in `internal/subsonic/api.go`:
   - Create a response struct with XML tags
   - Call `c.sendRequest()` with endpoint name and params
   - Return parsed data or error
3. Add unit tests to `internal/subsonic/api_test.go`:
   - Mock the HTTP response with sample XML
   - Test success and failure cases
4. Test: `go test -v ./internal/subsonic/... -run TestYourMethod`

### Extending Configuration

1. Add a field to the `Config` struct in `internal/config/config.go`
2. Update `LoadConfig()` to read from env var (with default if applicable)
3. Add validation logic if required

### UI Development

The `internal/ui/navigation.go` file provides navigation state management with tab support (Songs, Artists, Albums, AlbumArtists). This is a stub; full TUI rendering is in progress. Once the Bubbletea integration is complete, it will handle:
- Tab switching
- Search input
- Focus management
- Message dispatching

## Environment Files

- `.env` and `.env.*` are git-ignored (see `.gitignore`)
- `.env.example` is tracked (add if creating new env-based config)
- For testing, set env vars directly or in `.env.test` (not committed)

## Common Patterns in This Codebase

- **Interface-based HTTP client**: Makes testing straightforward without external dependencies
- **Concrete XML models**: Each Subsonic response type (Artist, Album, Song, etc.) has a dedicated struct with XML tags
- **Testify assertions**: Use `require` for test assertions (fails fast) instead of `assert` (logs and continues)
- **No error wrapping in tests**: Mock functions return simple errors; tests check error presence/messages

## Key Files to Understand

- **main.go**: Entry point; checks dev mode flag, loads config, prints startup info
- **internal/subsonic/client.go**: The HTTP client and request building logic (most critical)
- **internal/subsonic/api.go**: All public API methods; study one method to understand the pattern
- **internal/config/config.go**: Environment-based configuration loader


<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:ca08a54f -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files

## Session Completion

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd dolt push
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
<!-- END BEADS INTEGRATION -->
