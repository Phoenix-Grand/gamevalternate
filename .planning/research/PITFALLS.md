# Domain Pitfalls

**Domain:** Go/Wails desktop game client with Docker build pipeline, cloud saves, VNC/noVNC runtime
**Researched:** 2026-03-11
**Confidence note:** Web search and WebFetch were unavailable during this research session. All findings are from training knowledge (cutoff August 2025). Confidence levels reflect this constraint. Flags are included where official doc verification is strongly recommended before implementation.

---

## Critical Pitfalls

Mistakes that cause rewrites, data loss, or shipping blockers.

---

### Pitfall 1: CGo Mandatory for Wails — Cross-Compilation is Not Simple Go Cross-Compilation

**What goes wrong:** Developers assume `GOOS=windows GOARCH=amd64 go build` works for Wails because "it's just Go." Wails v2 has hard CGo dependencies (it embeds a WebView via CGo bindings on all platforms). Attempting standard Go cross-compilation fails immediately with CGo linker errors.

**Why it happens:** The Wails documentation emphasizes ease of use for the *native* platform. Cross-compilation is a separate, more complex topic. The project constraint "no CGo where avoidable" is incompatible with Wails v2 — CGo is unavoidable for the GUI layer.

**Consequences:**
- CI/CD pipeline fails entirely if built as plain Go
- macOS cross-compilation from Linux is practically impossible at the linker level (requires macOS SDK, which has redistribution restrictions)
- Windows cross-compilation from Linux requires `mingw-w64` toolchain and specific CGo flags
- Linux arm64 from amd64 requires `aarch64-linux-gnu-gcc` and matching sysroot

**Prevention:**
- Accept CGo as a first-class build constraint from day one
- Use platform-native builders for macOS targets: GitHub Actions `macos-latest` runner (not Linux Docker)
- For Windows cross-compilation from Linux: install `gcc-mingw-w64-x86-64`, set `CC=x86_64-w64-mingw32-gcc`, `CGO_ENABLED=1`, `GOOS=windows`
- For Linux arm64: use `gcc-aarch64-linux-gnu` or QEMU-based builds
- The multi-arch Docker cross-compile pipeline must use separate builder stages per platform, not a single `GOOS` loop

**Detection warning signs:**
- Build script uses `CGO_ENABLED=0` — this breaks Wails entirely
- Dockerfile `FROM golang:alpine` for a Wails build — alpine lacks the required C libraries
- GitHub Actions matrix uses a single `ubuntu-latest` runner for all platforms including macOS

**Phase mapping:** Address in the Docker build pipeline phase before any feature work. Wrong cross-compilation architecture discovered late causes full pipeline rebuilds.

**Confidence:** HIGH (Wails CGo requirement is documented; cross-compilation complexity is well-established Go ecosystem knowledge)

---

### Pitfall 2: macOS Target Cannot Be Built From Linux (Apple SDK Restriction)

**What goes wrong:** The multi-stage Dockerfile cross-compiles all five targets (linux/amd64, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64) from a single Linux container. The Darwin/macOS targets will fail.

**Why it happens:** macOS cross-compilation from Linux requires the macOS SDK. Apple's license prohibits redistribution of the SDK outside macOS hardware. Tools like `osxcross` exist but require users to extract the SDK themselves from a macOS machine — they cannot be shipped in a public Docker image.

**Consequences:**
- Public Docker Hub image cannot legally build macOS binaries
- Any automated CI/CD that tries to do macOS in Docker will either fail legally or technically
- If discovered after building the Docker pipeline, macOS support requires a separate approach

**Prevention:**
- Use GitHub Actions `macos-latest` runner for darwin/amd64 and darwin/arm64 targets — this is the standard approach for open-source projects
- Accept that "Docker builds all targets" is not achievable for macOS without significant SDK complexity
- Revise the build matrix: Docker for Linux targets, GitHub Actions native runners for macOS and optionally Windows
- Document this explicitly in the build guide

**Detection warning signs:**
- Dockerfile references `osxcross` or `darwin-cross`
- CI matrix has `darwin` targets running on `ubuntu-latest`
- Build docs promise "build all platforms from Docker"

**Phase mapping:** This is a project constraint correction — address it before writing any CI/CD code. The PROJECT.md requirement "multi-stage Dockerfile must cross-compile for darwin/amd64, darwin/arm64" is not achievable without SDK complexity or legal grey areas.

**Confidence:** HIGH (Apple SDK licensing is well-documented; osxcross existence and limitations are well-known)

---

### Pitfall 3: WebView2 Not Bundled on Windows — Silent Failure on Fresh Installs

**What goes wrong:** On Windows, Wails v2 uses WebView2 (Chromium-based). WebView2 is pre-installed on Windows 10 (1803+) and Windows 11 via Windows Update. However, fresh/minimal Windows installs, Windows Server editions, and some enterprise locked-down environments may not have it. The app launches but shows a blank window or a "WebView2 not installed" error — no graceful fallback.

**Why it happens:** Developers test on their own Windows machines where WebView2 is already present and never test against a clean image.

**Consequences:**
- App silently fails on some user environments
- Support burden: users report "blank window" with no actionable error
- Enterprise users (exactly the self-hosted crowd) often run Windows Server

**Prevention:**
- Bundle the WebView2 bootstrapper in the installer or detect missing WebView2 at startup and prompt for install
- Wails provides a build option to embed the WebView2 bootstrapper (`-webview2 embed` or `-webview2 download`) — verify current flag name in Wails docs
- Test against Windows Server 2019/2022 in CI using a fresh image
- Display a human-readable error if WebView2 initialization fails

**Detection warning signs:**
- No Windows installer (just a bare `.exe`) and no WebView2 detection code
- Tests only run on developer's Windows 11 machine

**Phase mapping:** Distribution/packaging phase. Must be addressed before the first Windows release.

**Confidence:** MEDIUM (Wails WebView2 handling is documented; specific `-webview2` flag names should be verified against current Wails v2 docs before use)

---

### Pitfall 4: Linux Runtime Requires webkit2gtk — Not Universal Across Distros

**What goes wrong:** Wails on Linux requires `libgtk-3` and `webkit2gtk-4.0` (or `webkit2gtk-4.1` on newer distros). The exact package name and version differ by distro: Ubuntu 22.04 uses `webkit2gtk-4.0`, Ubuntu 24.04 may require `webkit2gtk-4.1`, Arch uses `webkit2gtk`, Alpine Linux has no webkit2gtk in its standard repos.

**Why it happens:** Linux ecosystem fragmentation. Developers test on Ubuntu LTS and ship a binary that works on their machine.

**Consequences:**
- App fails to launch on Arch, Fedora, or Alpine with a missing `.so` linker error
- Docker runtime image must explicitly install the correct webkit2gtk version for the base image chosen
- Alpine-based Docker images cannot easily run Wails apps (critical for the VNC/noVNC Docker runtime)

**Prevention:**
- For the Docker VNC runtime image: use `debian:bookworm-slim` or `ubuntu:22.04` as base, not Alpine
- Install `libwebkit2gtk-4.0-dev` (or `4.1` depending on image) explicitly in the Dockerfile
- In release notes/install docs, list the required system packages per distro
- For the AppImage or raw binary distribution: document runtime deps clearly; consider using `linuxdeploy` to bundle gtk/webkit libs

**Detection warning signs:**
- Dockerfile FROM uses `alpine`
- No `apt-get install libwebkit2gtk*` in the runtime Dockerfile
- No distro-specific install instructions in README

**Phase mapping:** Infrastructure/Docker phase (runtime image), and distribution phase (binary packaging).

**Confidence:** HIGH (webkit2gtk dependency is a well-known Wails Linux requirement; distro naming differences are documented in Wails issues)

---

### Pitfall 5: Cloud Save Sync Race Condition — Data Loss on Concurrent Session or Crash

**What goes wrong:** Cloud save upload on game exit races with the game process terminating. If the app crashes, loses network, or the user force-quits between "game exited" detection and upload completion, the save is not synced. Worse: on next launch, the app downloads the old cloud save over a newer local save that was never uploaded.

**Why it happens:** "Upload on exit, download on launch" sounds simple but the gap between game exit and upload completion is a data loss window. Most implementations detect game exit via process monitoring and immediately trigger upload — but if that trigger fails, there's no recovery.

**Consequences:**
- Save game data loss — this is catastrophic for users
- Trust destruction; users stop using cloud save feature
- Hard to debug because it only happens on crash/network loss paths

**Prevention:**
- Implement a local "dirty flag" (file or DB record) set when local save is modified, cleared only after confirmed upload
- On app launch, check if dirty flag is set — offer to upload before downloading cloud save
- Version saves with timestamps on both local and server; implement a conflict resolution UI ("Local save is newer — use local or cloud?")
- Never overwrite local with cloud without checking timestamps
- Implement upload retry with exponential backoff, not fire-and-forget
- Store the last-uploaded hash/timestamp to detect drift

**Detection warning signs:**
- Cloud save code has `uploadOnExit()` followed immediately by `downloadOnLaunch()` with no ordering guarantee
- No conflict detection — always overwrites in one direction
- No local dirty state tracking

**Phase mapping:** Cloud save feature phase. Design the state machine before writing any sync code.

**Confidence:** HIGH (this is a well-known class of sync bug; the specific failure modes are domain knowledge from game backup tools and cloud save implementations)

---

### Pitfall 6: Large File Download Reliability — No Resume Support = Re-download Everything

**What goes wrong:** Game files are 1-100GB. A plain `http.Get()` with no resume capability means any network interruption (VPN disconnect, sleep, ISP hiccup) restarts the entire download. Users on slower connections may never successfully download a 50GB game.

**Why it happens:** HTTP range requests are optional; developers implement the happy path first and defer reliability.

**Consequences:**
- Users report games "never finish downloading"
- Partial downloads consume disk space with no cleanup
- No progress persistence between app restarts

**Prevention:**
- Implement HTTP Range request support from the start: `Range: bytes=<offset>-`
- Check if the GameVault backend supports `Accept-Ranges: bytes` response header (verify against backend source)
- Store download progress to disk (byte offset + ETag/Last-Modified for validation)
- On resume, validate the partial file before appending (hash check of downloaded portion if feasible, or use ETag)
- Implement checksum verification (SHA256) post-download if GameVault backend provides checksums
- Use a download queue with pause/resume UI controls

**Detection warning signs:**
- Download implementation uses `io.Copy(file, resp.Body)` with no offset tracking
- No download progress persisted to disk
- No retry/resume on error

**Phase mapping:** Game download feature phase. Build resume support before shipping download functionality.

**Confidence:** HIGH (HTTP range requests are standard; game download resume is well-established pattern)

---

## Moderate Pitfalls

---

### Pitfall 7: VNC/noVNC in Docker — Display Sizing and DPI Misconfiguration

**What goes wrong:** The VNC server (typically Xvfb + x11vnc or TigerVNC) starts with a hardcoded resolution (e.g. 1024x768). The Wails app renders at that resolution. When the user opens noVNC in a browser, the viewport doesn't match, resulting in scrollbars or blurry scaling. High-DPI displays see a tiny app.

**Why it happens:** VNC setup tutorials use example resolutions. Nobody tests what the app looks like at different browser window sizes.

**Consequences:**
- Poor user experience — scrollbars in a browser-based desktop app feel broken
- App layout breaks at small resolutions (responsive CSS needed)
- Retina/HiDPI users see a very small, unscalable UI

**Prevention:**
- Set initial Xvfb resolution to 1920x1080 as the default (reasonable for most browser windows)
- Configure noVNC with `resize=scale` mode so it scales to fill the browser viewport
- Or use a VNC server that supports dynamic resolution (TigerVNC with `randr` extension, or use `xrandr` to resize on connect)
- Test at 1280x720, 1920x1080, and 2560x1440 virtual resolutions
- Wails app should use responsive CSS (flexbox/grid), not fixed pixel sizes

**Detection warning signs:**
- Xvfb command uses `-screen 0 800x600x24`
- noVNC URL has no `resize=scale` parameter
- Frontend CSS uses `width: 1200px` fixed values

**Phase mapping:** Docker runtime phase.

**Confidence:** MEDIUM (Xvfb/noVNC resolution handling is well-known; specific noVNC URL parameters should be verified against current noVNC docs)

---

### Pitfall 8: VNC Authentication — Exposing Unprotected VNC Port

**What goes wrong:** The Docker container exposes VNC port 5900 and noVNC port 6080 without authentication. If the user maps these to public-facing ports (or uses `--net=host`), anyone on the network can control the desktop session.

**Why it happens:** Development setups skip auth for convenience. It makes it into the shipped Dockerfile.

**Consequences:**
- Security vulnerability in a self-hosted tool (self-hosted users are the exact audience)
- Credential exposure — if user is logged into their GameVault server, an attacker can see/steal the session

**Prevention:**
- Enable VNC password authentication (even a simple password is better than none)
- Set password via environment variable (`VNC_PASSWORD`)
- Prefer noVNC with a reverse proxy (nginx) that enforces HTTPS and basic auth
- Document that VNC/noVNC ports should NOT be exposed to the internet without a VPN or auth proxy
- Consider: noVNC token-based auth for slightly better UX

**Detection warning signs:**
- VNC server started without `-passwd` or `-SecurityTypes None`
- Docker Compose exposes `5900:5900` and `6080:6080` with no auth config

**Phase mapping:** Docker runtime phase.

**Confidence:** HIGH (VNC auth omission is a well-known container security mistake)

---

### Pitfall 9: Credential Storage — Storing Server Passwords in Plaintext

**What goes wrong:** The app saves GameVault server passwords (possibly multiple servers) to a config file or SQLite DB in plaintext. On a shared machine, another user reads the config. On a compromised machine, credentials are trivially extracted.

**Why it happens:** Secure credential storage APIs are platform-specific and seem complex. Developers use JSON config files for simplicity.

**Consequences:**
- Password exposure for self-hosted GameVault servers (and potentially reused passwords)
- On Linux in Docker, the config volume is readable by root/host

**Prevention:**
- Use OS keychain/keyring: `golang.org/x/crypto` for encryption, or platform-specific:
  - Windows: DPAPI via `github.com/danieljoos/wincred`
  - macOS: Keychain via `github.com/keybase/go-keychain`
  - Linux: Secret Service API via `github.com/zalando/go-keyring` (works with GNOME Keyring, KWallet)
- `go-keyring` provides a cross-platform abstraction that works on all three platforms
- For Docker/headless environments where no keyring is available: fall back to AES-256-GCM encrypted file with a machine-derived key (not plaintext)
- Never log credentials, even at debug level

**Detection warning signs:**
- Config struct has `Password string` with `json:"password"` tag written to disk
- No mention of keyring in dependencies
- Log statements that include auth headers or tokens

**Phase mapping:** Authentication phase (early). Credential storage design affects the config schema for the lifetime of the project.

**Confidence:** MEDIUM (go-keyring library existence is known from training; specific API and current maintenance status should be verified)

---

### Pitfall 10: GameVault API Breaking Changes — No Versioning Defense

**What goes wrong:** The GameVault backend team releases a new version that changes a response shape (adds required fields, renames a field, changes an enum). The Go client's JSON deserialization panics or silently drops data. Features break without a clear error.

**Why it happens:** Tight coupling to specific response shapes. Using `encoding/json` with strict struct matching that fails on unexpected fields (or vice versa: omitting fields that become required).

**Consequences:**
- App breaks for users who upgraded their GameVault backend
- Hard to diagnose: app shows empty game list but no obvious error
- Need to track GameVault backend changelog and release new client builds for each breaking change

**Prevention:**
- Use `omitempty` and pointer types for all optional API response fields to handle missing fields gracefully
- Set `json.Decoder.DisallowUnknownFields(false)` (default) — unknown fields are silently ignored, which is the right behavior for forward compatibility
- Version-check the backend on connect: call the `/api/info` or health endpoint, compare semver, warn if backend is outside supported range
- Define supported backend version range in the client (`minBackendVersion`, `maxTestedBackendVersion`)
- Write API integration tests against the actual GameVault backend using Docker Compose

**Detection warning signs:**
- API structs use strict types with no pointer optionals
- No backend version check on connection
- API client has no error handling for 404 or schema mismatches beyond "failed to decode"

**Phase mapping:** API client phase (foundational). Version check should be in the first API integration.

**Confidence:** MEDIUM (JSON handling patterns are standard Go knowledge; GameVault-specific versioning behavior requires verification against their API docs)

---

### Pitfall 11: Wails Frontend Asset Embedding — Hot Reload vs Production Build Mismatch

**What goes wrong:** In development mode, Wails serves the frontend via a dev server (Vite/webpack dev server). In production, assets are embedded via Go's `embed.FS`. Paths, base URLs, and asset references that work in dev silently break in production because the embedded asset path handling differs from the dev server.

**Why it happens:** Developers spend most time in `wails dev` mode and only build for production late in the cycle.

**Consequences:**
- Production binary has broken asset paths (CSS not loading, images 404ing)
- Relative imports that worked in dev server fail in embedded FS
- `wails build` succeeds but the binary shows a broken UI

**Prevention:**
- Run `wails build` and test the output binary early (after the first UI screen is built, not after all features)
- Use absolute paths for assets relative to the Wails embedded FS root
- Configure Vite `base` correctly for Wails embedding (usually `base: './'` or `base: '/'` — verify against Wails/Vite integration docs)
- Add `wails build` to CI from the start so broken production builds are caught immediately
- Test the built binary, not just the dev server, in the CI smoke test

**Detection warning signs:**
- CI only runs `go test`, never `wails build`
- Asset imports use absolute dev server paths (`http://localhost:5173/assets/...`)
- First production build attempt is late in the project

**Phase mapping:** First phase that includes frontend work. Add a `wails build` CI step before any feature is considered "done."

**Confidence:** HIGH (Wails embed.FS vs dev server mismatch is a documented common issue in the Wails community)

---

### Pitfall 12: GitHub Actions CGo Cross-Compilation Matrix — Wrong Toolchain Per Platform

**What goes wrong:** The CI matrix tries to build all platforms on `ubuntu-latest`. Windows CGo requires mingw-w64. Linux arm64 requires a cross-compiler. macOS requires macOS runners. Using the wrong runner/toolchain per platform produces linker errors that are easy to fix but slow to discover if the matrix is built incorrectly from the start.

**Why it happens:** GitHub Actions matrix syntax makes it look like you can parameterize OS as just another variable. CGo breaks this assumption.

**Consequences:**
- CI pipeline fails for cross-compiled targets
- Minutes-long CI runs just to see an obvious linker error
- Temptation to disable CGo (which breaks Wails) to "fix" the CI

**Prevention:**
Use this matrix mapping (not a generic OS loop):

```yaml
# Correct approach
jobs:
  build:
    strategy:
      matrix:
        include:
          - os: ubuntu-latest
            target: linux/amd64
            cc: gcc
          - os: ubuntu-latest
            target: linux/arm64
            cc: aarch64-linux-gnu-gcc
            packages: gcc-aarch64-linux-gnu
          - os: ubuntu-latest
            target: windows/amd64
            cc: x86_64-w64-mingw32-gcc
            packages: gcc-mingw-w64-x86-64
          - os: macos-latest
            target: darwin/amd64
          - os: macos-latest  # or macos-14 for M1
            target: darwin/arm64
```

- Install cross-compiler packages in the build step before running `wails build`
- Set `CC` environment variable per target
- Do not use `CGO_ENABLED=0` for any Wails target

**Detection warning signs:**
- Matrix has `os: [ubuntu-latest, windows-latest, macos-latest]` without target-specific toolchain config
- Single `ubuntu-latest` runner for all targets without cross-compiler install steps
- `CGO_ENABLED: 0` in CI environment variables

**Phase mapping:** CI/CD pipeline phase.

**Confidence:** MEDIUM (mingw-w64 pattern for Windows CGo is well-known; exact GitHub Actions syntax should be verified against current actions/setup-go docs)

---

## Minor Pitfalls

---

### Pitfall 13: Wails Context Cancellation — Goroutine Leaks on Window Close

**What goes wrong:** Background goroutines (polling for downloads, VNC health checks, sync status) are started but not tied to the Wails app context. When the window closes or the app exits, these goroutines keep running until the process terminates — or worse, panic on closed channels.

**Why it happens:** Go goroutines are "fire and forget" by default. The Wails lifecycle (OnStartup, OnShutdown) is easy to miss.

**Prevention:**
- All background goroutines must receive a `context.Context` derived from the Wails app context
- Cancel the context in `OnShutdown`
- Use `sync.WaitGroup` to drain goroutines before process exit
- Download manager, sync workers, and API polling loops must all be context-aware

**Phase mapping:** Any phase introducing background workers.

**Confidence:** HIGH (standard Go goroutine lifecycle; Wails OnStartup/OnShutdown lifecycle is documented)

---

### Pitfall 14: Process Monitoring for Game Launch — Platform-Specific Behavior

**What goes wrong:** Detecting "game exited" to trigger cloud save upload requires process monitoring. `os/exec` `Wait()` works for processes the app directly launched. But if the game is launched via a launcher (Steam, Lutris, Wine), the child process may exit while the game continues running, making `Wait()` return prematurely — triggering a save upload while the game is still running.

**Prevention:**
- For directly-launched executables: use `cmd.Wait()` — reliable
- For launcher-wrapped games: poll for the game binary by name using OS process listing, not just child process wait
- Document that "launched via another launcher" is a known limitation for save detection
- Provide a manual "game is done" button as fallback

**Phase mapping:** Game launch/run phase.

**Confidence:** MEDIUM (process monitoring complexity is known; GameVault-specific launch flow needs verification)

---

### Pitfall 15: Disk Space Management — No Cleanup on Failed/Cancelled Downloads

**What goes wrong:** A 40GB game download starts, the user cancels or the download fails. The partial file remains on disk. If this happens repeatedly, the user's disk fills up with orphaned partials.

**Prevention:**
- Track all in-progress download files in a manifest
- On app startup, check for orphaned partials (files in download dir not in the manifest or with incomplete marker)
- Provide a "Clean up incomplete downloads" UI action
- Partial files should use a `.part` extension and only be renamed on completion

**Phase mapping:** Game download phase.

**Confidence:** HIGH (partial file cleanup is standard practice; `.part` extension pattern is used by Firefox, wget, etc.)

---

### Pitfall 16: Wails App Not Respecting System Theme (Dark/Light Mode)

**What goes wrong:** The frontend CSS hardcodes colors rather than using `prefers-color-scheme` media queries. On macOS and Windows 11, users expect apps to follow system dark/light mode. The app always appears in one mode.

**Prevention:**
- Use CSS `prefers-color-scheme` from the start
- Or use a theme context/state that initializes from Wails' `runtime.WindowGetSystemTheme()` (verify API name in Wails docs)
- Design the color system with CSS custom properties (variables) so theme switching is a single variable change

**Phase mapping:** UI/design phase.

**Confidence:** MEDIUM (Wails theme API existence is likely but specific function name needs verification)

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Docker build pipeline | macOS cannot be cross-compiled from Linux | Use native macOS GitHub Actions runners for darwin targets |
| Docker build pipeline | CGo requires platform-specific C compilers | Install mingw-w64 (Windows), aarch64-gcc (arm64) per target in build stage |
| Docker VNC runtime | Alpine base image lacks webkit2gtk | Use debian:bookworm-slim or ubuntu:22.04 as runtime base |
| Docker VNC runtime | Unprotected VNC port | Set VNC_PASSWORD env var; document not to expose ports publicly |
| Docker VNC runtime | Fixed low resolution | Use 1920x1080 Xvfb; configure noVNC resize=scale |
| Authentication | Plaintext password storage | Use go-keyring or encrypted file fallback from the start |
| API client | Tight coupling to response shapes | Use pointer fields, omitempty, backend version check |
| Game download | No resume support | Implement HTTP Range requests before shipping download feature |
| Cloud saves | Sync race conditions | Implement dirty flag + conflict resolution before any upload/download |
| CI/CD | Wrong runner per CGo target | Use target-specific matrix entries, not generic OS array |
| Frontend | dev vs production asset path mismatch | Run wails build in CI from the first UI commit |
| Any background work | Goroutine leaks on shutdown | Tie all workers to Wails app context with cancellation |

---

## Sources

- Wails v2 documentation (training knowledge, cutoff August 2025) — verify at https://wails.io/docs/
- Apple macOS SDK redistribution restrictions — well-established; verify at https://developer.apple.com/support/xcode/
- GameVault backend API — verify at https://github.com/Phalcode/gamevault-backend and https://gamevault.app/docs
- go-keyring library — verify current status at https://github.com/zalando/go-keyring
- noVNC configuration — verify at https://github.com/novnc/noVNC
- HTTP range requests — RFC 7233 (stable standard)

**Overall confidence: MEDIUM** — Core pitfalls (CGo, macOS SDK, cloud save races, large file resume, credential storage, VNC auth) are HIGH confidence from well-established domain knowledge. Specific Wails API names, noVNC parameter syntax, and GitHub Actions specifics are MEDIUM and should be verified against current docs before implementation.
