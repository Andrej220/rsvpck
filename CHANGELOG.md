# Changelog


All notable changes to this project will be documented in this file.

## [v0.2.0] — 2025-10-19

### Added
- **Renderer / output model** with pluggable **text** and **table** renderers.
- **Execution policy: Optimized** — skip remaining probes when early checks fail (per mode).
- **Custom error code module** for DNS/HTTP/ICMP.
- **Config loading** from external files and embedded FS (defaults included).
- **Render configuration** knobs (formatting helpers).
- **Spinner/animation** while gathering information.
- **Version flag & render flags** (force ASCII, force renderer).
- **Speed test** scaffold for future use.
- **TLS certificate fetch via proxy/VPN** when direct internet is not available.

### Changed
- **Endpoints grouped** into Direct / VPN / Proxy sets; analyzer adapted.
- **Default policy** switched to **Exhaustive**.
- **Config moved to YAML** (plus GE defaults).
- **Unicode detection & render refactor**; optional forced ASCII headers/characters.
- **Timeouts increased** (TCP/TLS dialer ~10s).
- **Module metadata** updated (`go.mod`, `go.sum`).

### Fixed
- **TLS certificate fetching** reliability and context timeout handling (plus waiting spinner for slow TLS).
- Corrected **IP addresses** in defaults.
- Adjusted **result interpretation** logic in analyzer.
- Fixed **CRM-number collection** command in host info.

### CI
- **Race detector** (CGO enabled) and CI fixes; **v0.2**-aligned workflow; temporary disable/enable phases during migration; UPX compression retained from earlier work.

### Notes
- Initial v0.2 entry existed with high-level bullets; this update consolidates **all v0.2 work** up to Oct 19, 2025 into a single release, as requested.

[Unreleased]: https://github.com/azargarov/rsvpck/compare/v0.2.0...HEAD


## [v0.2.0] - 2025-10-07
- Domain model
- Hexagonal architecture: ports/adapters for DNS, HTTP, ICMP, Proxy
- Output model: renderers table/text as adapters.
- Config file for probes