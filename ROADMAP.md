# Roadmap

## v0.2.0 — DDD foundation + CLI UX
- Domain model: 
- Hexagonal architecture: ports/adapters for DNS, HTTP, ICMP, Proxy
- Output model: renderers (table/text/json) as adapters.
- makefile
- Concurrency model: worker pool for checks; timeouts & cancellation (context).
- Collect CRM info
- Packaging: smaller binaries via `-ldflags`
- Compress: upx --best --lzma rsvpck.exe
