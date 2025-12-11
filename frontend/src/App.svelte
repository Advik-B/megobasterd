<script>
  import { onMount } from 'svelte';
  import { AddDownload, GetDownloads, PauseDownload, RemoveDownload } from '../wailsjs/go/app/App';
  import { EventsOn } from '../wailsjs/runtime/runtime';

  let downloads = [];
  let url = '';
  let error = '';

  onMount(async () => {
    // Load initial downloads
    downloads = await GetDownloads();

    // Listen for download updates
    EventsOn('download-update', (download) => {
      const index = downloads.findIndex(d => d.id === download.id);
      if (index >= 0) {
        downloads[index] = download;
        downloads = [...downloads];
      } else {
        downloads = [...downloads, download];
      }
    });

    EventsOn('download-removed', (id) => {
      downloads = downloads.filter(d => d.id !== id);
    });
  });

  async function handleAddDownload(e) {
    e.preventDefault();
    error = '';

    if (!url.trim()) {
      error = 'Please enter a MEGA URL';
      return;
    }

    try {
      const download = await AddDownload(url);
      url = '';
      downloads = [...downloads, download];
    } catch (err) {
      error = err.message || 'Failed to add download';
    }
  }

  async function handlePause(id) {
    try {
      await PauseDownload(id);
    } catch (err) {
      console.error('Failed to pause download:', err);
    }
  }

  async function handleRemove(id) {
    try {
      await RemoveDownload(id);
    } catch (err) {
      console.error('Failed to remove download:', err);
    }
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i];
  }

  function formatSpeed(bytesPerSec) {
    return formatBytes(bytesPerSec) + '/s';
  }
</script>

<div class="app">
  <header class="header">
    <h1>🚀 MegaBasterd - Go Edition</h1>
    <p class="subtitle">Wails + Svelte powered MEGA downloader</p>
  </header>

  <div class="add-download">
    <form on:submit={handleAddDownload}>
      <input
        type="text"
        bind:value={url}
        placeholder="Enter MEGA URL (e.g., https://mega.nz/file/...)"
        class="url-input"
      />
      <button type="submit" class="btn btn-primary">
        Add Download
      </button>
    </form>
    {#if error}
      <div class="error">{error}</div>
    {/if}
  </div>

  <div class="downloads">
    <h2>Downloads</h2>
    {#if downloads.length === 0}
      <div class="empty-state">
        <p>No downloads yet. Add a MEGA URL to get started!</p>
      </div>
    {:else}
      <div class="downloads-list">
        {#each downloads as download (download.id)}
          <div class="download-item {download.status}">
            <div class="download-header">
              <div class="download-name">{download.fileName}</div>
              <div class="download-actions">
                {#if download.status === 'downloading'}
                  <button
                    on:click={() => handlePause(download.id)}
                    class="btn btn-small"
                  >
                    Pause
                  </button>
                {/if}
                <button
                  on:click={() => handleRemove(download.id)}
                  class="btn btn-small btn-danger"
                >
                  Remove
                </button>
              </div>
            </div>

            <div class="download-progress">
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  style="width: {download.progress}%"
                ></div>
                <div class="progress-text">
                  {download.progress.toFixed(1)}%
                </div>
              </div>
            </div>

            <div class="download-info">
              <span class="info-item">
                Size: {formatBytes(download.fileSize)}
              </span>
              <span class="info-item">
                Downloaded: {formatBytes(download.downloaded)}
              </span>
              <span class="info-item">
                Speed: {formatSpeed(download.speed)}
              </span>
              <span class="status-badge status-{download.status}">
                {download.status}
              </span>
            </div>

            {#if download.error}
              <div class="download-error">
                Error: {download.error}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
      'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
      sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  .app {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
  }

  .header {
    text-align: center;
    color: white;
    margin-bottom: 30px;
  }

  .header h1 {
    font-size: 2.5rem;
    margin-bottom: 10px;
  }

  .subtitle {
    font-size: 1.1rem;
    opacity: 0.9;
  }

  .add-download {
    background: white;
    border-radius: 12px;
    padding: 25px;
    margin-bottom: 30px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  }

  form {
    display: flex;
    gap: 10px;
  }

  .url-input {
    flex: 1;
    padding: 12px 18px;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    font-size: 1rem;
    transition: border-color 0.3s;
  }

  .url-input:focus {
    outline: none;
    border-color: #667eea;
  }

  .btn {
    padding: 12px 24px;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
  }

  .btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
  }

  .btn-small {
    padding: 6px 12px;
    font-size: 0.875rem;
  }

  .btn-danger {
    background: #e74c3c;
    color: white;
  }

  .btn-danger:hover {
    background: #c0392b;
  }

  .error {
    margin-top: 10px;
    padding: 10px;
    background: #fee;
    color: #c00;
    border-radius: 6px;
    border-left: 4px solid #c00;
  }

  .downloads {
    background: white;
    border-radius: 12px;
    padding: 25px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  }

  .downloads h2 {
    margin-bottom: 20px;
    color: #333;
  }

  .empty-state {
    text-align: center;
    padding: 40px;
    color: #999;
  }

  .downloads-list {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .download-item {
    border: 2px solid #e0e0e0;
    border-radius: 10px;
    padding: 15px;
    transition: all 0.3s;
  }

  .download-item:hover {
    border-color: #667eea;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.1);
  }

  .download-item.completed {
    background: #f0fff4;
    border-color: #38a169;
  }

  .download-item.failed {
    background: #fff5f5;
    border-color: #e53e3e;
  }

  .download-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .download-name {
    font-weight: 600;
    font-size: 1.1rem;
    color: #333;
  }

  .download-actions {
    display: flex;
    gap: 8px;
  }

  .download-progress {
    margin-bottom: 12px;
    position: relative;
  }

  .progress-bar {
    height: 24px;
    background: #e0e0e0;
    border-radius: 12px;
    overflow: hidden;
    position: relative;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    transition: width 0.3s ease;
  }

  .progress-text {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-weight: 600;
    color: #333;
    font-size: 0.875rem;
  }

  .download-info {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
    align-items: center;
    font-size: 0.9rem;
    color: #666;
  }

  .info-item {
    display: inline-block;
  }

  .status-badge {
    display: inline-block;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 0.875rem;
    font-weight: 600;
    text-transform: uppercase;
  }

  .status-queued {
    background: #e3f2fd;
    color: #1976d2;
  }

  .status-downloading {
    background: #fff3cd;
    color: #856404;
  }

  .status-completed {
    background: #d4edda;
    color: #155724;
  }

  .status-failed {
    background: #f8d7da;
    color: #721c24;
  }

  .status-paused {
    background: #e2e3e5;
    color: #383d41;
  }

  .download-error {
    margin-top: 10px;
    padding: 10px;
    background: #fff5f5;
    color: #e53e3e;
    border-radius: 6px;
    border-left: 4px solid #e53e3e;
    font-size: 0.875rem;
  }
</style>
