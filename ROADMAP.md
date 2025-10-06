# Roadmap

## v0.2.0 — DDD foundation + CLI UX
- Domain model: Connectivity check as core domain; aggregates for CheckPlan, ResultSet.
- Hexagonal architecture: ports/adapters for DNS, HTTP, ICMP, Proxy, Routing.
- Concurrency model: worker pool for checks; timeouts & cancellation (context).
- Output model: renderers (table/text/json) as adapters.
- Telemetry: structured logs, exit codes, error taxonomy.
- Packaging: smaller binaries via `-ldflags`, split adapters behind build tags.
- Compress: upx --best --lzma rsvpck.exe
