# Project Research Summary

**Project:** GameVault Go Client
**Domain:** Cross-platform desktop game library client (Go/Wails v2, Docker/VNC distribution)
**Researched:** 2026-03-11
**Confidence:** MEDIUM

## Executive Summary

GameVault Go Client is a native desktop replacement for the official Electron-based gamevault-app. It connects to self-hosted GameVault backends, enables game discovery and download, and supports cloud saves and progress tracking. The project's core value proposition is a smaller, more resource-efficient binary (50–150MB vs Electron's 300MB+) with first-class Linux support and a Docker/VNC/noVNC distribution path that allows the client to run headlessly on a NAS or home server and be accessed from any browser. The recommended build is Go 1.22+ with Wails v2 (explicitly not v3, which remained alpha as of August 2025), Svelte 5 for the frontend, SQLite via the pure-Go modernc driver for local state, and oapi-codegen v2 to generate a typed HTTP client from the GameVault OpenAPI spec.

The most important architectural decision is that the Go backend owns all business logic and API communication. The Svelte frontend is a rendering shell only — it calls bound Go methods and subscribes to Wails events; it never makes direct HTTP requests to the GameVault server. This discipline enables secure credential handling (OS keyring), download resume, and cloud save sync to be implemented entirely in Go without any frontend complexity leaking in. The layered service structure (API Client, Download Manager, Cloud Save Manager, Profile Store, EventBus) enables incremental delivery: each service can be built and tested independently before wiring into the App struct.

The project carries three hard build constraints that must be accepted from day one, not discovered late: (1) Wails requires CGo on all platforms — standard Go cross-compilation does not work; (2) macOS targets cannot be built from Linux due to Apple SDK redistribution restrictions — macOS-native GitHub Actions runners are mandatory; (3) the Docker runtime image cannot use Alpine Linux because WebKit2GTK requires glibc. These constraints force a specific CI/CD matrix and Docker base image. Getting these right in Phase 1 prevents costly pipeline rebuilds later.

## Key Findings

### Recommended Stack

The stack is anchored on Wails v2 + Go 1.22+. Wails v2 is the proven, community-supported release; v3 was still in alpha/RC at the knowledge cutoff and should not be used. Svelte 5 is the preferred frontend framework over React or Vue — it produces the smallest compiled output (~50KB) with zero runtime overhead, and Wails generates TypeScript bindings automatically. The pure-Go SQLite driver (modernc.org/sqlite) is mandatory over mattn/go-sqlite3 because the CGo-based driver breaks Docker cross-compilation; every build target would need a matching C toolchain for the SQLite dependency in addition to the Wails toolchain.

For the Docker/VNC runtime, the standard pattern is Xvfb + x11vnc + websockify + noVNC + supervisord on a debian:bookworm-slim base image. This is well-established for containerized GUI apps. The Wails binary runs as DISPLAY=:99; noVNC exposes a browser-accessible interface on port 8080. The Docker image will be 400–600MB due to the WebKit2GTK runtime dependency — this is unavoidable with Wails on Linux.

**Core technologies:**
- Go 1.22+ with Wails v2.9.x: desktop GUI framework — project-specified, proven cross-platform
- Svelte 5 + TypeScript + Vite: frontend UI — smallest bundle, Wails-native support, TS bindings auto-generated
- modernc.org/sqlite + GORM v2: local state storage — pure Go, no CGo, cross-compilation safe
- oapi-codegen v2: API client codegen from GameVault OpenAPI spec — avoids 3000+ line manual client
- debian:bookworm-slim + Xvfb + x11vnc + noVNC + supervisord: Docker/VNC runtime — standard headless GUI pattern
- go-keyring (zalando): cross-platform credential storage — OS keychain on macOS/Windows/Linux

### Expected Features

The feature set has a clear dependency chain: server configuration is the entry point, authentication gates everything else, and the core game loop (browse → download → install → launch) must be complete before any differentiators are useful. Cloud saves and progress tracking require the game launch flow to exist first, because they hook into launch/exit events. Multi-server support is independently implementable once the Profile Store exists.

**Must have (table stakes):**
- Server URL configuration + connection health check — entry point; without this nothing works
- User authentication (login/register) — gates entire library
- Game library browsing with search and filter — core purpose
- Game detail view (metadata, cover art) — prerequisite for download decision
- Game download with progress display and resume support — primary action; resume is mandatory for large files
- Game installation tracking (local state) — required to enable launch
- Game launch — the end-to-end goal of the app
- User profile view + password change — basic account management
- Offline/disconnected state handling — degrades gracefully on network loss

**Should have (differentiators):**
- Multiple server support (saved connection profiles) — power-user feature; Electron app handles this poorly
- Docker + VNC/noVNC distribution — core differentiator for NAS/headless users
- Cloud save sync (upload on exit, download on launch) — Plus feature; high user value
- Progress tracking (playtime timestamps to backend) — Plus feature
- Manual cloud save trigger — complements auto-sync; crash recovery

**Defer to v2+:**
- Achievement tracking — Plus feature, endpoint existence unconfirmed; verify before building
- Game notes/personal tags — local-only, low priority
- Per-server "new games" notification state — nice-to-have
- Multiple simultaneous downloads — serial queue sufficient for v1
- Save conflict resolution UI beyond basic prompt — basic newer-wins is acceptable for v1

### Architecture Approach

The architecture follows a strict layered pattern where Go owns all state and business logic, and the Svelte frontend is a pure rendering layer. Wails bindings expose typed Go methods to TypeScript; Wails events push async updates (download progress, sync status, auth expiry) from Go to the frontend. No frontend HTTP calls to the GameVault API are permitted. All services are instantiated in main.go via dependency injection — no global singletons — and the App struct is a thin delegation layer, not a monolith. SQLite is the single source of truth for all resumable state (download queue, save paths, active profile, library cache).

**Major components:**
1. App Struct (internal/app) — Wails entry point; thin delegation layer; exposes bound methods; owns service references
2. API Client (internal/api) — all HTTP communication to GameVault; auth token injection; retry/backoff; multi-server routing
3. Download Manager (internal/download) — concurrent queue; HTTP Range resume; progress events; extraction; state in SQLite
4. Cloud Save Manager (internal/cloudsave) — upload/download on launch/exit/manual; dirty flag; conflict detection; zip/unzip
5. Profile Store (internal/profile) — server profile CRUD; tokens in OS keyring; non-sensitive metadata in SQLite
6. EventBus (internal/events) — pub/sub bridge; forwards service events to Wails frontend via runtime.EventsEmit
7. SQLite State DB (internal/store) — schema: server_profiles, games_cache, downloads, save_paths, app_settings
8. Docker/VNC Layer — external wrapper; no Go code changes; Xvfb + x11vnc + noVNC + supervisord

### Critical Pitfalls

1. **CGo is mandatory for Wails — standard Go cross-compilation fails** — Accept CGo from day one; use platform-specific C toolchains per target; never set CGO_ENABLED=0; design the CI matrix for this from the start, not after feature work begins.

2. **macOS targets cannot be built from Linux (Apple SDK restriction)** — Use GitHub Actions macos-latest runners for darwin/amd64 and darwin/arm64; do not attempt Docker-based macOS cross-compilation; revise any PROJECT.md language implying Docker can build all five targets.

3. **Cloud save sync race condition causes save data loss** — Implement a local dirty flag (set on local save modification, cleared only after confirmed upload) before writing any sync code; never overwrite local without timestamp comparison; implement upload retry; always keep a local backup before overwriting.

4. **Large file download requires resume from the start** — Implement HTTP Range requests before shipping the download feature; track byte offset in SQLite; use .part files; validate with ETag on resume. Adding resume retroactively is significantly more complex.

5. **Credential plaintext storage is a trust-destroying mistake** — Use go-keyring for all tokens/passwords from the authentication phase; for Docker/headless environments where no keyring exists, fall back to AES-256-GCM encrypted file with machine-derived key — never plaintext JSON config.

6. **Wails dev vs production asset path mismatch** — Run wails build in CI from the first UI commit; test the built binary, not just the dev server; configure Vite base correctly for embedded FS.

## Implications for Roadmap

Based on research, the dependency graph is strict and the phase order is dictated by it. The build infrastructure must be correct before any feature work, because discovering CGo/macOS constraints late causes full pipeline rebuilds.

### Phase 1: Foundation + Build Infrastructure
**Rationale:** CGo cross-compilation constraints, Docker base image choice, and CI matrix structure must be resolved before any feature code exists. Discovering these late causes full rebuilds. SQLite schema must also be finalized before any service persists state — migrations are painful to retrofit.
**Delivers:** Working CI/CD pipeline that produces binaries for all targets; Wails window boots with Svelte scaffold; SQLite schema and store package; Profile Store + keyring integration; App struct skeleton with no-op bound methods.
**Addresses:** Server URL configuration (storage layer only), app settings persistence.
**Avoids:** Pitfall 1 (CGo assumption), Pitfall 2 (macOS SDK), Pitfall 4 (Linux webkit2gtk/Alpine), Pitfall 12 (wrong CI matrix), Pitfall 11 (dev/production asset mismatch caught early).
**Research flag:** Needs research-phase — CI matrix for CGo cross-compilation, Docker multi-stage setup, and Wails v2 scaffold configuration are all technically specific and easy to get wrong.

### Phase 2: API Client + Authentication
**Rationale:** All subsequent features require authenticated API access. Profile Store (built in Phase 1) provides credentials. oapi-codegen requires the GameVault OpenAPI spec — this must be fetched and evaluated before committing to codegen vs manual client.
**Delivers:** Full typed API client (or manual fallback); login/register flow; JWT token storage via keyring; backend version check on connect; frontend login view.
**Addresses:** User authentication, server connection health, user registration, user profile view, API compatibility checking.
**Avoids:** Pitfall 9 (plaintext credentials), Pitfall 10 (API breaking changes — version check implemented here).
**Research flag:** Needs research-phase — verify GameVault OpenAPI spec availability and completeness at `/api-json`; confirm endpoint paths (especially Plus-gated endpoints); confirm backend version check mechanism.

### Phase 3: Game Library + Core Loop
**Rationale:** Browse → detail → download → install → launch is the minimum viable user journey. Download resume must be built here, not deferred, because retrofitting it is significantly more complex. Game launch is the end-to-end validation of the core loop.
**Delivers:** Paginated game library with search/filter; game detail view with metadata and cover art; download queue with HTTP Range resume, progress events, .part file management, and disk cleanup; game installation tracking in SQLite; game launch via subprocess.
**Addresses:** Game library browsing, game detail view, game download (with resume), download progress display, game installation tracking, game launch, search and filter.
**Avoids:** Pitfall 6 (no resume support — built here), Pitfall 15 (orphaned partial files — .part tracking from the start), Pitfall 3 (Windows WebView2 — package with bootstrapper for first Windows release).
**Research flag:** Standard patterns for download queue and HTTP Range — skip research-phase. Verify GameVault's Accept-Ranges support and checksum endpoint availability.

### Phase 4: Account Management + Multi-Server Support
**Rationale:** Profile Store infrastructure is already built (Phase 1). Multi-server support is a differentiator and adds relatively low complexity at this stage. Multi-user local profiles share the same SQLite infrastructure.
**Delivers:** User profile view + password change; multiple saved server profiles; server switcher UI; per-server library state; multi-user local account switching.
**Addresses:** User profile, password change, multiple server support, multi-user account profiles, per-server auth state.
**Avoids:** Pitfall 10 (API changes — per-server version checking already in place from Phase 2).
**Research flag:** Standard patterns — skip research-phase.

### Phase 5: Cloud Saves + Progress Tracking (Plus Features)
**Rationale:** Game launch infrastructure exists from Phase 3. Cloud save sync hooks into launch/exit events. The dirty flag and conflict detection state machine must be designed before any sync code is written — not after.
**Delivers:** Progress tracking (playtime timestamps sent on launch/exit); cloud save path configuration per game; cloud save upload on exit + download on launch; manual sync trigger; dirty flag tracking; basic conflict resolution (newer-wins with user prompt on conflict); save path registry in SQLite.
**Addresses:** Progress tracking, cloud save sync, manual cloud save trigger, save conflict resolution (basic).
**Avoids:** Pitfall 5 (cloud save race condition — dirty flag + conflict detection designed upfront), Pitfall 14 (process monitoring complexity — document launcher limitation, provide manual fallback).
**Research flag:** Needs research-phase — verify `/api/saves` and `/api/progresses` endpoint shapes against live backend or source; verify save-path detection approach in gamevault-app reference client; confirm Plus-gated endpoint list.

### Phase 6: Docker + VNC/noVNC Distribution
**Rationale:** Requires a working binary (Phases 1–5). No Go code changes — wraps the existing binary in a container. VNC auth and display sizing must be correct before any public image is published.
**Delivers:** Multi-stage Dockerfile (linux/amd64 + linux/arm64); entrypoint.sh with Xvfb + x11vnc + websockify + noVNC + supervisord; VNC password auth via VNC_PASSWORD env var; 1920x1080 default resolution with noVNC resize=scale; auto-profile creation from GAMEVAULT_SERVER_URL env var on first run; Docker Hub/GHCR publishing via CI.
**Addresses:** Docker + VNC/noVNC distribution, headless/NAS deployment.
**Avoids:** Pitfall 4 (Alpine — use debian:bookworm-slim), Pitfall 7 (low VNC resolution), Pitfall 8 (unprotected VNC port).
**Research flag:** Standard pattern — skip research-phase for architecture. Verify current noVNC resize parameter syntax and supervisord config before implementation.

### Phase Ordering Rationale

- SQLite schema before any service: migrations retrofitted onto a schema shared by 5+ services are painful; finalize it in Phase 1 when cost is zero.
- API Client before Download Manager and Cloud Save Manager: both services require it for URLs and auth; this is a hard code dependency.
- Download + launch before cloud saves: cloud saves hook into game exit; the launch/exit lifecycle must exist first.
- Multi-server before cloud saves: cloud save paths are keyed by (serverID, gameID); the server profile concept must be stable before save paths are persisted.
- Docker last: wraps the complete binary; no Go changes needed; any feature added after Phase 6 requires rebuilding the image.
- macOS CI runners and mingw-w64 toolchain in Phase 1: discovering cross-compilation constraints after feature work begins causes full CI rebuild. Front-loading this is cheaper.

### Research Flags

Phases needing deeper research during planning (run /gsd:research-phase):
- **Phase 1:** CGo cross-compilation matrix setup, Wails v2 scaffold configuration, Docker multi-stage build pipeline — technically specific, easy to get wrong, high cost to fix late.
- **Phase 2:** GameVault OpenAPI spec availability and completeness, endpoint path confirmation, Plus-gated endpoint list, backend version check mechanism.
- **Phase 5:** `/api/saves` and `/api/progresses` endpoint shapes, save-path detection approach (cross-reference gamevault-app source), Plus feature gating confirmation.

Phases with standard patterns (skip research-phase):
- **Phase 3:** HTTP Range downloads, download queue state machines, subprocess launch — well-documented Go patterns.
- **Phase 4:** Multi-server profile switching, per-server state isolation — straightforward extension of Phase 1 infrastructure.
- **Phase 6:** Xvfb + x11vnc + noVNC + supervisord Docker pattern — established and well-documented.

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | MEDIUM | Wails v2 + Svelte is established community pattern. All version numbers must be verified at project start — training data cutoff August 2025. modernc/sqlite pure-Go rationale is HIGH confidence (well-documented). |
| Features | MEDIUM | Core feature set (auth, browse, download, launch) is HIGH confidence from GameVault public repo knowledge. Plus-feature endpoint shapes (saves, progresses, achievements) are LOW confidence — must be verified against live backend or source. |
| Architecture | MEDIUM-HIGH | Wails v2 binding model, EventsEmit pattern, and Go service patterns are HIGH confidence. GameVault-specific API response shapes and save-path conventions are unverified. |
| Pitfalls | MEDIUM-HIGH | Core pitfalls (CGo, macOS SDK, download resume, cloud save races, credential storage, VNC auth, webkit2gtk) are HIGH confidence from established domain knowledge. Specific API names (Wails runtime functions, noVNC parameters) are MEDIUM — verify before use. |

**Overall confidence:** MEDIUM

### Gaps to Address

- **GameVault OpenAPI spec:** Confirm it exists and is complete at `/api-json` before committing to oapi-codegen. If incomplete, handwritten client using interfaces is the fallback. Address in Phase 2 planning.
- **Plus endpoint shapes:** `/api/saves`, `/api/progresses`, achievement endpoints — verify against live backend source (`src/modules/`) before Phase 5 planning. These are the least verified areas of FEATURES.md.
- **Save-path detection approach:** The official gamevault-app Electron client likely has a per-game save-path config the user sets. Cross-reference the source before designing the Cloud Save Manager's path registry. Address before Phase 5 implementation.
- **WebView2 bundling flag:** The Wails `-webview2` flag name for embedding the WebView2 bootstrapper must be verified against current Wails v2 docs before the first Windows release. Training data has this at MEDIUM confidence.
- **Wails v3 status:** If the project timeline extends into late 2026, re-evaluate Wails v3 stability. The "revisit" criterion is a stable release tag + finalized migration guide.
- **GameVault backend version compatibility:** Define the `minBackendVersion` and `maxTestedBackendVersion` window during Phase 2 planning once the API is verified. This determines how aggressively to defend against breaking changes.

## Sources

### Primary (HIGH confidence)
- Wails v2 documentation (training knowledge, cutoff August 2025) — binding model, OnStartup context, EventsEmit, embed.FS behavior
- Go standard library — http.Client, goroutine patterns, context propagation, os/exec
- RFC 7233 — HTTP range requests (stable standard)
- Apple developer documentation — macOS SDK redistribution restrictions
- WebKit2GTK/glibc dependency — Wails Linux renderer constraint, well-established

### Secondary (MEDIUM confidence)
- GameVault backend (github.com/Phalcode/gamevault-backend, training data) — NestJS + OpenAPI spec likely available; endpoint paths inferred
- GameVault app (github.com/Phalcode/gamevault-app, training data) — feature set and UI patterns inferred
- modernc.org/sqlite — pure-Go SQLite driver, widely used for no-CGo cross-compilation
- zalando/go-keyring — cross-platform keyring abstraction
- oapi-codegen v2 — current community standard for OpenAPI → Go codegen
- Xvfb + x11vnc + noVNC + supervisord Docker GUI pattern — established, well-documented

### Tertiary (LOW confidence)
- GameVault Plus feature set (gamevau.lt/plus, not fetched) — cloud saves, achievements, advanced profiles; needs live verification
- GameVault official API docs (gamevau.lt/docs, not fetched) — endpoint paths assumed from training data; must be verified
- Specific Wails function names (runtime.WindowGetSystemTheme, -webview2 flag) — likely correct but need doc verification

---
*Research completed: 2026-03-11*
*Ready for roadmap: yes*
