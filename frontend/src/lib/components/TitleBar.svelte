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
      <svg width="12" height="2" viewBox="0 0 12 2"><rect width="12" height="2" fill="currentColor"/></svg>
    </button>
    <button class="control" onclick={toggleMaximise} aria-label="Maximize">
      <svg width="11" height="11" viewBox="0 0 11 11"><rect x="0.5" y="0.5" width="10" height="10" fill="none" stroke="currentColor" stroke-width="1.2"/></svg>
    </button>
    <button class="control close" onclick={quit} aria-label="Close">
      <svg width="11" height="11" viewBox="0 0 11 11">
        <line x1="1" y1="1" x2="10" y2="10" stroke="currentColor" stroke-width="1.4"/>
        <line x1="10" y1="1" x2="1" y2="10" stroke="currentColor" stroke-width="1.4"/>
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
    background-color: hsl(240 10% 12%);
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
    font-weight: 600;
    color: hsl(0 0% 85%);
    letter-spacing: 0.02em;
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
    color: hsl(0 0% 70%);
    transition: background-color 0.1s, color 0.1s;
  }

  .control:hover {
    background-color: hsl(240 5% 22%);
    color: hsl(0 0% 95%);
  }

  .control.close:hover {
    background-color: hsl(0 72% 51%);
    color: white;
  }
</style>
