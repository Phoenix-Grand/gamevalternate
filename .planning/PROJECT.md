# GameVault Go Client

## What This Is

A cross-platform desktop application written in Go (using Wails) that serves as a drop-in replacement for the official gamevault-app Electron client. It connects to any self-hosted GameVault backend and exposes all base and Plus features — profiles, cloud saves, game library management — with no feature tiers. Targets Linux, Windows, and macOS.

## Core Value

Users can manage and launch their self-hosted GameVault library from any platform with full feature parity to the official client, including cloud save sync and multi-user profiles.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Full GameVault backend API compatibility (base + Plus endpoints)
- [ ] Cross-platform GUI using Wails (Linux, Windows, macOS)
- [ ] Multiple GameVault server support (saved connection profiles)
- [ ] Multi-user account profiles on same install
- [ ] Full cloud save sync — upload and download on game launch/exit + manual
- [ ] Game library browsing, searching, filtering
- [ ] Game download and installation management
- [ ] User authentication (login/logout per server)
- [ ] Docker build pipeline (cross-compile for all targets)
- [ ] Docker runtime with VNC/noVNC (browser-accessible GUI, no X11 required)
- [ ] GitHub Actions CI/CD for builds and releases

### Out of Scope

- Feature tiers — all functionality included without paywalls
- Separate backend server — this is client-only, connects to existing GameVault backend
- Mobile app — desktop only for now

## Context

- **Reference client**: https://github.com/Phalcode/gamevault-app (Electron-based, used as feature baseline)
- **Backend**: https://github.com/Phalcode/gamevault-backend (must be fully compatible with all documented API endpoints including Plus tier)
- **GUI framework**: Wails v2 — Go backend, web frontend (HTML/CSS/JS), packages as single native binary
- **Container GUI**: VNC/noVNC approach — Docker container exposes a browser-accessible UI, no X11/XQuartz required from the user
- **Distribution**: Native binaries (via GitHub Actions releases) + Docker image (for users who prefer containerized workflow)

## Constraints

- **Tech Stack**: Go + Wails v2 — no Electron, no CGo where avoidable
- **API Compatibility**: Must implement all GameVault backend REST API endpoints (base and Plus) as documented
- **Docker Build**: Multi-stage Dockerfile must cross-compile for linux/amd64, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64
- **No Backend**: Client only — does not bundle or modify the GameVault backend
- **Single Binary**: Each platform target should ship as a single executable (Wails default)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Wails over Fyne/Gio | Web frontend flexibility, modern look, easier to match gamevault-app UX | — Pending |
| VNC/noVNC for Docker GUI | Works everywhere without X11 setup on host | — Pending |
| Multiple server support | Power users often run separate GameVault instances (home/work/friend) | — Pending |
| No feature tiers | Simplicity — one codebase, one experience | — Pending |

---
*Last updated: 2026-03-11 after initialization*
