package app

import (
"context"
"fmt"
"io"
"net/http"
"os"
"path/filepath"
"strings"
"sync"
"time"

"github.com/Advik-B/megobasterd/internal/api"
"github.com/Advik-B/megobasterd/internal/config"
"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
ctx      context.Context
config   *config.Config
client   *api.MegaClient
downloads map[string]*Download
mu       sync.RWMutex
}

// Download represents an active download
type Download struct {
ID          string    `json:"id"`
URL         string    `json:"url"`
FileName    string    `json:"fileName"`
FilePath    string    `json:"filePath"`
FileSize    int64     `json:"fileSize"`
Downloaded  int64     `json:"downloaded"`
Speed       float64   `json:"speed"`
Status      string    `json:"status"`
Progress    float64   `json:"progress"`
StartTime   time.Time `json:"startTime"`
Error       string    `json:"error,omitempty"`
cancelFunc  context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
cfg, err := config.Load()
if err != nil {
cfg = config.GetDefault()
}

return &App{
config:    cfg,
client:    api.NewMegaClient(),
downloads: make(map[string]*Download),
}
}

// Startup is called when the app starts
func (a *App) Startup(ctx context.Context) {
a.ctx = ctx
runtime.LogInfo(ctx, "MegaBasterd started")
}

// ParseMegaURL parses a MEGA URL and extracts file ID and key
func (a *App) ParseMegaURL(url string) (fileID, key string, err error) {
if !strings.Contains(url, "mega.nz") {
return "", "", fmt.Errorf("not a MEGA URL")
}

parts := strings.Split(url, "/")
if len(parts) < 2 {
return "", "", fmt.Errorf("invalid MEGA URL format")
}

lastPart := parts[len(parts)-1]
idAndKey := strings.Split(lastPart, "#")
if len(idAndKey) != 2 {
return "", "", fmt.Errorf("invalid MEGA URL format: missing key")
}

return idAndKey[0], idAndKey[1], nil
}

// AddDownload adds a new download from a MEGA URL
func (a *App) AddDownload(url string) (*Download, error) {
a.mu.Lock()
defer a.mu.Unlock()

// Parse MEGA URL
fileID, key, err := a.ParseMegaURL(url)
if err != nil {
return nil, err
}

// Get download URL from MEGA API
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

downloadURL, err := a.client.GetDownloadURL(ctx, fileID, key)
if err != nil {
return nil, fmt.Errorf("failed to get download URL: %w", err)
}

// Extract filename from URL or use file ID
fileName := fileID + ".bin"

// Create download struct
download := &Download{
ID:        fileID,
URL:       downloadURL,
FileName:  fileName,
FilePath:  filepath.Join(a.config.DownloadPath, fileName),
Status:    "queued",
StartTime: time.Now(),
}

a.downloads[fileID] = download

// Start download in background
go a.startDownload(download)

runtime.LogInfo(a.ctx, fmt.Sprintf("Download added: %s", fileName))
return download, nil
}

// startDownload performs the actual download
func (a *App) startDownload(download *Download) {
ctx, cancel := context.WithCancel(context.Background())
download.cancelFunc = cancel

// Update status
a.updateDownloadStatus(download.ID, "downloading")

// Make HTTP request
req, err := http.NewRequestWithContext(ctx, "GET", download.URL, nil)
if err != nil {
a.updateDownloadError(download.ID, err.Error())
return
}

client := &http.Client{
Timeout: 0, // No timeout for large downloads
}

resp, err := client.Do(req)
if err != nil {
a.updateDownloadError(download.ID, err.Error())
return
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
a.updateDownloadError(download.ID, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status))
return
}

// Get file size
download.FileSize = resp.ContentLength

// Create output file
out, err := os.Create(download.FilePath)
if err != nil {
a.updateDownloadError(download.ID, err.Error())
return
}
defer out.Close()

// Download with progress tracking
buffer := make([]byte, 32*1024) // 32KB buffer
var downloaded int64
lastUpdate := time.Now()

for {
n, err := resp.Body.Read(buffer)
if n > 0 {
written, writeErr := out.Write(buffer[:n])
if writeErr != nil {
a.updateDownloadError(download.ID, writeErr.Error())
return
}
downloaded += int64(written)

// Update progress every 500ms
if time.Since(lastUpdate) > 500*time.Millisecond {
a.updateDownloadProgress(download.ID, downloaded, download.FileSize)
lastUpdate = time.Now()
}
}

if err == io.EOF {
break
}
if err != nil {
a.updateDownloadError(download.ID, err.Error())
return
}
}

// Final update
a.updateDownloadProgress(download.ID, downloaded, download.FileSize)
a.updateDownloadStatus(download.ID, "completed")
runtime.LogInfo(a.ctx, fmt.Sprintf("Download completed: %s", download.FileName))
}

// updateDownloadStatus updates the download status
func (a *App) updateDownloadStatus(id, status string) {
a.mu.Lock()
defer a.mu.Unlock()

if download, exists := a.downloads[id]; exists {
download.Status = status
runtime.EventsEmit(a.ctx, "download-update", download)
}
}

// updateDownloadProgress updates download progress
func (a *App) updateDownloadProgress(id string, downloaded, total int64) {
a.mu.Lock()
defer a.mu.Unlock()

if download, exists := a.downloads[id]; exists {
download.Downloaded = downloaded
if total > 0 {
download.Progress = float64(downloaded) / float64(total) * 100
}

// Calculate speed
elapsed := time.Since(download.StartTime).Seconds()
if elapsed > 0 {
download.Speed = float64(downloaded) / elapsed
}

runtime.EventsEmit(a.ctx, "download-update", download)
}
}

// updateDownloadError sets download error status
func (a *App) updateDownloadError(id, errorMsg string) {
a.mu.Lock()
defer a.mu.Unlock()

if download, exists := a.downloads[id]; exists {
download.Status = "failed"
download.Error = errorMsg
runtime.EventsEmit(a.ctx, "download-update", download)
runtime.LogError(a.ctx, fmt.Sprintf("Download error for %s: %s", download.FileName, errorMsg))
}
}

// GetDownloads returns all downloads
func (a *App) GetDownloads() []*Download {
a.mu.RLock()
defer a.mu.RUnlock()

downloads := make([]*Download, 0, len(a.downloads))
for _, d := range a.downloads {
downloads = append(downloads, d)
}
return downloads
}

// PauseDownload pauses a download
func (a *App) PauseDownload(id string) error {
a.mu.Lock()
defer a.mu.Unlock()

download, exists := a.downloads[id]
if !exists {
return fmt.Errorf("download not found")
}

if download.cancelFunc != nil {
download.cancelFunc()
}
download.Status = "paused"
runtime.EventsEmit(a.ctx, "download-update", download)
return nil
}

// RemoveDownload removes a download
func (a *App) RemoveDownload(id string) error {
a.mu.Lock()
defer a.mu.Unlock()

download, exists := a.downloads[id]
if !exists {
return fmt.Errorf("download not found")
}

if download.cancelFunc != nil {
download.cancelFunc()
}

delete(a.downloads, id)
runtime.EventsEmit(a.ctx, "download-removed", id)
return nil
}

// GetConfig returns current configuration
func (a *App) GetConfig() *config.Config {
return a.config
}
