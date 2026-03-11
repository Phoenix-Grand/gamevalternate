<script lang="ts">
  // Wails injects window.runtime globally — no import needed
  const runtime = () => (window as any).runtime;

  function minimise() { runtime()?.WindowMinimise(); }
  function toggleMaximise() { runtime()?.WindowToggleMaximise(); }
  function quit() { runtime()?.Quit(); }
</script>

<div class="titlebar">
  <div class="titlebar-drag">
    <span class="titlebar-title">GameVault</span>
  </div>
  <div class="titlebar-controls">
    <button class="control" onclick={minimise} aria-label="Minimize">
      <svg width="10" height="1" viewBox="0 0 10 1"><rect width="10" height="1" fill="currentColor"/></svg>
    </button>
    <button class="control" onclick={toggleMaximise} aria-label="Maximize">
      <svg width="10" height="10" viewBox="0 0 10 10"><rect x="0.5" y="0.5" width="9" height="9" fill="none" stroke="currentColor"/></svg>
    </button>
    <button class="control close" onclick={quit} aria-label="Close">
      <svg width="10" height="10" viewBox="0 0 10 10">
        <line x1="0" y1="0" x2="10" y2="10" stroke="currentColor" stroke-width="1.2"/>
        <line x1="10" y1="0" x2="0" y2="10" stroke="currentColor" stroke-width="1.2"/>
      </svg>
    </button>
  </div>
</div>

<style>
  .titlebar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 32px;
    background-color: var(--sidebar);
    border-bottom: 1px solid var(--sidebar-border);
    flex-shrink: 0;
    user-select: none;
  }

  .titlebar-drag {
    flex: 1;
    display: flex;
    align-items: center;
    padding-left: 12px;
    height: 100%;
    /* stylelint-disable-next-line property-no-unknown */
    -webkit-app-region: drag;
  }

  .titlebar-title {
    font-size: 12px;
    font-weight: 500;
    color: var(--sidebar-foreground);
    opacity: 0.7;
  }

  .titlebar-controls {
    display: flex;
    height: 100%;
    /* stylelint-disable-next-line property-no-unknown */
    -webkit-app-region: no-drag;
  }

  .control {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 46px;
    height: 100%;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--sidebar-foreground);
    opacity: 0.7;
    transition: background-color 0.1s, opacity 0.1s;
  }

  .control:hover {
    background-color: var(--sidebar-accent);
    opacity: 1;
  }

  .control.close:hover {
    background-color: hsl(0 72.2% 50.6%);
    color: white;
    opacity: 1;
  }
</style>
