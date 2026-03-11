# Architecture Patterns

**Domain:** Cross-platform desktop game library client (Go/Wails v2)
**Researched:** 2026-03-11
**Confidence:** MEDIUM — Based on Wails v2 documentation knowledge (cutoff Aug 2025), Go ecosystem patterns, and GameVault API reference. External verification was unavailable during this research session.

---

## Recommended Architecture

The application follows a layered architecture where the Go backend owns all business logic and state, and the web frontend is a pure rendering layer that calls into Go via Wails bindings. This is the idiomatic Wails pattern and the critical architectural decision that separates it from an Electron app (where Node.js and renderer share state more freely).

```
┌─────────────────────────────────────────────────────────────────┐
│                        Wails Window                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    Frontend (JS/HTML)                      │  │
│  │   Library UI  │  Downloads UI  │  Settings UI  │  Auth UI  │  │
│  │                    Svelte or Vue SPA                       │  │
│  └────────────────────────┬──────────────────────────────────┘  │
│                            │  Wails JS Bindings (auto-generated) │
│  ┌─────────────────────────▼──────────────────────────────────┐  │
│  │                    App Struct (main.go)                    │  │
│  │         Wails context holder, entry point for bindings     │  │
│  └──┬──────────────┬──────────────┬──────────────┬───────────┘  │
│     │              │              │              │               │
│  ┌──▼───┐  ┌───────▼──┐  ┌───────▼──┐  ┌───────▼──┐           │
│  │ API  │  │Download  │  │CloudSave │  │ Profile  │           │
│  │Client│  │Manager   │  │Manager   │  │ Store    │           │
│  └──┬───┘  └───────┬──┘  └───────┬──┘  └───────┬──┘           │
│     │              │              │              │               │
│  ┌──▼───────────────▼──────────────▼──────────────▼──────────┐  │
│  │                  Shared Infrastructure                     │  │
│  │      EventBus │ Config (keyring) │ SQLite state DB         │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Component Boundaries

| Component | Responsibility | Communicates With |
|-----------|---------------|-------------------|
| **App Struct** | Wails entry point; aggregates service structs; exposes bound methods to frontend | All services (owns references); Wails runtime |
| **API Client** | All HTTP calls to GameVault backend; auth token management; retry/backoff; multi-server routing | Profile Store (reads server URL + token); EventBus (emits auth errors) |
| **Download Manager** | Concurrent game downloads; chunk-resume via HTTP Range; progress tracking; post-download extraction | API Client (constructs download URLs); EventBus (emits progress events); filesystem |
| **Cloud Save Manager** | Upload/download save files on game launch/exit and manual trigger; file watching; conflict detection | API Client (PUT/GET saves endpoints); EventBus (emits sync status); filesystem |
| **Profile Store** | Persists server connection profiles (URL, username, encrypted token); multi-server, multi-user state | keyring/OS credential store; SQLite (non-secret metadata); App Struct |
| **EventBus** | In-process pub/sub for backend-to-frontend push events (progress, errors, status) | Wails runtime (emits events to frontend via `runtime.EventsEmit`); all services (publish) |
| **Config / SQLite** | Persistent non-secret app state: active server, library metadata cache, download queue state | All services (read/write) |
| **Docker / VNC Layer** | Container runtime: Xvfb + Openbox/Fluxbox + x11vnc + noVNC proxy; exposes port 8080 | External to Go app; wraps the compiled binary as DISPLAY=:99 |

---

## Data Flow

### Authentication Flow

```
Frontend: login form submit
  → App.Login(serverURL, username, password)
    → APIClient.Authenticate(serverURL, credentials)
      → POST /api/users/login → JWT token
    → ProfileStore.SaveProfile(serverURL, username, token)
      → keyring.Set(key, token)
      → SQLite: upsert profile record
    → EventBus.Emit("auth:success", profileID)
  → Frontend: navigates to library view
```

### Library Browse Flow

```
Frontend: requests game list
  → App.GetGames(filters)
    → APIClient.ListGames(activeProfile.ServerURL, filters)
      → GET /api/games?search=&tags=
    → returns []GameDTO to frontend (JSON via Wails binding)
  → Frontend: renders grid
```

### Download Flow

```
Frontend: user clicks Download
  → App.StartDownload(gameID)
    → DownloadManager.Enqueue(gameID)
      → APIClient.GetDownloadURL(gameID) → presigned or direct URL
      → goroutine: HTTP GET with Range header support
        → writes chunks to tmp file
        → emits EventBus("download:progress", {gameID, bytes, total, speed})
      → on complete: moves tmp → install path, updates SQLite
    → EventBus → runtime.EventsEmit → Frontend: live progress bar
```

### Cloud Save Sync Flow (upload on exit)

```
Game process exits (watched via os/exec or process monitor)
  → CloudSaveManager.OnGameExit(gameID)
    → stat save directory → compute checksum
    → APIClient.GetSaveMetadata(gameID) → server checksum
    → if local != server:
        if local.mtime > server.mtime: upload
        else: conflict → emit EventBus("savesync:conflict", details)
    → APIClient.UploadSave(gameID, zipPayload)
    → EventBus.Emit("savesync:complete", gameID)
```

### Multi-Server State Flow

```
Profile Store holds: []ServerProfile{ID, URL, username, encryptedToken}
Active server = single pointer in SQLite (last_active_profile_id)

Frontend: switch server
  → App.SetActiveServer(profileID)
    → ProfileStore.SetActive(profileID)
    → APIClient.SetBaseURL(profile.ServerURL)
    → APIClient.SetToken(keyring.Get(profile.TokenKey))
    → App.GetGames() — frontend refreshes with new server context
```

---

## Component Definitions (Detailed)

### App Struct (internal/app)

The central Wails binding struct. In Wails v2, all methods on this struct tagged as exported become callable from JavaScript automatically. It is NOT a service itself — it delegates to services and returns results.

```go
type App struct {
    ctx           context.Context  // set by OnStartup
    apiClient     *api.Client
    downloadMgr   *download.Manager
    cloudSaveMgr  *cloudsave.Manager
    profileStore  *profile.Store
    eventBus      *events.Bus
}

// Wails lifecycle
func (a *App) OnStartup(ctx context.Context)   { a.ctx = ctx }
func (a *App) OnShutdown(ctx context.Context)  { /* flush queues */ }

// Bound methods (frontend callable)
func (a *App) Login(serverURL, user, pass string) error
func (a *App) GetGames(filters GameFilters) ([]GameDTO, error)
func (a *App) StartDownload(gameID int) error
func (a *App) GetDownloadProgress(gameID int) DownloadProgress
func (a *App) SyncSaves(gameID int) error
func (a *App) GetProfiles() []ProfileDTO
func (a *App) SetActiveProfile(profileID string) error
```

Key rule: App methods must return `(T, error)` — Wails serializes return values to JSON. Keep them thin; all logic lives in service packages.

### API Client (internal/api)

Owns all HTTP communication. Structured around a base `Client` type with per-resource methods.

```
internal/api/
  client.go       — http.Client wrapper, base URL, auth header injection, retry
  auth.go         — Login, Logout, RefreshToken
  games.go        — ListGames, GetGame, GetDownloadInfo
  saves.go        — GetSave, UploadSave, ListSaves
  users.go        — GetProfile, UpdateProfile
  media.go        — GetBoxArt, GetBackground (cached)
  types.go        — DTO structs matching backend JSON schema
```

Design decisions:
- Single `http.Client` with configurable transport (allows mock in tests)
- Auth token stored externally (ProfileStore injects on each request via `http.RoundTripper`)
- Retry: exponential backoff with jitter for 5xx and network errors; no retry on 401 (bubble to auth flow)
- Rate limiting: token bucket (golang.org/x/time/rate) to avoid hammering self-hosted servers
- Multi-server: `Client` is initialized per-profile, not singleton; App holds active client reference

### Download Manager (internal/download)

Concurrent download queue with resume support.

```
internal/download/
  manager.go      — queue, worker pool, public API
  job.go          — single download job state machine
  resume.go       — HTTP Range request logic, .part file management
  extractor.go    — zip/rar/7z extraction post-download
  progress.go     — progress calculation, speed smoothing
```

State machine per job:
```
Queued → Downloading → Paused → Downloading (resume)
                     → Extracting → Complete
                     → Failed (retryable or permanent)
```

Key decisions:
- Configurable worker pool size (default 2 concurrent downloads)
- Resume: download to `{gameID}.part`, track byte offset in SQLite, use `Range: bytes={offset}-`
- Progress events: throttled to emit at most every 500ms via time-based debounce
- Extraction: runs in separate goroutine after download; progress emitted for large archives
- Downloads survive app restart: queue state persisted in SQLite, `.part` files preserved

### Cloud Save Manager (internal/cloudsave)

```
internal/cloudsave/
  manager.go      — sync orchestration, game lifecycle hooks
  watcher.go      — filesystem watcher (fsnotify) for save directories
  conflict.go     — conflict detection and resolution strategy
  archive.go      — zip save directory for upload/unzip on download
```

Sync triggers:
1. Game launch: download saves from server before process starts
2. Game exit: upload saves after process exits
3. Manual: user-initiated sync from UI
4. File watch: optional auto-sync on save file change (configurable)

Conflict resolution:
- Primary strategy: last-write-wins (compare mtime)
- Secondary: expose conflict to user via EventBus event with both versions
- Never silently overwrite: always keep local backup before overwrite

Save path registry: stored in SQLite keyed by `(serverID, gameID)` → local path mapping. User configures once; remembered.

### Profile Store (internal/profile)

```
internal/profile/
  store.go        — CRUD for server profiles
  keyring.go      — OS credential store abstraction (zalando/go-keyring)
  types.go        — Profile struct
```

Storage split:
- Sensitive (token, password): OS keyring via `go-keyring` (Keychain on macOS, Secret Service on Linux, DPAPI on Windows)
- Non-sensitive (server URL, username, display name, active flag): SQLite

Multi-server: profiles are rows in `server_profiles` table. Active profile tracked as app setting. Switching servers: load profile from SQLite, fetch token from keyring, reinitialize APIClient.

### EventBus (internal/events)

Thin pub/sub over Wails `runtime.EventsEmit`. The bus receives events from backend services and forwards them to the Wails runtime (which pushes to frontend via WebSocket).

```go
type Bus struct {
    ctx context.Context // Wails context required for EventsEmit
}

func (b *Bus) Emit(event string, data any) {
    runtime.EventsEmit(b.ctx, event, data)
}
```

Standard event names (frontend subscribes via Wails JS `Events.On`):
```
download:progress    {gameID, bytesDownloaded, totalBytes, speedBps, eta}
download:complete    {gameID, installPath}
download:error       {gameID, error, retryable}
savesync:started     {gameID, direction}
savesync:complete    {gameID}
savesync:conflict    {gameID, localMtime, serverMtime}
auth:expired         {profileID}
```

### SQLite State DB (internal/store)

Single SQLite file at platform-appropriate user data directory:
- Linux: `$XDG_DATA_HOME/gamevault-go/state.db` (fallback `~/.local/share/gamevault-go/`)
- Windows: `%APPDATA%\gamevault-go\state.db`
- macOS: `~/Library/Application Support/gamevault-go/state.db`

Schema tables:
```sql
server_profiles (id, display_name, server_url, username, active, created_at)
games_cache     (server_id, game_id, title, metadata_json, cached_at)
downloads       (server_id, game_id, status, bytes_downloaded, total_bytes, install_path, part_path)
save_paths      (server_id, game_id, local_path, last_synced_at, last_checksum)
app_settings    (key, value)
```

Use `database/sql` with `modernc.org/sqlite` (pure Go, no CGo) or `mattn/go-sqlite3` (CGo). Prefer `modernc.org/sqlite` to satisfy the no-CGo constraint.

---

## Docker / VNC Architecture

The Docker deployment is a wrapper around the native binary, not a different application mode. The Go app itself has no Docker-specific code paths.

```
Docker Container
┌─────────────────────────────────────────────────────┐
│  entrypoint.sh                                       │
│    1. Xvfb :99 -screen 0 1280x800x24 &              │
│    2. openbox --display :99 &    (window manager)    │
│    3. DISPLAY=:99 ./gamevault-go &  (Wails app)      │
│    4. x11vnc -display :99 -forever -nopw &           │
│    5. websockify 8080 localhost:5900 &   (noVNC)     │
│                                                      │
│  Port 8080 → noVNC web UI (browser connects here)    │
│  Volume /data → SQLite state DB                      │
│  Volume /games → game install directory              │
│  Env: GAMEVAULT_SERVER_URL, GAMEVAULT_USERNAME, etc  │
└─────────────────────────────────────────────────────┘
```

Base image: `debian:bookworm-slim` (not Alpine — Wails uses WebKit2GTK on Linux which needs glibc).

Dockerfile stages:
1. `builder` stage: Go cross-compile (linux/amd64 and linux/arm64 via GOARCH)
2. `runtime` stage: Xvfb, x11vnc, websockify, WebKit2GTK runtime libs, noVNC static files

Critical: Wails on Linux renders via WebKit2GTK (libwebkit2gtk-4.0). This must be in the runtime image. The Docker image will be ~400-600MB for this reason. There is no way around it without switching away from Wails.

Auto-configuration from environment variables: on first startup in container mode, if `GAMEVAULT_SERVER_URL` is set, the app should auto-create a profile. Detect container via `/.dockerenv` or env var `GAMEVAULT_DOCKER=1`.

---

## Frontend Architecture

The frontend is a SPA embedded in the Wails binary via `//go:embed frontend/dist`.

```
frontend/
  src/
    lib/
      api.ts          — wailsjs/ auto-generated bindings wrapper
      stores/         — Svelte stores (or Pinia for Vue) for UI state
        library.ts    — game list, filters, pagination
        downloads.ts  — download queue state
        auth.ts       — active profile, login state
      components/
        GameCard.vue
        DownloadBar.vue
        ServerSwitcher.vue
    views/
      Library.vue
      GameDetail.vue
      Downloads.vue
      Settings.vue
      Login.vue
    wailsjs/          — auto-generated by `wails generate module`
      go/
        main/         — bound App methods as TypeScript functions
      runtime/        — Wails EventsOn, EventsOff, etc.
```

Frontend framework recommendation: Vue 3 + Pinia (aligns with gamevault-app reference client which is Angular/similar; Vue is closer in DX and component model than React or raw JS). Svelte is also viable and produces smaller bundles.

The frontend MUST NOT make direct HTTP calls to the GameVault backend. All API calls route through Go backend bindings. This is the key architectural discipline — the frontend is a rendering shell, not an API client.

---

## Patterns to Follow

### Pattern 1: Service Initialization via Dependency Injection

All services instantiated in `main.go` before `wails.Run()`. No global singletons. App struct owns service references.

```go
func main() {
    db := store.Open(store.DefaultPath())
    profileStore := profile.NewStore(db)
    eventBus := events.NewBus() // ctx set on startup
    apiClient := api.NewClient(profileStore)
    downloadMgr := download.NewManager(apiClient, eventBus, db)
    cloudSaveMgr := cloudsave.NewManager(apiClient, eventBus, db)

    app := &App{
        apiClient:    apiClient,
        downloadMgr:  downloadMgr,
        cloudSaveMgr: cloudSaveMgr,
        profileStore: profileStore,
        eventBus:     eventBus,
    }
    wails.Run(&options.App{
        OnStartup: app.OnStartup,
        Bind:      []interface{}{app},
    })
}
```

### Pattern 2: Wails Context Propagation

Wails passes `context.Context` via `OnStartup`. The EventBus and any service needing `runtime.EventsEmit` must receive this context. Store it; do not copy it.

```go
func (a *App) OnStartup(ctx context.Context) {
    a.ctx = ctx
    a.eventBus.SetContext(ctx)
}
```

### Pattern 3: Progress via Events, Not Return Values

Long-running operations (downloads, sync) must not block bound methods. Start them as goroutines, return immediately, and communicate progress via EventBus.

```go
// Correct: returns immediately
func (a *App) StartDownload(gameID int) error {
    return a.downloadMgr.Enqueue(gameID)
}
// Progress delivered via: EventsEmit("download:progress", ...)
```

### Pattern 4: Error Categorization

Errors returned from bound methods appear as JavaScript exceptions. Categorize them:

```go
type AppError struct {
    Code    string `json:"code"`    // "AUTH_EXPIRED", "NETWORK_ERROR", etc.
    Message string `json:"message"`
    Retry   bool   `json:"retry"`
}
```

Frontend can pattern-match on `code` to show appropriate UX (e.g., `AUTH_EXPIRED` triggers re-login flow).

### Pattern 5: SQLite as Source of Truth for Resumable State

All state that must survive app restarts (download queue, save path registry, active profile) lives in SQLite. In-memory caches are derived from SQLite on startup. Never rely on process memory for durable state.

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: Frontend Calling Backend API Directly

**What:** Frontend JS making `fetch()` calls to the GameVault server.
**Why bad:** Bypasses auth token management, multi-server routing, error handling. Creates two API clients to maintain.
**Instead:** All API calls go through bound Go methods. Frontend is dumb rendering layer.

### Anti-Pattern 2: Global State in Go Packages

**What:** Package-level `var activeClient *api.Client` used across packages.
**Why bad:** Untestable, race conditions, makes multi-server switching fragile.
**Instead:** Dependency injection through struct constructors. Pass explicit references.

### Anti-Pattern 3: Blocking Bound Methods

**What:** `func (a *App) DownloadGame(gameID int) error` that blocks until download completes.
**Why bad:** Freezes the Wails window; frontend JS awaits the call, blocking the UI thread.
**Instead:** Return immediately, emit progress events. Frontend subscribes to events.

### Anti-Pattern 4: Storing Tokens in SQLite Plaintext

**What:** Saving JWT tokens directly in the SQLite database.
**Why bad:** Database file is user-readable plaintext; tokens are credentials.
**Instead:** Use OS keyring (go-keyring). SQLite stores only non-sensitive profile metadata.

### Anti-Pattern 5: Single Wails Binding Struct as Monolith

**What:** All 50+ methods on one `App` struct, all logic inline.
**Why bad:** Untestable, unmaintainable, methods share global mutable state implicitly.
**Instead:** Thin App struct delegates to focused service packages. Services are independently testable.

### Anti-Pattern 6: Alpine Linux for Docker

**What:** Using `alpine` as Docker base image.
**Why bad:** Wails on Linux requires WebKit2GTK which depends on glibc. Alpine uses musl libc; WebKit2GTK won't run.
**Instead:** `debian:bookworm-slim` or `ubuntu:22.04-minimal`.

---

## Suggested Build Order (Phase Dependencies)

The components have a strict dependency graph that determines what must be built first:

```
Phase 1: Foundation
  SQLite schema + store package
  Profile Store + keyring integration
  App Struct skeleton (no-op bound methods)
  Wails window boots, frontend scaffold loads
  (No external calls yet — testable locally)

Phase 2: API Client
  Requires: Profile Store (for credentials)
  HTTP client, auth flow, game list endpoint
  Frontend: login view + basic library grid
  (Can test against real GameVault server)

Phase 3: Download Manager
  Requires: API Client (download URLs), SQLite (queue state)
  Queued downloads, progress events, resume
  Frontend: download queue UI, progress bars

Phase 4: Cloud Save Manager
  Requires: API Client (save endpoints), SQLite (save paths)
  Upload/download saves, launch/exit hooks
  Frontend: sync status indicators

Phase 5: Multi-Server / Multi-Profile Polish
  Requires: Profile Store (already built), API Client (reinit on switch)
  Server switcher UI, profile management

Phase 6: Docker / VNC
  Requires: Working binary (Phases 1-5)
  Dockerfile, entrypoint.sh, CI/CD
  No Go code changes (wraps existing binary)
```

Critical dependency: Download Manager cannot be built before API Client. Cloud Save Manager cannot be built before API Client. Profile Store must exist before API Client (it provides credentials). SQLite schema must be finalized before any service that persists state — schema migrations are painful to retrofit.

---

## Scalability Considerations

This is a single-user desktop client connecting to a self-hosted server. Scalability concerns are different from web services:

| Concern | Practical Bound | Approach |
|---------|-----------------|----------|
| Game library size | 100-10,000 games | SQLite cache with pagination; don't load all at once |
| Concurrent downloads | Limited by disk I/O and bandwidth | Configurable worker pool; default 2 |
| Save file sizes | 1MB-10GB per game (some saves are large) | Stream uploads/downloads; don't buffer in memory |
| Multiple servers | 2-10 profiles typical | Switch reinitializes API client; no per-server goroutine pool |
| File watching | 1-10 save directories watched | fsnotify handles this fine; coalesce rapid events |

---

## Sources

- Wails v2 documentation (wails.io/docs) — HIGH confidence for binding model, OnStartup context, EventsEmit API
- Go standard library patterns — HIGH confidence for http.Client, goroutine patterns, context propagation
- `modernc.org/sqlite` — MEDIUM confidence for pure-Go SQLite (widely used; satisfies no-CGo constraint)
- `zalando/go-keyring` — MEDIUM confidence for cross-platform keyring (known working on Linux/macOS/Windows)
- `fsnotify/fsnotify` — HIGH confidence for cross-platform filesystem watching
- VNC/noVNC Docker pattern — HIGH confidence (well-established container GUI approach)
- WebKit2GTK/glibc requirement — HIGH confidence (Wails Linux renderer; known constraint)
- GameVault reference client (github.com/Phalcode/gamevault-app) — not directly inspected this session; architecture inferred from project description and GameVault documentation
