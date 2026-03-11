# Feature Landscape

**Domain:** GameVault Desktop Client (Go/Wails replacement for gamevault-app)
**Researched:** 2026-03-11
**Confidence note:** Web tools unavailable during this research session. Findings are based on training data (cutoff August 2025) from the GameVault GitHub repos, official docs at gamevau.lt, and community discussion. Confidence levels reflect this limitation. All LOW/MEDIUM findings should be verified against the live repos before implementation.

---

## Table Stakes

Features users expect from any GameVault client. Missing = product feels incomplete or broken.

| Feature | Why Expected | Complexity | Confidence | Notes |
|---------|--------------|------------|------------|-------|
| User authentication (login/logout) | Every server is multi-tenant; session management is the entry point | Low | HIGH | Username/password POST to `/api/users/login`, JWT returned |
| Server URL configuration | Client is useless without pointing at a backend | Low | HIGH | Must persist across sessions; multiple saved servers is a UX bonus |
| Game library browsing | Core purpose of the client | Medium | HIGH | Paginated game list from `/api/games`; shows title, cover art, metadata |
| Game detail view | Users need to see metadata before downloading | Low | HIGH | `/api/games/:id` — title, description, genres, release date, box art |
| Game download | Primary action; users download installers/archives from the server | High | HIGH | `/api/games/:id/files` delivers the file; download queue with progress |
| Download progress display | Downloads can be large (multi-GB); no feedback = assumed broken | Medium | HIGH | Byte progress, speed, ETA |
| Game installation tracking | Client must know where a game is installed locally to launch it | Medium | HIGH | Local state (install path, version) stored client-side |
| Game launch | Invoking installed games is the end-to-end goal | Medium | HIGH | Shell exec of installer/executable; may need per-game launch config |
| User registration (first-run) | Fresh installs need account creation if server allows it | Low | HIGH | POST `/api/users/register`; admin may disable this |
| Search and filter games | Library grows; browsing unfiltered is unusable at scale | Medium | HIGH | Server-side search on title; client-side filter on genre/platform/status |
| Metadata display | Genres, developer, publisher, release year — context for games | Low | HIGH | Populated from RAWG.io integration on the backend |
| Cover art / box art display | Visual library is a key UX expectation for any game client | Low | HIGH | Images served via backend media endpoints |
| User profile view | See your own account info (username, role, registered date) | Low | HIGH | GET `/api/users/:id` |
| Password change | Basic account management | Low | HIGH | PUT `/api/users/:id` |
| Server connection health | User needs to know if the server is reachable | Low | MEDIUM | Health check endpoint; surfaced in UI as connection status |
| Offline/disconnected state handling | Network goes down; client should degrade gracefully, not crash | Medium | MEDIUM | Cached state for library + clear offline indicator |

## Differentiators

Features that set the Go client apart from the official Electron app. Not table stakes, but high-value for target users.

| Feature | Value Proposition | Complexity | Confidence | Notes |
|---------|-------------------|------------|------------|-------|
| Cloud save sync (upload + download) | GameVault Plus feature; sync save files on launch/exit automatically | High | HIGH | Backend: POST/GET `/api/saves`; client must detect save paths per game engine |
| Multiple server support (saved profiles) | Power users run multiple GameVault instances; Electron app handles this poorly | Medium | MEDIUM | Local profile store: {alias, url, username, auth_token} per server |
| Multi-user account profiles on same install | Shared PC scenario; each user has their own library state | Medium | MEDIUM | Per-user local state directory; switch-user without full re-auth |
| Linux native first-class support | Official Electron app has Linux issues; Go binary is cleaner | Low | MEDIUM | Wails builds native; no Electron overhead or packaging pain |
| Docker + VNC/noVNC distribution | Headless server users, NAS users — run the client in a container | High | HIGH (design) | noVNC exposes browser-accessible GUI; no X11 on host required |
| Single native binary (no Electron) | Lower memory, faster startup, no Chromium bundled | Low | HIGH | Wails compiles to one binary; 50–150MB vs Electron's 300MB+ |
| Resource efficiency | Meaningful on low-power machines (NUC, RPi 5, NAS) | Low | MEDIUM | Go runtime vs Node.js/Chromium; measurable but hard to quantify upfront |
| Manual cloud save trigger | Upload/download on demand, not just on launch/exit | Low | MEDIUM | Complements auto-sync; useful after crash or manual save copy |
| Save conflict resolution UI | What happens when server save is newer than local? | Medium | LOW | Backend may handle conflict strategy; client needs to surface it clearly |
| Per-server notification state | Track what's new on each connected server | Medium | LOW | "New games added" since last visit per server profile |
| Game notes / personal tags | Private annotations on games (local only) | Low | LOW | Client-side only; no backend equivalent |

## GameVault Plus Features (Require Plus-enabled Backend)

These features are gated behind the GameVault Plus subscription on the backend. The Go client will expose all of them without its own tier — but they require the server to have Plus enabled.

| Feature | API Area | Client Complexity | Confidence | Notes |
|---------|----------|-------------------|------------|-------|
| Cloud save sync | `/api/saves` (POST upload, GET download, GET list) | High | MEDIUM | Save path detection is the hard part; engine-specific paths |
| Advanced user profiles / avatar | `/api/users/:id` extended fields | Low | LOW | Plus may unlock profile image uploads and extended metadata |
| Progress tracking (playtime) | `/api/progresses` — game-specific play time per user | Medium | MEDIUM | Server tracks launch/exit events; client sends timestamps |
| Achievement tracking | Unknown endpoint | High | LOW | May exist in Plus; not confirmed in training data — verify against live API |

**Save engine support (cloud saves):** GameVault's backend stores save files as-is. The client is responsible for locating save files. Common engine save paths:

| Engine | Default Save Location (Windows) | Default Save Location (Linux) |
|--------|----------------------------------|-------------------------------|
| Unity | `%APPDATA%\[game]` or `%LOCALAPPDATA%\[game]` | `~/.config/unity3d/[company]/[game]` |
| Unreal Engine | `%LOCALAPPDATA%\[game]\Saved\SaveGames` | `~/.local/share/[game]` |
| GameMaker | `%LOCALAPPDATA%\[game]` | `~/.local/share/[game]` |
| RPG Maker MV/MZ | `%APPDATA%\[game]` | `~/.local/share/[game]` |
| Generic / Custom | Game-specific; user-defined fallback | Same |

**Confidence: LOW** — GameVault cloud saves rely on user-configured paths per game, not automatic engine detection. The official client likely uses a per-game save-path config that the user sets. Verify against `gamevault-app` source for actual implementation.

## Anti-Features

Features to explicitly NOT build in v1 (and why).

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| Built-in game store / shop | This is a self-hosted client; no store exists | N/A — not in GameVault's scope |
| Social/friends system | Backend has no social graph; don't invent one | Per-server user list is enough context |
| Game mods management | Out of scope for GameVault entirely | Launch game, let user manage mods externally |
| Torrent/P2P download | Backend serves files over HTTP; no torrent protocol | Standard HTTP download with resume |
| Built-in game launcher overlays | Electron app doesn't have this; adds complexity | Launch subprocess, track PID, that's sufficient |
| Automatic metadata scraping (client-side) | Backend handles RAWG.io integration; don't duplicate | Display what backend returns |
| Platform integration (Steam/GOG linking) | Not a GameVault feature; diverges from the protocol | Out of scope entirely |
| Achievements (custom, client-invented) | Backend may have Plus achievements; don't invent new ones | Expose only what the API provides |
| Chat / messaging | Not in GameVault API | N/A |
| Game recommendations / ML features | Over-engineering for a self-hosted client | N/A |
| Admin panel / server management | The backend has its own admin UI | Client is user-facing only; admin = backend responsibility |
| Multiple simultaneous downloads (v1) | Complex queue management; defer until v2 | Serial download queue in v1 |

## Feature Dependencies

```
Server URL config → Authentication → (everything else)

Authentication → Game library browsing
Authentication → User profile view
Authentication → Cloud save sync (Plus)
Authentication → Progress tracking (Plus)

Game library browsing → Game detail view
Game detail view → Game download
Game download → Game installation tracking
Game installation tracking → Game launch
Game launch → Cloud save sync (auto on launch/exit)
Game launch → Progress tracking (send start timestamp)

Cloud save sync → Save path config per game
Cloud save sync → Upload on exit
Cloud save sync → Download before launch

Multiple server support → Per-server auth state
Multiple server support → Per-server library state
Multiple server support → Profile switching UI
```

## MVP Recommendation

### Phase 1 — Core Loop (Table Stakes)
Prioritize the minimum viable user journey: connect → browse → download → play.

1. Server URL configuration + connection health check
2. User authentication (login / register)
3. Game library browsing with search and basic filter
4. Game detail view (metadata, cover art)
5. Game download with progress display
6. Game installation tracking (local state)
7. Game launch

### Phase 2 — Account + Profile
1. User profile view + password change
2. Multiple server support (saved connection profiles)
3. Multi-user account profiles (local switch)

### Phase 3 — Plus Features (Cloud Saves + Progress)
1. Progress tracking (send play timestamps)
2. Cloud save path configuration per game
3. Cloud save upload on exit / download on launch
4. Manual cloud save trigger
5. Save conflict resolution (basic: newer-wins or prompt)

### Phase 4 — Polish + Docker
1. Docker + VNC/noVNC distribution
2. Offline state handling improvements
3. Save conflict resolution UI (if deferred)
4. Per-server "new games" notification state

### Defer to later
- Achievement tracking: verify API existence first
- Game notes / personal tags: local-only, low priority
- Per-server notification state: nice-to-have, not blocking

---

## Sources

All findings based on training data (cutoff August 2025). Live verification required.

- GameVault App (Electron reference): https://github.com/Phalcode/gamevault-app — MEDIUM confidence (training data)
- GameVault Backend: https://github.com/Phalcode/gamevault-backend — MEDIUM confidence (training data)
- GameVault official docs: https://gamevau.lt/docs — LOW confidence (not fetched, URL from training)
- GameVault Plus info: https://gamevau.lt/plus — LOW confidence (not fetched)

**Verification priorities before implementation:**
1. Fetch `/api` schema from a live backend or the OpenAPI spec in the repo to confirm all endpoint paths
2. Confirm Plus-gated endpoint list from backend source (`src/modules/` directory structure)
3. Verify save-path detection approach from `gamevault-app` source (search for "save" in the Electron app repo)
4. Confirm progress tracking endpoint shape (`/api/progresses` or similar)
5. Check whether achievement tracking exists as a Plus feature
