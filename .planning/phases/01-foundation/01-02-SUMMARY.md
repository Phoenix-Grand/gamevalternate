---
phase: 01-foundation
plan: 02
subsystem: ui
tags: [svelte5, svelte-runes, wails, tailwindcss, shadcn-svelte, typescript, vite]

# Dependency graph
requires:
  - phase: 01-foundation-01
    provides: Wails project scaffold with Go backend, frontend/ directory, shadcn-svelte template with CSS variables

provides:
  - Svelte 5 app shell with full-chrome sidebar + content area flex layout (App.svelte)
  - Left sidebar navigation with Library/Downloads/Settings nav items and SVG icons (Sidebar.svelte)
  - Bottom-of-sidebar server/account status indicator (ServerIndicator.svelte)
  - Pre-auth centered connect card with URL input and inert Connect button (ConnectCard.svelte)
  - Updated app.css with full-viewport #app layout for the shell

affects: [02-auth, 03-library, 04-downloads, 05-settings]

# Tech tracking
tech-stack:
  added: [nvm v0.39.7, node v22.22.1 (upgraded from 18 for Vite 7 compatibility)]
  patterns: [Svelte 5 $state/$derived runes, CSS custom property consumption via var(--token), flex full-chrome shell layout, component-scoped styles without hsl() wrapper]

key-files:
  created:
    - frontend/src/lib/components/Sidebar.svelte
    - frontend/src/lib/components/ServerIndicator.svelte
    - frontend/src/lib/components/ConnectCard.svelte
  modified:
    - frontend/src/App.svelte
    - frontend/src/app.css

key-decisions:
  - "CSS variables in this template are already full hsl() values (--background: hsl(...)), so component styles use var(--background) directly, not hsl(var(--background))"
  - "Template uses App.svelte (capital A) not app.svelte — main.ts imports App.svelte"
  - "Node.js 22 required for Vite 7; installed via nvm (no sudo) since system Node was 18.19.1"
  - "Wails full binary build requires system WebKit2GTK + pkg-config which are not installed in this dev environment; frontend compilation verified via npm run build"

patterns-established:
  - "Svelte 5 rune syntax: $state() for reactive variables, $derived() for computed values, onclick= attribute for events"
  - "Component CSS uses var(--token) directly (not hsl(var(--token))) because template CSS vars contain full hsl() values"
  - "App shell layout: display:flex on #app wrapper + height:100vh; sidebar fixed width 200px; content-area flex:1"

requirements-completed: [FOUND-02]

# Metrics
duration: 6min
completed: 2026-03-11
---

# Phase 1 Plan 02: App Shell Summary

**Svelte 5 full-chrome shell with sidebar nav (Library/Downloads/Settings + ServerIndicator) and centered ConnectCard, all using $state/$derived runes and shadcn-svelte CSS tokens**

## Performance

- **Duration:** 6 min
- **Started:** 2026-03-11T19:12:07Z
- **Completed:** 2026-03-11T19:17:33Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments
- Replaced scaffold App.svelte with full-chrome flex layout (sidebar left, content right)
- Created Sidebar.svelte with three SVG-icon nav items and active state tracking via $state rune
- Created ServerIndicator.svelte showing "No server / Click to connect" in sidebar bottom slot
- Created ConnectCard.svelte with URL input and Connect button (disabled when empty, inert in Phase 1)
- Fixed app.css #app constraints to allow full-viewport layout

## Task Commits

Each task was committed atomically:

1. **Task 1: App shell layout and sidebar navigation** - `8db2831` (feat)
2. **Task 2: Connect card and full Wails build verification** - `4b868f2` (feat)

**Plan metadata:** committed with this SUMMARY.md

## Files Created/Modified
- `frontend/src/App.svelte` - Root app shell: flex layout with Sidebar + ConnectCard (when !isConnected)
- `frontend/src/app.css` - Reset #app from centered 1280px to full-viewport 100vh
- `frontend/src/lib/components/Sidebar.svelte` - Left nav with Library/Downloads/Settings + ServerIndicator slot
- `frontend/src/lib/components/ServerIndicator.svelte` - Disconnected status indicator at sidebar bottom
- `frontend/src/lib/components/ConnectCard.svelte` - Centered pre-auth card with server URL input + Connect button

## Decisions Made
- CSS variables in this template contain full `hsl()` values (e.g., `--background: hsl(0 0% 100%)`), so component styles reference them as `var(--background)` directly without wrapping in `hsl()`. The plan's code samples used `hsl(var(--background))` which would have caused double-wrapping errors.
- Template uses `App.svelte` (capital A), not `app.svelte`. The plan referred to `app.svelte` but `main.ts` imports from `App.svelte` — used the correct filename.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Upgraded Node.js from 18 to 22 via nvm**
- **Found during:** Task 1 (frontend npm run build verification)
- **Issue:** System Node.js was v18.19.1; Vite 7 requires Node 20.19+ or 22.12+. Build failed with `SyntaxError: The requested module 'node:util' does not provide an export named 'styleText'`
- **Fix:** Installed nvm (user-level, no sudo), then installed Node.js v22.22.1; reinstalled npm packages to rebuild native bindings
- **Files modified:** None (nvm installs to ~/.nvm, no project files changed)
- **Verification:** `npm run build` exits 0 after Node.js upgrade
- **Committed in:** Part of Task 1 environment setup (not committed to repo)

**2. [Rule 3 - Blocking] CSS variable format adaptation**
- **Found during:** Task 1 (component authoring)
- **Issue:** Plan code samples used `hsl(var(--background))` syntax, but the shadcn-svelte template already stores full hsl() values in CSS vars (e.g., `--background: hsl(0 0% 100%)`). Using `hsl(var(--background))` would have produced invalid double-wrapping.
- **Fix:** Changed all CSS references to use `var(--token)` directly instead of `hsl(var(--token))`
- **Files modified:** Sidebar.svelte, ServerIndicator.svelte, ConnectCard.svelte, App.svelte
- **Verification:** `npm run build` exits 0 with no CSS errors
- **Committed in:** 8db2831, 4b868f2

---

**Total deviations:** 2 auto-fixed (2 blocking)
**Impact on plan:** Both necessary for correct operation. No scope creep.

## Issues Encountered
- Wails full binary build (`wails build -platform linux/amd64`) cannot complete in this WSL dev environment because `pkg-config` and `libwebkit2gtk-4.0-dev` are not installed (require sudo which is password-protected). The frontend compilation step ("Compiling frontend: Done") succeeded within the Wails build output, confirming all Svelte/TypeScript code is correct. The CGo/binary link step requires system packages only installable by the user with sudo.

## User Setup Required
None — no external service configuration required.

## Next Phase Readiness
- App shell structure is established; Phase 2 can wire real auth state into the `isConnected` variable in App.svelte
- ConnectCard.svelte is ready for Phase 2 to add real wailsjs backend calls replacing the inert console.log
- Sidebar.svelte nav items are ready for Phase 3 (Library), Phase 4 (Downloads/Settings) to add content routing

---
*Phase: 01-foundation*
*Completed: 2026-03-11*
