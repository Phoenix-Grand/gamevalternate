---
phase: 1
slug: foundation
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-03-11
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (build smoke tests) + GitHub Actions CI matrix |
| **Config file** | .github/workflows/ci.yml (Wave 0 creates it) |
| **Quick run command** | `go build ./...` |
| **Full suite command** | `go build ./... && wails build -platform linux/amd64` |
| **Estimated runtime** | ~30 seconds (local build check) |

---

## Sampling Rate

- **After every task commit:** Run `go build ./...`
- **After every plan wave:** Run full build for the target platform
- **Before `/gsd:verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|-----------|-------------------|-------------|--------|
| 1-01-01 | 01 | 0 | FOUND-01 | build | `go build ./...` | ✅ W0 | ⬜ pending |
| 1-01-02 | 01 | 1 | FOUND-02 | build | `wails build -platform linux/amd64` | ✅ W0 | ⬜ pending |
| 1-01-03 | 01 | 1 | FOUND-03 | unit | `go test ./internal/db/...` | ✅ W0 | ⬜ pending |
| 1-01-04 | 01 | 2 | FOUND-01 | CI | CI matrix run | ✅ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `go.mod` — module initialized with correct dependencies
- [ ] `main.go` — Wails v2 entry point scaffold
- [ ] `internal/db/db.go` — SQLite init stub for testing
- [ ] `.github/workflows/ci.yml` — CI matrix skeleton (5 targets)

*Wave 0 establishes the build system itself — all other verification depends on it.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| App window opens on Linux | FOUND-02 | Requires display/desktop environment | Run `wails dev` on Linux with display; verify window appears |
| App window opens on Windows | FOUND-02 | Requires Windows OS | Run built .exe; verify window opens, no crash |
| App window opens on macOS | FOUND-02 | Requires macOS | Run built .app; verify window opens on both Intel and Apple Silicon |
| SQLite DB created on first launch | FOUND-03 | Requires running binary | Launch app fresh; check `~/.config/gamevault/state.db` exists with correct tables |
| CI produces artifacts for all 5 targets | FOUND-01 | Requires GitHub Actions run | Push a tag; verify release artifacts appear for all 5 targets |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
