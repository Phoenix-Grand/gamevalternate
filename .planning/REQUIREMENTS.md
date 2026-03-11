# Requirements: GameVault Go Client

**Defined:** 2026-03-11
**Core Value:** Users can manage and launch their self-hosted GameVault library from any platform with full feature parity to the official client, including cloud save sync and multi-user profiles.

## v1 Requirements

### Foundation

- [ ] **FOUND-01**: App builds as native binary for Linux (amd64/arm64), Windows (amd64), macOS (amd64/arm64) via GitHub Actions CI matrix
- [ ] **FOUND-02**: App window opens on Linux, Windows, and macOS using Wails v2 + Svelte 5 scaffold
- [ ] **FOUND-03**: SQLite local state database initializes on first launch with schema for server_profiles, games_cache, downloads, save_paths, and app_settings
- [ ] **FOUND-04**: Docker image runs the app in a VNC/noVNC container accessible via browser (no X11 required on host)
- [ ] **FOUND-05**: User can set GAMEVAULT_SERVER_URL env var to pre-configure server on Docker first launch

### Authentication

- [ ] **AUTH-01**: User can enter a GameVault server URL and verify backend connectivity before logging in
- [ ] **AUTH-02**: User can register a new account on the connected GameVault server
- [ ] **AUTH-03**: User can log in with username and password; JWT token stored in OS keyring (never plaintext)
- [ ] **AUTH-04**: User can log out, clearing stored credentials from the OS keyring
- [ ] **AUTH-05**: User can save multiple server connection profiles (URL + credentials) and switch between them
- [ ] **AUTH-06**: Multiple local users can maintain separate GameVault profiles on the same machine install

### Game Library

- [ ] **LIB-01**: User can browse the server game library in paginated grid and list views
- [ ] **LIB-02**: User can search games by title
- [ ] **LIB-03**: User can filter games by genre, platform, and installation status
- [ ] **LIB-04**: User can view a game detail page showing cover art, metadata, description, and file size
- [ ] **LIB-05**: App displays per-game installation state (not downloaded / downloading / installed / updateable)

### Downloads

- [ ] **DL-01**: User can queue a game for download with real-time progress display (bytes downloaded, speed, ETA)
- [ ] **DL-02**: Interrupted downloads resume from the last byte offset using HTTP Range requests
- [ ] **DL-03**: User can configure the maximum number of concurrent downloads (1 = serial, N = parallel)
- [ ] **DL-04**: User can pause, cancel, and restart queued or in-progress downloads
- [ ] **DL-05**: App verifies download integrity after completion (checksum if backend provides one)

### Game Launch

- [ ] **LAUNCH-01**: User can launch an installed game from the library or game detail view
- [ ] **LAUNCH-02**: App detects when a launched game process is running and when it exits
- [ ] **LAUNCH-03**: User can uninstall a game and remove its local files from disk

### Plus Features

- [ ] **PLUS-01**: App sends playtime start/end timestamps to the GameVault backend on game launch and exit
- [ ] **PLUS-02**: App automatically downloads cloud saves from the server when a game is launched (if server version is newer)
- [ ] **PLUS-03**: App automatically uploads cloud saves to the server when a game exits
- [ ] **PLUS-04**: User can manually trigger a cloud save sync from the UI at any time
- [ ] **PLUS-05**: User can configure the local save file path per game (stored in SQLite)
- [ ] **PLUS-06**: App uses a dirty flag to prevent overwriting local saves after a crash or network loss

### Account Management

- [ ] **ACCT-01**: User can view their own profile (username, avatar, account details from server)
- [ ] **ACCT-02**: User can change their account password
- [ ] **ACCT-03**: User can update their profile (avatar image, display name)

## v2 Requirements

### Notifications

- **NOTF-01**: App displays in-app notification when new games are added to a server
- **NOTF-02**: Per-server "new games since last visit" badge on server profile switcher

### Advanced Cloud Saves

- **CSAV-01**: Save conflict resolution UI — user chooses between local and server version when timestamps conflict
- **CSAV-02**: Save history — view and restore previous cloud save versions

### Achievements

- **ACHV-01**: User can view achievements for a game (pending verification that Plus endpoint exists)
- **ACHV-02**: App syncs achievement state with backend

### Quality of Life

- **QOL-01**: User can add personal notes and tags to games (local-only)
- **QOL-02**: User can mark games as favorites with dedicated filter
- **QOL-03**: App remembers scroll position and filter state per server

## Out of Scope

| Feature | Reason |
|---------|--------|
| Feature tiers / paywalls | All functionality included — this client has no concept of tiers |
| Bundled backend server | Client-only; connects to user-hosted GameVault backend |
| Mobile app | Desktop only for v1 |
| Admin panel / server management | Backend ships its own web admin UI; client should not replicate it |
| macOS Docker cross-compile | Apple SDK redistribution restrictions prevent Linux-based macOS cross-compilation; uses native CI runners instead |
| Automatic save path detection | Per-game paths are user-configured; auto-detection is out of scope for v1 |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| FOUND-01 | Phase 1 | Pending |
| FOUND-02 | Phase 1 | Pending |
| FOUND-03 | Phase 1 | Pending |
| FOUND-04 | Phase 6 | Pending |
| FOUND-05 | Phase 6 | Pending |
| AUTH-01 | Phase 2 | Pending |
| AUTH-02 | Phase 2 | Pending |
| AUTH-03 | Phase 2 | Pending |
| AUTH-04 | Phase 2 | Pending |
| AUTH-05 | Phase 4 | Pending |
| AUTH-06 | Phase 4 | Pending |
| LIB-01 | Phase 3 | Pending |
| LIB-02 | Phase 3 | Pending |
| LIB-03 | Phase 3 | Pending |
| LIB-04 | Phase 3 | Pending |
| LIB-05 | Phase 3 | Pending |
| DL-01 | Phase 3 | Pending |
| DL-02 | Phase 3 | Pending |
| DL-03 | Phase 3 | Pending |
| DL-04 | Phase 3 | Pending |
| DL-05 | Phase 3 | Pending |
| LAUNCH-01 | Phase 3 | Pending |
| LAUNCH-02 | Phase 3 | Pending |
| LAUNCH-03 | Phase 3 | Pending |
| PLUS-01 | Phase 5 | Pending |
| PLUS-02 | Phase 5 | Pending |
| PLUS-03 | Phase 5 | Pending |
| PLUS-04 | Phase 5 | Pending |
| PLUS-05 | Phase 5 | Pending |
| PLUS-06 | Phase 5 | Pending |
| ACCT-01 | Phase 4 | Pending |
| ACCT-02 | Phase 4 | Pending |
| ACCT-03 | Phase 4 | Pending |

**Coverage:**
- v1 requirements: 33 total
- Mapped to phases: 33
- Unmapped: 0 ✓

---
*Requirements defined: 2026-03-11*
*Last updated: 2026-03-11 after roadmap creation — all 33 requirements mapped to 6 phases*
