<script lang="ts">
  import Sidebar from '$lib/components/Sidebar.svelte';
  import ConnectCard from '$lib/components/ConnectCard.svelte';
  import TitleBar from '$lib/components/TitleBar.svelte';

  // Phase 1: always false — Phase 2 wires real auth state
  let isConnected = $state(false);
</script>

<div class="root">
  <TitleBar />
  <div class="app-shell">
    <Sidebar />
    <main class="content-area">
      {#if !isConnected}
        <ConnectCard />
      {/if}
    </main>
  </div>
</div>

<style>
  .root {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
  }

  .app-shell {
    display: flex;
    flex: 1;
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
