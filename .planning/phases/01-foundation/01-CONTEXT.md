# Phase 1: Foundation - Context

**Gathered:** 2026-03-11
**Status:** Ready for planning

<domain>
## Phase Boundary

Establish the complete build pipeline: Wails v2 + Svelte 5 scaffold that opens a window on all 5 targets, SQLite local state database initialized on first launch, and GitHub Actions CI/CD matrix producing release artifacts for linux/amd64, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64. No feature logic — this phase is purely infrastructure and skeleton.

</domain>

<decisions>
## Implementation Decisions

### App Shell Structure
- Full-chrome shell is present from first launch — sidebar and layout chrome appear immediately, not just post-login
- Left sidebar navigation: Library, Downloads, Settings (icons + labels)
- Server/account indicator pinned to the bottom of the left sidebar (VS Code / Slack pattern) — clicking opens server switcher or account settings
- Main content area before server connection: centered card with "Connect to a GameVault server" prompt — URL input + Connect button

### Frontend Framework
- Svelte 5 (specified in REQUIREMENTS.md — already decided)

### SQLite Driver
- modernc.org/sqlite (pure-Go, no CGo beyond Wails) — from project research

### macOS CI
- darwin/amd64 and darwin/arm64 require native macOS GitHub Actions runners (Apple SDK restriction prevents Linux-based cross-compilation)

### Claude's Discretion
- Go package layout and directory structure (cmd/, internal/, frontend/ conventions)
- CI/CD trigger policy (tags, PRs, caching strategy)
- SQLite schema field-level detail (tables are fixed: server_profiles, games_cache, downloads, save_paths, app_settings)
- CGo toolchain matrix specifics per target
- Loading/transition behavior between the connect prompt and main shell post-auth

</decisions>

<code_context>
## Existing Code Insights

### Reusable Assets
- None — greenfield project, no existing code

### Established Patterns
- None yet — patterns established in this phase become the baseline for all subsequent phases

### Integration Points
- Phase 2 builds the auth flow into the shell established here (the "Connect" card becomes the Phase 2 login flow)
- Phase 3 fills the Library content area established by the sidebar nav
- Downloads sidebar item connects to Phase 3 download management
- Settings sidebar item used by Phase 4 (account) and Phase 5 (save paths)

</code_context>

<specifics>
## Specific Ideas

- No specific UI references provided — standard modern desktop app chrome
- The "Connect to a GameVault server" pre-auth screen should be clean and minimal (not a full wizard)

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-foundation*
*Context gathered: 2026-03-11*
