# Contributing to imgn

## Architecture

```
cmd/imgn/          CLI entry point and command definitions (cobra)
internal/config/   Configuration loading (env → config file → flags)
internal/client/   Gemini API HTTP client
internal/models/   Model definitions and validation
internal/output/   File writing and JSON output formatting
```

## Code Style

- **Go 1.24+**, modules
- `gofmt` formatted (enforced in CI)
- No global mutable state
- Errors wrapped with context: `fmt.Errorf("operation: %w", err)`
- Structured logging via `log/slog`
- Context propagation for all I/O operations
- Table-driven tests

## Running Tests

```bash
go test ./...
```

Tests use `httptest` for API client testing — no real API calls needed.

## PR Guidelines

1. Run `go test ./...` and `go vet ./...` before submitting
2. Add tests for new features
3. Update docs if adding flags or commands
4. Keep commits focused and well-described
