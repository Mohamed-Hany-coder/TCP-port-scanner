# TCP Port Scanner in Go

A command-line TCP port scanner for authorized testing and system-programming education. It supports single ports, inclusive ranges, comma-separated expressions, timeouts, bounded concurrency, human-readable output, and JSON output.

## Responsible use

Scan only hosts and networks that you own or are explicitly authorized to test. This project intentionally does not implement stealth scanning, packet crafting, exploitation, UDP scanning, or security-control bypasses.

## Requirements

- Go 1.22 or newer
- Windows, Linux, or macOS

## Build

```powershell
go build -o scanner.exe .\cmd\scanner
```

## Usage

```powershell
.\scanner.exe -host 127.0.0.1 -ports 22,80,443,8000-8010
.\scanner.exe -host localhost -ports 8000-8010 -timeout 250ms -workers 16
.\scanner.exe -host localhost -ports 80,443 -json
```

Flags:

- `-host`: required hostname or IP address.
- `-ports`: required expression such as `80`, `20-22`, or `22,80,443,8000-8010`.
- `-timeout`: timeout per connection; default `500ms`.
- `-workers`: maximum concurrent workers; default `1`.
- `-json`: emit JSON instead of a terminal table.
- `-version`: print the version.

Exit codes: `0` completed, `2` invalid input, `3` fatal setup/output error, `4` cancelled.

## Test and quality checks

```powershell
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
go test -cover ./...
```

The integration tests use local TCP listeners and do not depend on public internet targets.

## Project structure

```text
cmd/scanner       CLI entry point
internal/ports    port expression parser and validation
internal/scanner  timeout-aware TCP checks and worker pool
internal/output   human-readable and JSON rendering
docs              project report and management material
```

## Direct modification workflow

1. Create a branch: `git switch -c feature/short-description`.
2. Modify the smallest relevant package.
3. Add or update tests beside the changed code.
4. Run `gofmt -w .`, `go vet ./...`, and `go test -race ./...`.
5. Build and manually run a localhost scan.
6. Commit with a Jira key, for example `TPS-123 add JSON output`.
7. Open a pull request and request review.

## Advanced extension roadmap

1. Add a service-name column using the standard `net` package's port-service lookup.
2. Add configurable output files while keeping stdout stable for scripts.
3. Add retry policy for transient network failures.
4. Add a streaming result channel for very large port lists.
5. Add benchmarks comparing sequential and worker-pool scans.
6. Add configuration-file support after the CLI contract is stable.
7. Add IPv6-specific integration tests.
8. Add a rate limiter if the scanner will be used in larger authorized lab environments.

Do not add advanced features until the core parser, timeout, bounded concurrency, and test suite remain stable.
