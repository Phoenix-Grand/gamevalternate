# Roadmap: GameVault Go Client

## Overview

Six phases that carry the project from a working build pipeline to a fully deployed, cross-platform game library client. The order is dictated by hard dependency: CI and build infrastructure must be correct before any feature code exists (CGo and macOS SDK constraints are expensive to discover late), authentication gates every API call, the core browse-download-launch loop is the minimum viable product, account management and multi-server support extend that foundation, Plus features (cloud saves, progress tracking) hook into the launch/exit lifecycle established in Phase 3, and Docker packaging wraps the complete binary at the end.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation** - Build pipeline, Wails scaffold, SQLite schema, and CI/CD matrix for all targets
- [ ] **Phase 2: API Client + Authentication** - Typed API client, login/register flow, JWT keyring storage, server health check
- [ ] **Phase 3: Game Library + Core Loop** - Browse, search, filter, download with resume, install tracking, game launch
- [ ] **Phase 4: Account Management + Multi-Server** - Profile view, password change, avatar, saved server profiles, multi-user local accounts
- [ ] **Phase 5: Plus Features** - Playtime tracking, cloud save sync (auto + manual), dirty flag, save path configuration
- [ ] **Phase 6: Docker + VNC Distribution** - Containerized GUI via noVNC, env-var server pre-configuration, CI image publishing

## Phase Details

### Phase 1: Foundation
**Goal**: The project compiles and ships binaries for all five targets; the Wails window boots; local state storage is initialized and ready for every service that follows.
**Depends on**: Nothing (first phase)
**Requirements**: FOUND-01, FOUND-02, FOUND-03
**Success Criteria** (what must be TRUE):
  1. GitHub Actions CI produces release artifacts for linux/amd64, linux/arm64, windows/amd64, darwin/amd64, and darwin/arm64 on every tagged commit
  2. The app window opens on Linux, Windows, and macOS displaying the Wails + Svelte 5 scaffold (no blank screen, no crash)
  3. On first launch, a SQLite database file is created with all required tables: server_profiles, games_cache, downloads, save_paths, app_settings
  4. The CI matrix uses CGo-enabled toolchains per target and macOS-native runners for darwin targets — no CGO_ENABLED=0 anywhere
**Plans**: 3 plans
Plans:
- [ ] 01-01-PLAN.md — Wails project scaffold, Go module, App struct, SQLite store + unit tests
- [ ] 01-02-PLAN.md — Svelte 5 app shell: sidebar navigation, connect card, server indicator
- [ ] 01-03-PLAN.md — GitHub Actions CI matrix (5 targets) + human verification checkpoint

### Phase 2: API Client + Authentication
**Goal**: Users can connect to a GameVault server, verify it is reachable, and authenticate — with credentials stored securely in the OS keyring.
**Depends on**: Phase 1
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04
**Success Criteria** (what must be TRUE):
  1. User can enter a server URL and see a connection status indicator confirming the backend is reachable before attempting login
  2. User can register a new account on the connected server from within the app
  3. User can log in with username and password; subsequent app launches do not ask for credentials again (token persisted in OS keyring)
  4. User can log out, after which the app returns to the server connection screen and no credentials remain in the keyring
**Plans**: 3 plans
Plans:
- [ ] 01-01-PLAN.md — Wails project scaffold, Go module, App struct, SQLite store + unit tests
- [ ] 01-02-PLAN.md — Svelte 5 app shell: sidebar navigation, connect card, server indicator
- [ ] 01-03-PLAN.md — GitHub Actions CI matrix (5 targets) + human verification checkpoint

### Phase 3: Game Library + Core Loop
**Goal**: Users can browse the server library, download a game with reliable resume support, and launch it — the complete end-to-end workflow of the app.
**Depends on**: Phase 2
**Requirements**: LIB-01, LIB-02, LIB-03, LIB-04, LIB-05, DL-01, DL-02, DL-03, DL-04, DL-05, LAUNCH-01, LAUNCH-02, LAUNCH-03
**Success Criteria** (what must be TRUE):
  1. User can browse the game library in grid and list views, search by title, and filter by genre, platform, and installation status
  2. User can view a game detail page showing cover art, metadata, description, and file size
  3. User can start a download and see real-time progress (bytes, speed, ETA); if the connection drops mid-download, resuming continues from the last byte offset without restarting
  4. User can pause, cancel, and restart a download; download integrity is verified via checksum after completion
  5. User can launch an installed game; the app detects when the game process exits and can uninstall a game and remove its files
**Plans**: 3 plans
Plans:
- [ ] 01-01-PLAN.md — Wails project scaffold, Go module, App struct, SQLite store + unit tests
- [ ] 01-02-PLAN.md — Svelte 5 app shell: sidebar navigation, connect card, server indicator
- [ ] 01-03-PLAN.md — GitHub Actions CI matrix (5 targets) + human verification checkpoint

### Phase 4: Account Management + Multi-Server
**Goal**: Users can manage their account profile and maintain saved connections to multiple GameVault servers, with per-server state kept isolated.
**Depends on**: Phase 3
**Requirements**: ACCT-01, ACCT-02, ACCT-03, AUTH-05, AUTH-06
**Success Criteria** (what must be TRUE):
  1. User can view their profile (username, avatar, account details) fetched from the connected server
  2. User can change their password and update their avatar and display name from within the app
  3. User can save multiple server connection profiles and switch between them; each server shows its own library and account state independently
  4. Multiple local OS users can each maintain separate GameVault profiles on the same machine install without credential or state bleed
**Plans**: 3 plans
Plans:
- [ ] 01-01-PLAN.md — Wails project scaffold, Go module, App struct, SQLite store + unit tests
- [ ] 01-02-PLAN.md — Svelte 5 app shell: sidebar navigation, connect card, server indicator
- [ ] 01-03-PLAN.md — GitHub Actions CI matrix (5 targets) + human verification checkpoint

### Phase 5: Plus Features
**Goal**: Game sessions are tracked for playtime, and cloud saves automatically sync on game launch and exit — with a dirty flag ensuring no save data is silently overwritten after a crash or network loss.
**Depends on**: Phase 4
**Requirements**: PLUS-01, PLUS-02, PLUS-03, PLUS-04, PLUS-05, PLUS-06
**Success Criteria** (what must be TRUE):
  1. When a game is launched, the app records a playtime start timestamp to the backend; when the game exits, the end timestamp is sent
  2. When a game is launched, the app downloads cloud saves from the server if the server version is newer than the local copy
  3. When a game exits, the app uploads local saves to the server; if upload fails, the dirty flag prevents any future auto-download from overwriting unsaved local changes
  4. User can configure the local save file path per game from the UI, and the setting persists across sessions
  5. User can manually trigger a cloud save sync at any time from the game detail or library view
**Plans**: 3 plans
Plans:
- [ ] 01-01-PLAN.md — Wails project scaffold, Go module, App struct, SQLite store + unit tests
- [ ] 01-02-PLAN.md — Svelte 5 app shell: sidebar navigation, connect card, server indicator
- [ ] 01-03-PLAN.md — GitHub Actions CI matrix (5 targets) + human verification checkpoint

### Phase 6: Docker + VNC Distribution
**Goal**: Users who prefer a containerized or headless deployment can run the app in a browser-accessible Docker container with no X11 setup required on the host.
**Depends on**: Phase 5
**Requirements**: FOUND-04, FOUND-05
**Success Criteria** (what must be TRUE):
  1. A Docker image built on debian:bookworm-slim starts the app via Xvfb + x11vnc + noVNC + supervisord; the UI is accessible at http://localhost:8080 in any browser with no X11 on the host
  2. Setting the GAMEVAULT_SERVER_URL environment variable pre-populates the server URL on first launch so the user goes directly to the login screen
  3. The CI pipeline builds and publishes linux/amd64 and linux/arm64 Docker images to GHCR on every tagged release
**Plans**: 3 plans
Plans:
- [ ] 01-01-PLAN.md — Wails project scaffold, Go module, App struct, SQLite store + unit tests
- [ ] 01-02-PLAN.md — Svelte 5 app shell: sidebar navigation, connect card, server indicator
- [ ] 01-03-PLAN.md — GitHub Actions CI matrix (5 targets) + human verification checkpoint

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 1/3 | In Progress|  |
| 2. API Client + Authentication | 0/TBD | Not started | - |
| 3. Game Library + Core Loop | 0/TBD | Not started | - |
| 4. Account Management + Multi-Server | 0/TBD | Not started | - |
| 5. Plus Features | 0/TBD | Not started | - |
| 6. Docker + VNC Distribution | 0/TBD | Not started | - |
