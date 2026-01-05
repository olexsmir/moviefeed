# moviefeed - Agent Guidelines
## Build/Test/Lint Commands

```bash
# Build the application
go build -o moviefeed .

# Run the application
go run . -config config.yaml

# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Run tests with verbose output
go test -v ./...
```

## Code Style Guidelines
### Go Version
- Target Go 1.24 or later
- Use standard library features where possible

### Imports
- Group imports: standard library, third-party, internal (blank lines between)
- Use explicit import paths (no `.` imports)

### Formatting
- Use `go fmt` for all code formatting
- Maximum line length: 100 characters (soft limit)
- Use tabs for indentation (Go standard)

### Struct Tags
- Use json tags for JSON serialization: `json:"field_name"`
- Use yaml tags for YAML serialization: `yaml:"field_name"`
- Match tag names to external API format exactly (e.g., TMDB API snake_case)

### Error Handling
- Always return errors from functions that can fail
- Use `fmt.Errorf` with `%w` for error wrapping: `fmt.Errorf("failed: %w", err)`
- Validate inputs and configuration early
- Use `slog.Error` for failures with structured context
- Use `slog.Warn` for non-critical issues that shouldn't halt execution
- Continue processing on individual failures when appropriate (e.g., fetching episodes per show)

### Resource Management
- Always `defer` cleanup operations (file closing, response body closing)
- Check errors immediately after operations that return them
- Use `defer resp.Body.Close()` for HTTP responses

### Logging
- Use `log/slog` for all logging
- Log levels: Error for failures, Warn for non-critical issues, Info for normal operation
- Include relevant context using key-value pairs: `slog.Warn("msg", "key", value)`
- Avoid logging secrets or sensitive data

### Configuration
- Support both YAML and JSON config formats via file extension detection
- Provide sensible defaults for optional fields (e.g., port defaults to "8000")
- Validate required fields after loading config (API key, at least one show)
- Return clear error messages for validation failures

### HTTP/API
- Use net/http for HTTP client operations
- Set appropriate Content-Type headers
- Use gorilla/feeds for RSS feed generation
- Return proper HTTP status codes (500 for internal errors)
- Use defer for response body cleanup

### TMDB API Integration
- Support both IMDB IDs (tt*) and TMDB IDs directly
- Use find endpoint to convert IMDB to TMDB ID
- Fetch only first season and latest season (not all seasons)
- Filter episodes by air date (last 30 days)
- Handle missing or invalid air dates gracefully

### Testing
- Write tests for all public functions
- Use table-driven tests for multiple test cases
- Mock external API calls in tests
- Test error paths and edge cases
- Test config validation with both valid and invalid inputs

### Comments
- Exported functions should have godoc comments
- Keep comments brief and accurate
- Update comments when code changes
- Don't comment obvious code
