# Roadmap

## v0.2.0 — DDD foundation + CLI UX
- Domain model: 
- Hexagonal architecture: ports/adapters for DNS, HTTP, ICMP, Proxy, Routing.
- Output model: renderers (table/text/json) as adapters.
- makefile
- Concurrency model: worker pool for checks; timeouts & cancellation (context).
- Packaging: smaller binaries via `-ldflags`
- Compress: upx --best --lzma rsvpck.exe
