# TCP Port Scanner — Explanation Report

## Summary

This project is a Go command-line application that checks whether TCP ports on an authorized host accept connections. It accepts a target host, parses a flexible port expression, attempts TCP connections with a per-port timeout, and reports the result as open, closed, or error.

## Architecture

The project is separated into four layers:

- `cmd/scanner`: parses command-line flags, handles signals, maps failures to exit codes, and selects output mode.
- `internal/ports`: converts user input such as `22,80,8000-8010` into a validated, deduplicated, sorted list.
- `internal/scanner`: performs timeout-aware TCP checks. It includes sequential behavior through the same core check and a bounded worker pool for concurrency.
- `internal/output`: sorts and renders scan results as a terminal table or JSON.

This separation makes the network engine testable without invoking the CLI and makes future output formats independent of scanning logic.

## Important implementation decisions

### TCP connection method

The scanner uses `net.Dialer` with `DialContext`. A successful connection means the port accepted a TCP connection during the timeout window. A refused connection or timeout is reported as closed with the original diagnostic retained in the internal result. Context cancellation is reported as an error.

### Timeout handling

Every connection has its own timeout. This prevents a nonresponsive target from blocking the entire scan indefinitely. The command-line default is 500 milliseconds and can be changed with `-timeout`.

### Port ranges

The parser accepts individual ports, inclusive ranges, and comma-separated combinations. It rejects invalid boundaries, malformed values, reversed ranges, empty items, and ports outside 1 through 65535. A set removes duplicates and sorting provides deterministic output.

### Concurrency

The scanner uses a fixed worker pool rather than creating an unbounded goroutine for every port. This gives the user control through `-workers`, limits local resource usage, and makes the design safer for authorized lab testing. Results are sorted before output so concurrency does not change the displayed order.

### Testing

Unit tests cover the parser and output. Integration tests create local TCP listeners using `net.Listen`, so they do not depend on external websites. The race detector is part of the required quality checks because the concurrent scanner uses multiple goroutines and channels.

## Direct test procedure

### 1. Install Go

Install Go 1.22 or newer from the official Go website, then open a new PowerShell window and confirm:

```powershell
go version
```

### 2. Run automated tests

From the repository root:

```powershell
go test ./...
go vet ./...
go test -race ./...
go test -cover ./...
```

All commands should complete without failures.

### 3. Build the program

```powershell
go build -o scanner.exe .\cmd\scanner
```

### 4. Start a known local TCP service

For a simple Windows test, use any local service already running. Alternatively, run a temporary listener using a small helper program or test against a known local development server. Do not use a public target without authorization.

### 5. Run a human-readable scan

```powershell
.\scanner.exe -host 127.0.0.1 -ports 80,443,8000-8010 -timeout 300ms -workers 8
```

Expected behavior is a row per port, with `open` for accepting ports and `closed` for refused or timed-out ports.

### 6. Test JSON output

```powershell
.\scanner.exe -host localhost -ports 80,443 -json
```

The output should be valid JSON containing `host`, `port`, and `status` fields.

### 7. Test validation errors

```powershell
.\scanner.exe -host localhost -ports 0
.\scanner.exe -host localhost -ports 100-20
.\scanner.exe -host localhost -ports 80,,81
```

Each command should print an actionable error and return exit code 2.

## How to modify the project

- Change accepted port syntax in `internal/ports/parser.go` and add table-driven tests in `parser_test.go`.
- Change connection behavior in `internal/scanner/scanner.go`; preserve timeout and context handling.
- Change terminal or JSON formatting in `internal/output/output.go` and update output tests.
- Add CLI flags in `cmd/scanner/main.go`, then document them in `README.md`.
- Keep networking, parsing, and presentation separate; do not move all logic into `main.go`.

After each change:

```powershell
gofmt -w .
go vet ./...
go test -race ./...
go build -o scanner.exe .\cmd\scanner
```

## Advanced development roadmap

- Add a rate limiter to control requests per second.
- Add a context-aware streaming API for large scans.
- Add a `-verbose` mode for detailed timeout/refusal diagnostics.
- Add optional service-name lookup and carefully scoped banner detection.
- Add output files and CSV format without breaking stdout compatibility.
- Add configuration files after the command-line interface is stable.
- Add IPv4/IPv6 selection flags and dedicated IPv6 tests.
- Add benchmarks and profiling before optimizing.
- Add a release workflow that builds Windows, Linux, and macOS binaries.

Advanced work must retain the responsible-use restriction and must continue to pass the race detector and local integration tests.
