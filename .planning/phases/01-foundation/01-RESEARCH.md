# Phase 1: Foundation - Research

**Researched:** 2026-03-11
**Domain:** Wails v2 + Svelte 5 scaffold, CGo cross-compilation CI matrix, SQLite schema initialization
**Confidence:** MEDIUM-HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- Full-chrome shell is present from first launch — sidebar and layout chrome appear immediately, not just post-login
- Left sidebar navigation: Library, Downloads, Settings (icons + labels)
- Server/account indicator pinned to the bottom of the left sidebar (VS Code / Slack pattern) — clicking opens server switcher or account settings
- Main content area before server connection: centered card with "Connect to a GameVault server" prompt — URL input + Connect button
- Frontend: Svelte 5 (specified in REQUIREMENTS.md)
- SQLite driver: modernc.org/sqlite (pure-Go, no CGo beyond Wails)
- macOS CI: darwin/amd64 and darwin/arm64 require native macOS GitHub Actions runners (Apple SDK restriction)

### Claude's Discretion
- Go package layout and directory structure (cmd/, internal/, frontend/ conventions)
- CI/CD trigger policy (tags, PRs, caching strategy)
- SQLite schema field-level detail (tables are fixed: server_profiles, games_cache, downloads, save_paths, app_settings)
- CGo toolchain matrix specifics per target
- Loading/transition behavior between the connect prompt and main shell post-auth

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| FOUND-01 | App builds as native binary for Linux (amd64/arm64), Windows (amd64), macOS (amd64/arm64) via GitHub Actions CI matrix | CI matrix strategy section; CGo toolchain per target; native runner strategy for darwin and linux/arm64 |
| FOUND-02 | App window opens on Linux, Windows, and macOS using Wails v2 + Svelte 5 scaffold | Wails v2.9.x init command; community Svelte 5 template; project structure section |
| FOUND-03 | SQLite local state database initializes on first launch with schema for server_profiles, games_cache, downloads, save_paths, and app_settings | glebarez/sqlite GORM adapter; schema design; platform data dir resolution |
</phase_requirements>

---

## Summary

Phase 1 delivers the build infrastructure and app skeleton that every subsequent phase depends on. The three deliverables are: (1) a GitHub Actions CI matrix that produces release binaries for all five targets on tagged commits, (2) a Wails v2 window that opens displaying the Svelte 5 scaffold with the full app shell, and (3) a SQLite database initialized on first launch with the complete schema.

The most important decision in this phase is the CI matrix design. Wails v2 requires CGo on all platforms, which breaks standard Go cross-compilation. The consequence is that each build target needs a different approach: linux/amd64 and windows/amd64 can build on `ubuntu-latest` with appropriate cross-compiler toolchains; linux/arm64 is best served with a native `ubuntu-24.04-arm` runner (GitHub's arm64 hosted runners are publicly available as of early 2025); darwin/amd64 and darwin/arm64 require native macOS runners. There is no cross-compilation shortcut for macOS. This is the known blocker flagged in STATE.md and must be the first task in this phase.

The SQLite driver decision is locked: `modernc.org/sqlite` (pure Go, no CGo). The GORM integration path is `glebarez/sqlite`, a GORM adapter wrapping modernc.org/sqlite. The schema tables are fixed; only field-level detail is at Claude's discretion. Database file location follows platform conventions via `os.UserConfigDir()`.

**Primary recommendation:** Use a 5-job explicit include matrix (not a generic OS array), use `ubuntu-24.04-arm` for linux/arm64, use `macos-latest` for darwin targets, and use `glebarez/sqlite` for GORM integration. Initialize the Svelte 5 scaffold from a community template since Wails' built-in template predates Svelte 5.

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.22+ | Application language | Wails v2.9.3 go.mod requires Go 1.22; range-over-func, improved perf |
| Wails | v2.9.3 | Desktop GUI framework | Latest stable v2.x (2025-02-13). v3 is alpha.74 as of late 2025 — not production stable |
| Svelte | 5.x | Frontend UI | Locked decision. Smallest compiled output; Wails generates TS bindings |
| TypeScript | 5.x | Frontend type safety | Required for Wails auto-generated bindings |
| Vite | 7.x | Frontend build tool | Used by community Svelte 5 Wails templates (requires Node 20.19+ or 22.12+) |
| modernc.org/sqlite | v1.46.1 | SQLite driver | Locked decision: pure Go, no CGo, cross-compile safe |
| glebarez/sqlite | v1.11.0 | GORM adapter for modernc.org/sqlite | Provides `gorm.Open(sqlite.Open(...))` without CGo |
| gorm.io/gorm | v1.25.x | ORM + schema migration | AutoMigrate for schema init; struct-based table definitions |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| actions/checkout | v4 | CI: source checkout | Every CI job |
| actions/setup-go | v5 | CI: Go toolchain setup | Every CI job |
| actions/cache | v4 | CI: Go module cache | Speeds up repeated builds |
| softprops/action-gh-release | v2 | CI: Upload binaries to GitHub Release | On tag push, after all build jobs complete |
| Node.js | 20.19+ or 22.12+ | Frontend build runtime | Required by Vite 7 in the Svelte 5 template |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| glebarez/sqlite | mattn/go-sqlite3 | mattn requires CGo — breaks cross-compilation pipeline |
| glebarez/sqlite | raw database/sql + modernc.org/sqlite | Valid; avoids extra dependency, but requires manual schema migration code |
| ubuntu-24.04-arm for linux/arm64 | Cross-compile from ubuntu-latest with aarch64-linux-gnu-gcc | Cross-compile fails for Wails v2 (documented issue #1921, #3719); native runner is reliable |
| glebarez/sqlite | bboehmke/gorm-sqlite | Both wrap modernc; glebarez has more stars (833) and more community adoption |
| Wails community Svelte 5 template | Built-in `svelte-ts` template then upgrade | Built-in template ships Svelte 4; upgrading to Svelte 5 has known breaking API changes; community template is pre-configured |

**Installation:**

```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Scaffold with community Svelte 5 template (Svelte 5 + TS + Vite 7 + Tailwind + shadcn-svelte)
wails init -n gamevault-go -t https://github.com/bnema/wails-vite-svelte5-ts-taildwind-shadcn-template

# Go backend dependencies
go get github.com/wailsapp/wails/v2
go get modernc.org/sqlite
go get github.com/glebarez/sqlite
go get gorm.io/gorm

# Frontend (from frontend/)
# Already configured by template — npm install to install deps
```

---

## Architecture Patterns

### Recommended Project Structure

```
gamevault-go/
├── cmd/
│   └── main.go              # wails.Run() entry point — instantiates services, starts app
├── internal/
│   ├── app/
│   │   └── app.go           # App struct: Wails bindings entry point, thin delegation layer
│   ├── store/
│   │   ├── store.go         # Open DB, run migrations, return *gorm.DB
│   │   ├── schema.go        # GORM model structs for all 5 tables
│   │   └── paths.go         # Platform-aware DB file location (os.UserConfigDir)
│   └── events/
│       └── bus.go           # EventBus: wraps runtime.EventsEmit
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/  # Reusable Svelte components
│   │   │   └── stores/      # Svelte 5 $state stores
│   │   ├── routes/          # SvelteKit-style route pages (if using SvelteKit static)
│   │   │   └── +page.svelte # Main app shell
│   │   └── wailsjs/         # Auto-generated by Wails — DO NOT EDIT
│   ├── package.json
│   └── vite.config.ts
├── build/
│   ├── appicon.png          # Required by Wails — 512x512 RGBA
│   └── darwin/
│       └── Info.plist       # macOS bundle metadata
├── wails.json               # Wails project config
└── go.mod
```

### Pattern 1: Wails App Struct (Binding Entry Point)

**What:** A single exported struct whose methods are automatically exposed to the frontend as TypeScript functions.
**When to use:** Always — this is the Wails idiomatic pattern.

```go
// internal/app/app.go
type App struct {
    ctx   context.Context
    db    *gorm.DB
    // Phase 1: only db needed; later phases add apiClient, downloadMgr, etc.
}

func NewApp(db *gorm.DB) *App {
    return &App{db: db}
}

// OnStartup is called by Wails after the window is created
func (a *App) OnStartup(ctx context.Context) {
    a.ctx = ctx
}

// OnShutdown is called by Wails when the window is closed
func (a *App) OnShutdown(ctx context.Context) {
    // flush queues, close resources
}

// Example bound method — returns AppInfo to frontend
func (a *App) GetAppInfo() AppInfo {
    return AppInfo{Version: "0.1.0"}
}
```

```go
// cmd/main.go
func main() {
    db := store.Open(store.DefaultPath())
    app := app.NewApp(db)

    err := wails.Run(&options.App{
        Title:     "GameVault",
        Width:     1200,
        Height:    800,
        MinWidth:  900,
        MinHeight: 600,
        AssetServer: &assetserver.Options{
            Assets: assets, // embedded frontend/dist
        },
        OnStartup:  app.OnStartup,
        OnShutdown: app.OnShutdown,
        Bind:       []interface{}{app},
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

### Pattern 2: SQLite Schema Initialization via GORM AutoMigrate

**What:** Define schema as GORM model structs; call `db.AutoMigrate()` on startup to create tables if absent.
**When to use:** First app launch and every subsequent launch (AutoMigrate is idempotent — adds missing columns, never drops).

```go
// internal/store/schema.go
type ServerProfile struct {
    ID          uint      `gorm:"primarykey"`
    DisplayName string    `gorm:"not null"`
    ServerURL   string    `gorm:"not null;uniqueIndex"`
    Username    string    `gorm:"not null"`
    Active      bool      `gorm:"default:false"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type GameCache struct {
    ID           uint      `gorm:"primarykey"`
    ServerID     uint      `gorm:"not null;index"`
    GameID       int       `gorm:"not null"`
    Title        string
    MetadataJSON string    `gorm:"type:text"`
    CachedAt     time.Time
}

type Download struct {
    ID              uint   `gorm:"primarykey"`
    ServerID        uint   `gorm:"not null;index"`
    GameID          int    `gorm:"not null"`
    Status          string `gorm:"not null;default:'queued'"` // queued|downloading|paused|complete|failed
    BytesDownloaded int64  `gorm:"default:0"`
    TotalBytes      int64  `gorm:"default:0"`
    InstallPath     string
    PartPath        string
    UpdatedAt       time.Time
}

type SavePath struct {
    ID           uint      `gorm:"primarykey"`
    ServerID     uint      `gorm:"not null;index"`
    GameID       int       `gorm:"not null"`
    LocalPath    string    `gorm:"not null"`
    LastSyncedAt time.Time
    LastChecksum string
}

type AppSetting struct {
    Key       string `gorm:"primarykey"`
    Value     string `gorm:"type:text"`
    UpdatedAt time.Time
}
```

```go
// internal/store/store.go
func Open(dbPath string) *gorm.DB {
    if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
        log.Fatalf("failed to create data dir: %v", err)
    }
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }
    if err := db.AutoMigrate(
        &ServerProfile{},
        &GameCache{},
        &Download{},
        &SavePath{},
        &AppSetting{},
    ); err != nil {
        log.Fatalf("failed to migrate schema: %v", err)
    }
    return db
}

// DefaultPath returns platform-appropriate user data directory
// Linux:   ~/.config/gamevault-go/state.db  (or $XDG_CONFIG_HOME/gamevault-go/state.db)
// Windows: %APPDATA%\gamevault-go\state.db
// macOS:   ~/Library/Application Support/gamevault-go/state.db
func DefaultPath() string {
    configDir, err := os.UserConfigDir()
    if err != nil {
        configDir = filepath.Join(os.Getenv("HOME"), ".config")
    }
    return filepath.Join(configDir, "gamevault-go", "state.db")
}
```

### Pattern 3: GitHub Actions CI Matrix (5 Explicit Include Jobs)

**What:** Separate, explicit matrix entries per target — not a generic OS array. Each entry specifies its runner, any cross-compiler install, and the `GOOS`/`GOARCH` passed to `wails build`.
**When to use:** Always for CGo projects requiring multi-platform CI.

```yaml
jobs:
  build:
    strategy:
      fail-fast: false
      matrix:
        include:
          # linux/amd64 — native build, no cross-compiler needed
          - os: ubuntu-latest
            platform: linux/amd64
            artifact: gamevault-go-linux-amd64

          # linux/arm64 — native arm64 runner (GitHub public preview, Jan 2025)
          # Cross-compilation from amd64 fails for Wails v2 (CGo assembler errors)
          - os: ubuntu-24.04-arm
            platform: linux/arm64
            artifact: gamevault-go-linux-arm64

          # windows/amd64 — cross-compile from Linux using mingw-w64
          - os: ubuntu-latest
            platform: windows/amd64
            cc: x86_64-w64-mingw32-gcc
            packages: gcc-mingw-w64-x86-64
            artifact: gamevault-go-windows-amd64.exe

          # darwin/amd64 — must use native macOS runner (Apple SDK restriction)
          - os: macos-latest
            platform: darwin/amd64
            artifact: gamevault-go-darwin-amd64

          # darwin/arm64 — native macOS runner (macos-latest is arm64 as of late 2024)
          - os: macos-latest
            platform: darwin/arm64
            artifact: gamevault-go-darwin-arm64

    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install Node
        uses: actions/setup-node@v4
        with:
          node-version: '22'

      - name: Install Wails CLI
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Install Linux system deps (ubuntu-latest and ubuntu-24.04-arm)
        if: startsWith(matrix.os, 'ubuntu')
        run: |
          sudo apt-get update
          sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev ${{ matrix.packages || '' }}

      - name: Install macOS system deps
        if: matrix.os == 'macos-latest'
        run: |
          # WebKit is bundled with macOS SDK — no additional packages needed

      - name: Build
        env:
          CGO_ENABLED: 1
          CC: ${{ matrix.cc || '' }}
        run: wails build -platform ${{ matrix.platform }} -o ${{ matrix.artifact }}

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: ${{ matrix.artifact }}
          path: build/bin/${{ matrix.artifact }}
```

### Pattern 4: Svelte 5 App Shell Structure

**What:** The full-chrome shell with sidebar navigation, rendered immediately on app launch.
**When to use:** Phase 1 scaffold — this becomes the permanent app frame all subsequent phases populate.

```svelte
<!-- frontend/src/routes/+page.svelte (or App.svelte if not using SvelteKit) -->
<script lang="ts">
  import Sidebar from '$lib/components/Sidebar.svelte';
  import MainContent from '$lib/components/MainContent.svelte';

  // Phase 1: no server connection; main content shows connect prompt
  let isConnected = $state(false);
</script>

<div class="app-shell">
  <Sidebar />
  <main class="content-area">
    {#if !isConnected}
      <!-- Connect prompt — Phase 2 replaces this with auth flow -->
      <div class="connect-card">
        <h2>Connect to a GameVault server</h2>
        <input type="url" placeholder="https://your-server.example.com" />
        <button>Connect</button>
      </div>
    {/if}
  </main>
</div>
```

```svelte
<!-- frontend/src/lib/components/Sidebar.svelte -->
<script lang="ts">
  // Active route tracking — Phase 1 scaffold, filled by later phases
  let activeRoute = $state('library');
</script>

<nav class="sidebar">
  <!-- Top nav items -->
  <div class="nav-items">
    <button class:active={activeRoute === 'library'} onclick={() => activeRoute = 'library'}>
      <!-- Library icon + label -->
      Library
    </button>
    <button class:active={activeRoute === 'downloads'} onclick={() => activeRoute = 'downloads'}>
      Downloads
    </button>
    <button class:active={activeRoute === 'settings'} onclick={() => activeRoute = 'settings'}>
      Settings
    </button>
  </div>

  <!-- Bottom: server/account indicator -->
  <div class="server-indicator">
    <!-- Phase 1: shows "No server" or connect prompt trigger -->
    <!-- Phase 2+: shows active server URL + username -->
  </div>
</nav>
```

### Anti-Patterns to Avoid

- **`CGO_ENABLED=0` anywhere in CI:** Wails requires CGo on all platforms. Setting this to 0 silently breaks the build without obvious error in some configurations.
- **Generic OS matrix (`os: [ubuntu-latest, windows-latest, macos-latest]`):** Does not correctly handle linux/arm64; tries to build darwin on ubuntu; produces incorrect artifacts. Always use explicit `include` entries.
- **Using `mattn/go-sqlite3` instead of `modernc.org/sqlite`:** mattn requires a C compiler per target for the SQLite layer, adding toolchain complexity beyond the Wails CGo requirement.
- **Using `FROM golang:alpine` in any build Docker image:** Alpine lacks glibc; WebKit2GTK (Wails Linux renderer) will not link. Use debian-based images.
- **Committing `wailsjs/` directory:** This directory is auto-generated by Wails on build. It should be in `.gitignore`; regeneration is part of the build step.
- **Testing only `wails dev`:** Always run `wails build` and test the compiled binary in CI. Dev server uses Vite's hot-reload path; the production binary uses `embed.FS` — the two behave differently for asset paths.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| SQLite schema migration | Manual `CREATE TABLE IF NOT EXISTS` SQL | GORM AutoMigrate | AutoMigrate handles additive migrations (new columns), is idempotent, and validates struct tags — hand-rolled SQL is error-prone to maintain across schema versions |
| Platform data dir resolution | Custom `switch runtime.GOOS` path logic | `os.UserConfigDir()` | Standard library handles XDG on Linux, APPDATA on Windows, Library/Application Support on macOS — no custom logic needed |
| Wails TypeScript bindings | Manual JS/TS fetch wrappers for Go functions | Wails auto-generated `wailsjs/` | Wails generates type-safe bindings automatically; hand-rolled wrappers drift from the Go API and cause runtime type errors |
| CGo cross-compiler setup | Custom Dockerfile with compiler chains | Native GH Actions runners (darwin, linux/arm64) + mingw-w64 (windows) | Cross-compilation of Wails has known failure modes; native runners are the supported path |

**Key insight:** The SQLite and path resolution problems each have exactly one right answer in the Go standard library or the GORM ecosystem. The CGo cross-compilation problem has been solved by the community via native runners and mingw-w64; attempting a custom Docker-based solution for darwin is blocked by Apple's SDK license.

---

## Common Pitfalls

### Pitfall 1: linux/arm64 Cross-Compilation Fails for Wails v2

**What goes wrong:** Attempting to build linux/arm64 on an `ubuntu-latest` (x86_64) runner with `aarch64-linux-gnu-gcc` produces assembler errors in `runtime/cgo`. This is a documented Wails issue (#1921, #3719) with no current workaround.
**Why it happens:** Wails v2 CGo bindings for WebKit2GTK do not cleanly support aarch64 cross-compilation from a different host architecture.
**How to avoid:** Use a native `ubuntu-24.04-arm` runner. GitHub added linux/arm64 hosted runners to public repositories in January 2025 (free in public preview). The `ubuntu-24.04-arm` label is available for public repos.
**Warning signs:** CI matrix with `os: ubuntu-latest` for all Linux targets; cross-compile step with `GOARCH=arm64 CC=aarch64-linux-gnu-gcc`.

### Pitfall 2: darwin Cannot Be Cross-Compiled from Linux

**What goes wrong:** Attempting to build macOS targets from a Linux runner fails or requires the macOS SDK (which Apple prohibits redistributing outside macOS hardware).
**Why it happens:** Apple SDK redistribution restrictions. Tools like osxcross exist but cannot be shipped in a public Docker image.
**How to avoid:** Use `macos-latest` runner for both darwin/amd64 and darwin/arm64. As of late 2024, `macos-latest` on GitHub Actions points to an M1 runner (arm64 hardware), which can build both architectures natively using `wails build -platform darwin/amd64` and `wails build -platform darwin/arm64`.
**Warning signs:** Any CI job with `os: ubuntu-latest` and `platform: darwin/*`.

### Pitfall 3: Wails Built-in Svelte Template is Svelte 4

**What goes wrong:** Running `wails init -n myapp -t svelte-ts` scaffolds a Svelte 4 project. Upgrading to Svelte 5 manually involves significant breaking API changes (stores become runes, `$:` reactive statements become `$derived`, component event model changes).
**Why it happens:** Wails' built-in template was not updated for Svelte 5 as of v2.9.x.
**How to avoid:** Use the community template `https://github.com/bnema/wails-vite-svelte5-ts-taildwind-shadcn-template` which ships Svelte 5 + Vite 7 + TypeScript pre-configured.
**Warning signs:** Scaffold shows `"svelte": "^4.x"` in `frontend/package.json`.

### Pitfall 4: WebKit2GTK Version Mismatch on Build Runner vs Runtime

**What goes wrong:** The build runner (e.g., `ubuntu-latest` = Ubuntu 24.04) links against `libwebkit2gtk-4.1` while the target runtime environment (e.g., user's Ubuntu 22.04) has only `libwebkit2gtk-4.0`. The binary launches but shows a blank window or a runtime linker error.
**Why it happens:** Ubuntu 24.04 defaults to webkit2gtk-4.1; Ubuntu 22.04 uses 4.0. They are separate packages.
**How to avoid:** Explicitly install `libwebkit2gtk-4.0-dev` (not `4.1`) on the CI build runner, or use `ubuntu-22.04` as the runner for linux builds, or build on `ubuntu-22.04` which defaults to 4.0. Wails v2 targets `webkit2gtk-4.0` as its primary supported version.
**Warning signs:** `ldd gamevault-go | grep webkit` shows `libwebkit2gtk-4.1` on the build machine.

### Pitfall 5: GORM AutoMigrate Does Not Drop Columns

**What goes wrong:** A developer renames a column in the GORM struct. AutoMigrate adds the new column but leaves the old one. The DB accumulates orphaned columns; queries that depend on column order may silently use wrong data.
**Why it happens:** AutoMigrate is intentionally additive-only to prevent accidental data loss.
**How to avoid:** For Phase 1 (greenfield), AutoMigrate is safe because no data exists to protect. Design the schema correctly up front. Any future column renames require explicit `db.Migrator().RenameColumn()` calls — not just struct updates.
**Warning signs:** Struct field renamed without a corresponding migration script.

### Pitfall 6: Wails dev vs Production Binary Asset Path Mismatch

**What goes wrong:** `wails dev` serves frontend assets via Vite's dev server (live reload). `wails build` embeds `frontend/dist` via `//go:embed`. Relative asset paths that work in dev fail in production, or the embedded FS root differs from expectations.
**Why it happens:** Vite base URL configuration matters for embed compatibility.
**How to avoid:** Ensure `vite.config.ts` has `base: './'`. Run `wails build` and test the output binary in CI from the very first commit that includes any frontend code. Never rely solely on `wails dev` as a correctness check.
**Warning signs:** CI only runs `go test`, never `wails build`; first production build attempt is deferred to near release.

---

## Code Examples

Verified patterns from official sources and current documentation:

### Wails v2: Embedding Frontend Assets

```go
// cmd/main.go
// Source: Wails v2 docs — https://wails.io/docs/reference/options

//go:embed all:frontend/dist
var assets embed.FS

func main() {
    err := wails.Run(&options.App{
        AssetServer: &assetserver.Options{
            Assets: assets,
        },
        // ...
    })
}
```

### glebarez/sqlite with GORM

```go
// Source: github.com/glebarez/sqlite README
import (
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

db, err := gorm.Open(sqlite.Open("/path/to/state.db"), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}
```

### os.UserConfigDir() for Platform-Aware DB Path

```go
// Source: Go standard library — stable API
configDir, err := os.UserConfigDir()
// Linux:   $XDG_CONFIG_HOME or ~/.config
// Windows: %APPDATA%
// macOS:   ~/Library/Application Support
dbPath := filepath.Join(configDir, "gamevault-go", "state.db")
```

### Wails EventsEmit from Go to Frontend

```go
// Source: Wails v2 runtime docs
// Go side — emit from any service that holds the Wails context
runtime.EventsEmit(ctx, "example:event", map[string]any{"key": "value"})
```

```typescript
// Frontend side (wailsjs/runtime/runtime.ts auto-generated)
import { EventsOn } from '../wailsjs/runtime/runtime';
EventsOn("example:event", (data) => {
    console.log(data.key); // "value"
});
```

### Svelte 5 Runes (not Svelte 4 stores)

```svelte
<script lang="ts">
  // Svelte 5: use $state rune, not writable() stores
  let isConnected = $state(false);
  let serverURL = $state('');

  // Svelte 5: use $derived, not $: reactive statements
  let canConnect = $derived(serverURL.length > 0);
</script>
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `wails init -t svelte-ts` (Svelte 4) | Community Svelte 5 template | Svelte 5 stable Oct 2024 | Built-in template is stale; use `-t` with community template URL |
| QEMU emulation for linux/arm64 CI | Native `ubuntu-24.04-arm` runner | GitHub public preview Jan 2025 | No more QEMU overhead; native builds are fast and reliable |
| Cross-compile linux/arm64 with aarch64-gcc | Native arm64 runner | Wails CGo limitation (ongoing) | Cross-compile fails; native runner is the correct path |
| `mattn/go-sqlite3` (CGo SQLite) | `modernc.org/sqlite` (pure Go) | Stable pattern since ~2021 | Eliminates CGo surface for SQLite layer; simplifies cross-compile |
| Wails v3 alpha | Wails v2.9.x stable | v3 still alpha.74 as of Dec 2025 | Use v2; v3 not production ready |
| `gorm.io/driver/sqlite` (CGo mattn) | `github.com/glebarez/sqlite` (pure Go) | Stable since 2022 | CGo-free GORM integration |

**Deprecated/outdated:**
- `wailsapp/xgo` Docker image: last updated 2021, does not support linux/arm64; superseded by native runners.
- `CGO_ENABLED=0` for Wails: never valid; breaks the webview binding layer.

---

## Open Questions

1. **Wails v2.9.3 linux/arm64 on `ubuntu-24.04-arm` runner: verified?**
   - What we know: GitHub added ubuntu-24.04-arm as a free hosted runner for public repos (Jan 2025). Wails CGo cross-compilation from amd64 fails. Native builds should work.
   - What's unclear: No Wails-specific confirmation of ubuntu-24.04-arm compatibility was found. The runner is new enough that community documentation hasn't caught up.
   - Recommendation: Create a minimal smoke-test job on ubuntu-24.04-arm in the first CI iteration; fail-fast: false so it doesn't block the other 4 targets if there's an issue.

2. **`libwebkit2gtk-4.0-dev` vs `4.1-dev` on ubuntu-latest (24.04)**
   - What we know: Ubuntu 24.04 (ubuntu-latest as of mid-2024) may have transitioned to webkit2gtk 4.1. Wails v2 documentation targets 4.0.
   - What's unclear: Whether Wails v2.9.x links against 4.0 or 4.1, and which version `ubuntu-latest` provides.
   - Recommendation: In the CI setup step, install `libwebkit2gtk-4.0-dev` explicitly. If that package is absent on the runner, also try `libwebkit2gtk-4.1-dev`. Verify by checking `dpkg -l | grep webkit` in a test CI run.

3. **macos-latest runner architecture (amd64 vs arm64)**
   - What we know: GitHub changed `macos-latest` to point to M1 (arm64) hardware as of late 2024.
   - What's unclear: Whether a single `macos-latest` runner can produce both `darwin/amd64` and `darwin/arm64` binaries via `wails build -platform darwin/amd64` cross-targeting.
   - Recommendation: Use two separate matrix entries: one for `darwin/amd64` and one for `darwin/arm64`, both on `macos-latest`. macOS supports cross-targeting within the Apple ecosystem (amd64 from arm64 runner is possible with the right SDK flags).

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Go `testing` stdlib + `go test ./...` |
| Config file | None required — standard Go test discovery |
| Quick run command | `go test ./internal/store/... -v -count=1` |
| Full suite command | `go test ./... -v -count=1` |
| Build smoke test | `wails build -platform linux/amd64 -skipbindings` (in CI) |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FOUND-01 | CI matrix produces binaries for all 5 targets | CI smoke | `wails build -platform linux/amd64` (per-job) | Wave 0 — CI YAML |
| FOUND-02 | App window opens with Wails + Svelte 5 scaffold | Smoke (manual) | Launch built binary, observe window | Manual verify |
| FOUND-03 | SQLite DB created on first launch with all 5 tables | Unit | `go test ./internal/store/... -run TestSchemaInit -v` | Wave 0 — test file |
| FOUND-03 | All 5 table names present in DB after init | Unit | `go test ./internal/store/... -run TestTableNames -v` | Wave 0 — test file |
| FOUND-03 | DB path resolves correctly per platform | Unit | `go test ./internal/store/... -run TestDefaultPath -v` | Wave 0 — test file |

**Note on FOUND-02:** Window-open verification requires a display server (Xvfb or native OS). CI can verify `wails build` succeeds (no crash during compilation + binary exists). Full window-open verification is a manual acceptance test or can be automated with Xvfb + screenshot in a later phase.

### Sampling Rate
- **Per task commit:** `go test ./internal/store/... -v -count=1` (fast; under 5 seconds)
- **Per wave merge:** `go test ./... -v -count=1` + `wails build -platform linux/amd64`
- **Phase gate:** All 5 CI matrix jobs green on a test tag push before `/gsd:verify-work`

### Wave 0 Gaps

- [ ] `internal/store/store_test.go` — covers FOUND-03: `TestSchemaInit`, `TestTableNames`, `TestDefaultPath`
- [ ] `.github/workflows/release.yml` — covers FOUND-01: 5-job matrix building all targets
- [ ] `internal/store/schema.go` — GORM model definitions for all 5 tables

*(If schema.go is created in Wave 1 task, test file is its Wave 0 counterpart)*

---

## Sources

### Primary (HIGH confidence)
- Wails v2 documentation — binding model, OnStartup context, embed.FS, EventsEmit API
- Go standard library — `os.UserConfigDir()`, context propagation, embed directive
- GitHub Actions changelog — ubuntu-24.04-arm runner availability (Jan 2025), arm64 standard runners for private repos (Jan 2026)
- Apple developer documentation — macOS SDK redistribution restrictions

### Secondary (MEDIUM confidence)
- `pkg.go.dev/modernc.org/sqlite` — v1.46.1, published Feb 18 2026; pure-Go SQLite 3.51.2
- `github.com/glebarez/sqlite` — v1.11.0 (Mar 2024); GORM adapter wrapping modernc.org/sqlite; 833 stars
- Wails v2.9.3 release — Feb 13 2025; Go 1.22 minimum confirmed via issue #4147
- `github.com/bnema/wails-vite-svelte5-ts-taildwind-shadcn-template` — Svelte 5 + Vite 7 + Wails v2.11+ template (Nov 2024)
- Chris Wheeler blog: cross-compilation with Wails — linux/arm64 cross-compile confirmed broken; darwin cross-compile confirmed unsupported
- Wails issue #1921, #3719 — linux/arm64 cross-compilation failure documentation
- madin.dev/cross-wails — macOS confirmed requires native macOS runner

### Tertiary (LOW confidence)
- `ubuntu-24.04-arm` runner + Wails v2 — infrastructure confirmed available; Wails-specific compatibility not directly documented in community resources found

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — library choices are locked decisions or well-verified (versions confirmed via pkg.go.dev and GitHub releases)
- CI matrix strategy: HIGH — darwin runner requirement is Apple legal constraint; linux/arm64 native runner strategy is verified from GitHub changelog; windows/amd64 mingw-w64 pattern is established
- Architecture patterns: HIGH — Wails binding model and GORM AutoMigrate are well-documented
- SQLite schema: MEDIUM — table names are locked; field-level design is at discretion; GORM model patterns are standard
- linux/arm64 on ubuntu-24.04-arm: MEDIUM — runner exists and is free for public repos; Wails-specific behavior unconfirmed

**Research date:** 2026-03-11
**Valid until:** 2026-06-11 (90 days — Wails v2 is in maintenance mode; no breaking changes expected in the window)
