<script lang="ts">
  // Phase 1: inert — Connect button does nothing
  // Phase 2 will wire this to the real auth backend via wailsjs bindings
  let serverURL = $state('');
  let canConnect = $derived(serverURL.trim().length > 0);
</script>

<div class="connect-card">
  <div class="card-header">
    <h1 class="card-title">GameVault</h1>
    <p class="card-subtitle">Connect to a GameVault server</p>
  </div>

  <div class="card-form">
    <label for="server-url" class="form-label">Server URL</label>
    <input
      id="server-url"
      type="url"
      class="form-input"
      placeholder="https://your-gamevault-server.example.com"
      bind:value={serverURL}
    />
    <button
      class="connect-button"
      disabled={!canConnect}
      onclick={() => {
        // Phase 2 wires real connection logic here
        console.log('Connect to:', serverURL);
      }}
    >
      Connect
    </button>
  </div>
</div>

<style>
  .connect-card {
    background-color: var(--card);
    border: 1px solid var(--border);
    border-radius: 12px;
    padding: 40px;
    width: 100%;
    max-width: 420px;
    display: flex;
    flex-direction: column;
    gap: 28px;
  }

  .card-header {
    text-align: center;
  }

  .card-title {
    font-size: 28px;
    font-weight: 700;
    color: var(--foreground);
    margin: 0 0 8px 0;
  }

  .card-subtitle {
    font-size: 14px;
    color: var(--muted-foreground);
    margin: 0;
  }

  .card-form {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .form-label {
    font-size: 13px;
    font-weight: 500;
    color: var(--foreground);
  }

  .form-input {
    width: 100%;
    padding: 10px 14px;
    border-radius: 6px;
    border: 1px solid var(--border);
    background-color: var(--background);
    color: var(--foreground);
    font-size: 14px;
    box-sizing: border-box;
    outline: none;
  }

  .form-input:focus {
    border-color: var(--ring);
    box-shadow: 0 0 0 2px color-mix(in srgb, var(--ring) 20%, transparent);
  }

  .connect-button {
    width: 100%;
    padding: 10px 16px;
    border-radius: 6px;
    border: none;
    background-color: var(--primary);
    color: var(--primary-foreground);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.15s;
  }

  .connect-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .connect-button:not(:disabled):hover {
    opacity: 0.9;
  }
</style>
