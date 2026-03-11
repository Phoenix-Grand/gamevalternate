<script lang="ts">
  import Sidebar from '$lib/components/Sidebar.svelte';
  import ConnectCard from '$lib/components/ConnectCard.svelte';

  // Phase 1: always false — Phase 2 wires real auth state
  let isConnected = $state(false);
</script>

<div class="app-shell">
  <Sidebar />
  <main class="content-area">
    {#if !isConnected}
      <ConnectCard />
    {/if}
  </main>
</div>

<style>
  .app-shell {
    display: flex;
    height: 100vh;
    overflow: hidden;
    background-color: var(--background);
    color: var(--foreground);
  }
  .content-area {
    flex: 1;
    overflow: auto;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
