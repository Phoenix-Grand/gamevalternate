# Technology Stack

**Project:** GameVault Go Client
**Researched:** 2026-03-11
**Confidence Note:** External network tools (WebSearch, WebFetch, Context7) were unavailable during this research session. All findings are based on training data with knowledge cutoff August 2025. Confidence levels reflect this limitation. Version numbers should be verified against official sources before locking.

---

## Recommended Stack

### Core Framework

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Go | 1.22+ | Application language | Module workspace support, improved toolchain, CGo still required by Wails webview; 1.22 is the minimum for range-over-func and improved performance |
| Wails | v2 (v2.9.x) | Desktop GUI framework | PROJECT.md explicitly specifies v2. v3 was in active alpha/RC as of mid-2025 — not production-stable. v2 is the proven, documented, community-supported release |
| Svelte + SvelteKit (static) | Svelte 5.x | Frontend UI layer | Smallest compiled output of all Wails-compatible frontends (~50KB vs React ~130KB+). No virtual DOM — zero runtime overhead. Fastest TTI for a desktop app where startup feel matters. Wails scaffolding supports it natively |
| TypeScript | 5.x | Frontend type safety | Wails generates TypeScript bindings for Go functions automatically; TS is required to use them correctly |
| Vite | 5.x | Frontend build tool | Wails v2 uses Vite as the dev server and bundler. This is the default; do not swap it out |

**Confidence:** MEDIUM — Wails v2 + Svelte is well-established community pattern. Version numbers need verification against current releases.

### Wails v3 Assessment

**Do not use Wails v3 for this project.**

As of August 2025, Wails v3 was in late alpha/RC. Key reasons to avoid:
- Breaking API changes between RC releases were still occurring
- Multi-window support (v3's headline feature) is not needed by this project
- Community examples, Stack Overflow answers, and tutorials overwhelmingly cover v2
- The `wails build` cross-compilation ecosystem (Docker images, community CI scripts) is built around v2
- v3 drops the Vite dev server integration in favor of a different approach — migrations are non-trivial

**Revisit when:** v3 has a stable release tag and migration guide is finalized (likely 2026).

### GameVault API Client

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| oapi-codegen | v2.x | OpenAPI → Go client codegen | GameVault backend (NestJS) exposes a Swagger/OpenAPI spec at `/api` or `/api-json`. oapi-codegen generates idiomatic Go interfaces + net/http-based client. Avoids 3000-line manual client maintenance |
| net/http (stdlib) | stdlib | HTTP transport | oapi-codegen targets stdlib by default. Do not add a third-party HTTP client — no dependency needed |

**Confidence:** MEDIUM — GameVault backend being NestJS makes OpenAPI spec availability very likely (NestJS generates it by default). oapi-codegen v2 is the current maintained fork. Verify spec URL by checking `https://github.com/Phalcode/gamevault-backend`.

**If OpenAPI spec is unavailable or incomplete:** Fall back to a handwritten client using stdlib `net/http`. Structure it as an interface so it can be swapped. Do not use `go-resty` or `gentleman` — they add dependencies without meaningful benefit for a typed API client.

### State Management (Frontend)

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Svelte stores (built-in) | Svelte 5 runes | Reactive state | Svelte 5's rune-based reactivity (`$state`, `$derived`) replaces the old store API. No external state library needed — Svelte's built-in primitives handle the complexity of this app |
| Tauri-style event bridge | Wails built-in | Go ↔ frontend events | Wails provides `runtime.EventsEmit` (Go side) and `Events.On` (JS side) for async notifications. Use for download progress, cloud save events |

**Confidence:** HIGH — Svelte store patterns are framework internals, not external dependency.

### Persistent Storage

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| SQLite via modernc/sqlite | v1.x | Local app database | Stores server profiles, user accounts, game metadata cache, download state. `modernc.org/sqlite` is the pure-Go SQLite driver — no CGo required, works in cross-compilation targets without host toolchain matching |
| GORM | v2.x | ORM | Pragmatic choice for this app's complexity level. Schema migrations, associations for profiles/servers. Alternative: `sqlc` if you want compile-time SQL safety, but adds codegen step |

**Confidence:** MEDIUM — `modernc/sqlite` vs `mattn/go-sqlite3` is the critical choice. `modernc` is pure Go which is essential for Docker cross-compilation without CGo complications. `mattn/go-sqlite3` requires CGo and a C compiler for each target — this breaks the cross-compile pipeline.

### Docker Build Pipeline

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `crazymax/xgo` | latest | CGo cross-compilation base | Provides pre-built cross-compiler toolchains for linux/amd64, linux/arm64, windows/amd64. Handles the MinGW toolchain for Windows targets. Alternative: `goreleaser/goreleaser-cross` but it requires goreleaser config |
| Wails CLI | v2.9.x | Build orchestration | `wails build -platform linux/amd64` etc. Run inside the cross-compile container |
| Docker BuildKit | 1.x | Multi-stage builds | Use `--platform` flag for QEMU-based arm64 builds when native cross-compile isn't sufficient |
| osxcross | via docker image | macOS cross-compilation | macOS targets require osxcross + macOS SDK (legal: must use SDK from your own Mac). This is the hard problem — see PITFALLS.md |

**Confidence:** MEDIUM — CGo cross-compilation for Wails is genuinely hard. The specific Docker image names/tags need verification. The macOS SDK legal constraint is a real blocker that must be addressed in planning.

**Critical constraint:** Wails uses webview2 (Windows) and WebKit (Linux/macOS), both of which require CGo bindings. You cannot avoid CGo in this project. This means every Docker build target needs the appropriate C cross-compiler for that target platform.

### VNC/noVNC Runtime Container

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| Xvfb | OS package | Virtual framebuffer | Creates a virtual display `:99` inside the container. Required by any X11 GUI app in headless Docker |
| x11vnc | OS package | VNC server | Reads from Xvfb display and serves it over VNC protocol. Lighter than TigerVNC for this use case (single app, no desktop environment) |
| noVNC | 1.4.x | Browser-based VNC client | Converts VNC to WebSocket, serves HTML5 canvas UI. Users access `http://host:6080/vnc.html` — no client software needed |
| websockify | 0.11.x | VNC ↔ WebSocket bridge | noVNC requires WebSocket; x11vnc speaks raw VNC. websockify is the standard bridge. Ships with noVNC |
| supervisord | via python3-supervisor | Process manager | Manages Xvfb + x11vnc + websockify + the Go app as a single Docker entrypoint. Standard pattern for multi-process Docker containers |

**Base image recommendation:** `debian:bookworm-slim` — good package availability, smaller than ubuntu, stable. Install: `xvfb x11vnc python3-websockify novnc supervisor`.

**Do not use:** TigerVNC server — overkill for a single-app container, and requires more setup. KasmVNC — adds unnecessary complexity. VirtualGL — for GPU-accelerated apps only (not needed here).

**Confidence:** MEDIUM — This is a well-established pattern for containerized GUI apps. x11vnc + noVNC + Xvfb + supervisord is the standard Linux stack as of 2025. noVNC version needs verification.

### GitHub Actions CI/CD

| Technology | Version | Purpose | Why |
|------------|---------|---------|-----|
| `actions/checkout` | v4 | Source checkout | Current stable major |
| `actions/setup-go` | v5 | Go toolchain | Current stable major |
| `docker/setup-buildx-action` | v3 | Docker BuildKit | Required for multi-platform builds |
| `docker/build-push-action` | v5 | Build + push Docker image | Standard for Docker Hub / GHCR publishing |
| goreleaser | v2.x | Release artifact management | Generates changelogs, creates GitHub releases, uploads binaries. Works alongside Wails build |
| `softprops/action-gh-release` | v2 | GitHub release creation | Upload Wails-built binaries to GitHub Releases if not using goreleaser |

**Build matrix strategy:** Use separate jobs per target (not matrix) for cross-compilation — cross-compile jobs take different times and have different failure modes. Linux native builds on `ubuntu-latest`; Windows builds via cross-compilation from Linux; macOS builds on `macos-latest` runners (required for macOS SDK legality).

**Confidence:** MEDIUM — Action versions need verification. The macOS-on-macOS-runner constraint is a firm legal/practical requirement, not a preference.

---

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Wails version | v2 (stable) | v3 (alpha/RC) | v3 API unstable as of Aug 2025; no production case studies; cross-compile tooling not mature |
| Frontend | Svelte 5 | React 18 | React bundle is ~3x larger; heavier runtime; no meaningful advantage for a desktop app where bundle size affects startup |
| Frontend | Svelte 5 | Vue 3 | Vue is a valid choice but Svelte has stronger Wails community adoption and smaller output |
| SQLite driver | modernc/sqlite (pure Go) | mattn/go-sqlite3 (CGo) | mattn requires CGo — breaks Docker cross-compilation pipeline; each target needs matching C toolchain |
| API client | oapi-codegen | go-swagger | go-swagger is older, heavier, generates more boilerplate; oapi-codegen v2 is the current community standard |
| API client | oapi-codegen | manual net/http | Manual is fine for small APIs but GameVault has 50+ endpoints — maintenance burden too high |
| VNC server | x11vnc | TigerVNC | TigerVNC designed for full desktop sessions; x11vnc simpler for single-app use; less config |
| Container process mgmt | supervisord | s6-overlay | s6-overlay is excellent but adds learning curve; supervisord is universal knowledge |
| Build pipeline | Docker cross-compile | Native runners per platform | GitHub Actions macOS minutes cost 10x Linux; cross-compile where possible, use macOS runners only for macOS targets |

---

## Installation

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Scaffold new project with Svelte + TypeScript template
wails init -n gamevault-client -t svelte-ts

# Go dependencies (add to go.mod)
go get github.com/wailsapp/wails/v2
go get modernc.org/sqlite
go get gorm.io/gorm
go get gorm.io/driver/sqlite        # GORM SQLite driver (modernc adapter)
go get github.com/oapi-codegen/oapi-codegen/v2  # Install as tool, not dep

# Frontend dependencies (from frontend/ directory)
npm install                          # Vite + Svelte already in template
npm install -D typescript @tsconfig/svelte svelte-check
```

---

## Key Version Pins

These should be verified against current releases at project start:

| Package | Expected Current | Verify At |
|---------|-----------------|-----------|
| `github.com/wailsapp/wails/v2` | v2.9.x | https://github.com/wailsapp/wails/releases |
| `svelte` | 5.x | https://github.com/sveltejs/svelte/releases |
| `modernc.org/sqlite` | v1.x | https://pkg.go.dev/modernc.org/sqlite |
| `gorm.io/gorm` | v1.25.x | https://github.com/go-gorm/gorm/releases |
| `oapi-codegen` | v2.x | https://github.com/oapi-codegen/oapi-codegen/releases |
| noVNC | 1.4.x | https://github.com/novnc/noVNC/releases |

---

## Sources

- Training data (knowledge cutoff August 2025) — all findings LOW to MEDIUM confidence
- PROJECT.md constraints (Wails v2, no CGo where avoidable) — HIGH confidence (project-defined)
- Wails v2 documentation pattern (Svelte/React/Vue support) — MEDIUM confidence
- modernc/sqlite pure-Go rationale — HIGH confidence (well-documented CGo avoidance pattern)
- x11vnc + noVNC + Xvfb supervisord pattern — MEDIUM confidence (standard headless GUI Docker pattern)
- **NOTE:** Context7, WebSearch, WebFetch, and Bash tools were all unavailable during this research session. All version numbers MUST be verified against official sources before use.
