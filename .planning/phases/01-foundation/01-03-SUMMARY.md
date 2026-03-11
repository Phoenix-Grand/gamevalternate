---
phase: 01-foundation
plan: "03"
subsystem: infra
tags: [github-actions, ci-cd, wails, cross-compile, webkit2gtk, release, mingw-w64]

# Dependency graph
requires:
  - phase: 01-foundation
    provides: "Wails scaffold, SQLite store, Svelte frontend shell (Plans 01 and 02)"
provides:
  - "3-job GitHub Actions release matrix covering linux/amd64, linux/arm64, windows/amd64"
  - "Automated binary publishing to GitHub Releases on tag push via softprops/action-gh-release"
  - "Native ubuntu-24.04-arm runner for linux/arm64 (avoids Wails CGo cross-compile failures)"
affects:
  - "All future phases — every v* tag push produces release artifacts automatically"

# Tech tracking
tech-stack:
  added:
    - "GitHub Actions: ubuntu-24.04-arm native arm64 runner"
    - "softprops/action-gh-release@v2 for release publishing"
    - "libwebkit2gtk-4.1-dev + webkit2_41 build tag for Ubuntu 24.04 Noble compatibility"
    - "Native windows-latest runner for windows/amd64 (replaced mingw-w64 cross-compile)"
  patterns:
    - "Explicit matrix include (not exclude) pattern for CI — avoids accidental coverage gaps"
    - "fail-fast: false to preserve artifacts from successful targets when one fails"
    - "Pinned Wails CLI version @v2.9.3 in CI — no @latest"
    - "npm ci (not npm install) for reproducible CI installs"

key-files:
  created:
    - ".github/workflows/release.yml"
  modified: []

key-decisions:
  - "macOS dropped from CI matrix: runner issues during verification; linux/windows cover primary user base; may revisit in future phase"
  - "libwebkit2gtk-4.1-dev required on Ubuntu 24.04 Noble — 4.0-dev is not available in Noble apt repos"
  - "webkit2_41 build tag required for Wails v2 CGo on Ubuntu 24.04 — Wails hardcodes webkit2gtk-4.0 but tag overrides to 4.1"
  - "Windows switched to native windows-latest runner — mingw-w64 cross-compile produced linking failures"
  - "ubuntu-24.04-arm native runner for linux/arm64 — cross-compile from amd64 fails for Wails v2 CGo (issues #1921/#3719)"
  - "Pinned Wails CLI @v2.9.3 in CI — avoids breaking changes from @latest"

patterns-established:
  - "CI matrix: use explicit include entries, not combinatorial matrix, for multi-platform Wails builds"
  - "Platform compatibility: always verify apt package availability for the specific Ubuntu version in use before coding CI"

requirements-completed: [FOUND-01]

# Metrics
duration: 45min
completed: 2026-03-11
---

# Phase 1 Plan 03: CI/CD Release Matrix Summary

**3-job GitHub Actions release matrix (linux/amd64, linux/arm64, windows/amd64) with automated GitHub Release publishing; macOS intentionally dropped after runner issues; Ubuntu 24.04 WebKit compatibility resolved via iterative CI fixes.**

## Performance

- **Duration:** ~45 min (initial write + CI verification + iterative fixes)
- **Started:** 2026-03-11T19:19:54Z
- **Completed:** 2026-03-11
- **Tasks:** 2 of 2 complete
- **Files modified:** 1

## Accomplishments

- Created `.github/workflows/release.yml` with multi-platform build matrix
- All 3 confirmed targets (linux/amd64, linux/arm64, windows/amd64) pass green in CI
- Binary artifacts published automatically to GitHub Releases on tag push
- Ubuntu 24.04 WebKit compatibility resolved (libwebkit2gtk-4.1-dev + webkit2_41 build tag)
- Windows job switched to native runner, eliminating cross-compile linking failures
- Native ubuntu-24.04-arm runner confirmed working for linux/arm64

## Task Commits

Each task was committed atomically:

1. **Task 1: Write GitHub Actions release workflow** - `d7e5c26` (feat)
2. **Task 2: Verify complete Phase 1 foundation — CI pipeline** - Human-verified green

**CI fix commits applied during verification:**
- `3295bb4` fix(ci): use libwebkit2gtk-4.1-dev on Ubuntu 24.04 (noble)
- `7a85d6f` fix(ci): pass webkit2_41 build tag for Linux
- `84c0bde` fix(ci): zip macOS .app bundle before upload; split wails_output from artifact name
- `43897fb` fix(ci): drop macOS, switch Windows to native runner

## Files Created/Modified

- `.github/workflows/release.yml` - 3-job CI matrix triggering on `v*` tags; builds linux/amd64, linux/arm64, windows/amd64; publishes binaries to GitHub Release

## Decisions Made

- **macOS dropped from scope:** macOS CI jobs encountered runner issues during verification. Final pipeline covers linux/amd64, linux/arm64, and windows/amd64. macOS may be revisited in a future phase.
- **libwebkit2gtk-4.1-dev on Ubuntu 24.04:** The original plan specified 4.0-dev (based on research for older Ubuntu versions). Ubuntu 24.04 Noble only ships 4.1-dev; adopted 4.1-dev with the `webkit2_41` build tag.
- **Windows switched to native runner:** Original plan cross-compiled via mingw-w64 on ubuntu-latest. Native `windows-latest` runner eliminated linking issues.
- **ubuntu-24.04-arm confirmed for linux/arm64:** Native arm64 runner strategy validated — cross-compile from amd64 fails for Wails v2 CGo.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] libwebkit2gtk-4.0-dev not available on Ubuntu 24.04**
- **Found during:** Task 2 (CI verification)
- **Issue:** Ubuntu 24.04 Noble does not include libwebkit2gtk-4.0-dev; apt install failed in CI
- **Fix:** Switched to libwebkit2gtk-4.1-dev and added `webkit2_41` build tag for Wails CGo compatibility
- **Files modified:** .github/workflows/release.yml
- **Verification:** linux/amd64 and linux/arm64 CI jobs passed green
- **Committed in:** 3295bb4, 7a85d6f

**2. [Rule 1 - Bug] Windows cross-compile via mingw-w64 produced linking failures**
- **Found during:** Task 2 (CI verification)
- **Issue:** mingw-w64 cross-compile on ubuntu-latest generated linker errors for the Wails Windows build
- **Fix:** Switched windows/amd64 job to native windows-latest runner
- **Files modified:** .github/workflows/release.yml
- **Verification:** windows/amd64 CI job passed green
- **Committed in:** 43897fb

**3. [Scope Change - User Approved] macOS targets removed from matrix**
- **Found during:** Task 2 (CI verification)
- **Issue:** macOS CI jobs encountered runner issues during the test tag push
- **Resolution:** User approved dropping macOS from scope; final matrix covers 3 targets
- **Files modified:** .github/workflows/release.yml
- **Committed in:** 43897fb

---

**Total deviations:** 3 (2 auto-fixed bugs, 1 user-approved scope reduction)
**Impact on plan:** Original plan specified 5 CI targets; delivered pipeline covers 3 confirmed-green targets. macOS explicitly deferred. All fixes required for CI to function correctly.

## Issues Encountered

- Ubuntu 24.04 Noble requires libwebkit2gtk-4.1-dev (not 4.0-dev as documented in Wails v2 docs targeting older Ubuntu)
- Wails v2 CGo hardcodes webkit2gtk-4.0 internally; `webkit2_41` build tag must be passed to override on Ubuntu 24.04
- Windows cross-compilation via mingw-w64 on Linux was unreliable; native runner resolved this cleanly

## User Setup Required

None - CI pipeline is fully operational. No additional external service configuration required.

## Next Phase Readiness

- CI/CD pipeline is operational; every `v*` tag push produces release artifacts for linux/amd64, linux/arm64, and windows/amd64
- Phase 1 foundation is complete: Go backend, Svelte frontend, and release pipeline all confirmed working
- Phase 2 can begin: GameVault API client implementation
- Note: macOS builds are not currently produced by CI; macOS users must build locally until revisited

---
*Phase: 01-foundation*
*Completed: 2026-03-11*
