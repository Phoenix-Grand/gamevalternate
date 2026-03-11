<script lang="ts">
  import ServerIndicator from './ServerIndicator.svelte';

  type Route = 'library' | 'downloads' | 'settings';
  let activeRoute = $state<Route>('library');
</script>

<nav class="sidebar">
  <div class="nav-items">
    <button
      class="nav-item"
      class:active={activeRoute === 'library'}
      onclick={() => activeRoute = 'library'}
      title="Library"
    >
      <!-- Library icon: grid -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <rect x="2" y="3" width="7" height="7"/><rect x="15" y="3" width="7" height="7"/>
        <rect x="2" y="14" width="7" height="7"/><rect x="15" y="14" width="7" height="7"/>
      </svg>
      <span>Library</span>
    </button>

    <button
      class="nav-item"
      class:active={activeRoute === 'downloads'}
      onclick={() => activeRoute = 'downloads'}
      title="Downloads"
    >
      <!-- Download icon: arrow pointing down -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
      </svg>
      <span>Downloads</span>
    </button>

    <button
      class="nav-item"
      class:active={activeRoute === 'settings'}
      onclick={() => activeRoute = 'settings'}
      title="Settings"
    >
      <!-- Settings icon: gear -->
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="3"/>
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
      </svg>
      <span>Settings</span>
    </button>
  </div>

  <div class="server-section">
    <ServerIndicator />
  </div>
</nav>

<style>
  .sidebar {
    width: 200px;
    min-width: 200px;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 12px 8px;
    border-right: 1px solid var(--border);
    background-color: var(--card);
  }

  .nav-items {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    border-radius: 6px;
    border: none;
    background: transparent;
    color: var(--muted-foreground);
    cursor: pointer;
    text-align: left;
    font-size: 14px;
    width: 100%;
    transition: background-color 0.15s, color 0.15s;
  }

  .nav-item:hover {
    background-color: var(--accent);
    color: var(--accent-foreground);
  }

  .nav-item.active {
    background-color: var(--accent);
    color: var(--foreground);
    font-weight: 500;
  }

  .server-section {
    border-top: 1px solid var(--border);
    padding-top: 8px;
  }
</style>
