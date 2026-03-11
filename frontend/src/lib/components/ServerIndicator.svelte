<script lang="ts">
  // Phase 1: static — Phase 2 replaces with real server/account state
  // Phase 2+ will pass server URL and username via props or a shared $state store
</script>

<button class="server-indicator" title="No server connected — click to connect">
  <div class="indicator-dot disconnected"></div>
  <div class="indicator-text">
    <span class="label">No server</span>
    <span class="sublabel">Click to connect</span>
  </div>
</button>

<style>
  .server-indicator {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 8px 12px;
    border-radius: 6px;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
  }

  .server-indicator:hover {
    background-color: var(--accent);
  }

  .indicator-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  .disconnected {
    background-color: var(--muted-foreground);
  }

  .indicator-text {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .label {
    font-size: 13px;
    font-weight: 500;
    color: var(--foreground);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .sublabel {
    font-size: 11px;
    color: var(--muted-foreground);
  }
</style>
