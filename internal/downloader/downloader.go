package downloader

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrDownloadCanceled is returned when download is canceled
	ErrDownloadCanceled = errors.New("download canceled")
	
	// ErrInvalidURL is returned when URL is invalid
	ErrInvalidURL = errors.New("invalid URL")
)

// Status represents download status
type Status string

const (
	StatusQueued     Status = "Queued"
	StatusDownloading Status = "Downloading"
	StatusPaused     Status = "Paused"
	StatusCompleted  Status = "Completed"
	StatusFailed     Status = "Failed"
)

// Download represents a file download
type Download struct {
	ID          string
	URL         string
	FilePath    string
	FileName    string
	FileSize    int64
	Workers     int
	Status      Status
	Progress    int64
	Speed       float64
	
	chunks      []*Chunk
	mu          sync.RWMutex
	cancelFunc  context.CancelFunc
}

// Chunk represents a download chunk
type Chunk struct {
	ID          int
	Start       int64
	End         int64
	Downloaded  int64
	Status      Status
}

// NewDownload creates a new download
func NewDownload(url, filePath string, workers int) *Download {
	return &Download{
		URL:      url,
		FilePath: filePath,
		Workers:  workers,
		Status:   StatusQueued,
		chunks:   make([]*Chunk, 0),
	}
}

// Start starts the download
func (d *Download) Start(ctx context.Context) error {
	d.mu.Lock()
	d.Status = StatusDownloading
	d.mu.Unlock()
	
	// Create cancelable context
	ctx, cancel := context.WithCancel(ctx)
	d.cancelFunc = cancel
	
	// Initialize chunks
	if err := d.initChunks(); err != nil {
		d.setStatus(StatusFailed)
		return err
	}
	
	// Download chunks in parallel
	g, ctx := errgroup.WithContext(ctx)
	
	// Limit concurrency with semaphore
	sem := make(chan struct{}, d.Workers)
	
	for _, chunk := range d.chunks {
		chunk := chunk // Capture loop variable
		
		g.Go(func() error {
			select {
			case sem <- struct{}{}: // Acquire
				defer func() { <-sem }() // Release
				return d.downloadChunk(ctx, chunk)
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}
	
	// Wait for all chunks to complete
	if err := g.Wait(); err != nil {
		d.setStatus(StatusFailed)
		return err
	}
	
	// Merge chunks
	if err := d.mergeChunks(); err != nil {
		d.setStatus(StatusFailed)
		return err
	}
	
	d.setStatus(StatusCompleted)
	return nil
}

// Pause pauses the download
func (d *Download) Pause() {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	if d.cancelFunc != nil {
		d.cancelFunc()
	}
	d.Status = StatusPaused
}

// Resume resumes the download
func (d *Download) Resume(ctx context.Context) error {
	d.mu.RLock()
	if d.Status != StatusPaused {
		d.mu.RUnlock()
		return errors.New("download is not paused")
	}
	d.mu.RUnlock()
	
	return d.Start(ctx)
}

// GetProgress returns download progress (0-100)
func (d *Download) GetProgress() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	if d.FileSize == 0 {
		return 0
	}
	
	return float64(d.Progress) / float64(d.FileSize) * 100
}

// GetSpeed returns current download speed in bytes/sec
func (d *Download) GetSpeed() float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Speed
}

// setStatus sets the download status
func (d *Download) setStatus(status Status) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Status = status
}

// initChunks initializes download chunks
func (d *Download) initChunks() error {
	// TODO: Implement chunk initialization
	// This would divide the file into chunks based on file size
	// For now, return a placeholder
	return errors.New("chunk initialization not yet implemented")
}

// downloadChunk downloads a single chunk
func (d *Download) downloadChunk(ctx context.Context, chunk *Chunk) error {
	// TODO: Implement actual chunk download
	// This would:
	// 1. Make HTTP request with Range header
	// 2. Download chunk data
	// 3. Write to temporary file
	// 4. Update progress
	
	return errors.New("chunk download not yet implemented")
}

// mergeChunks merges all chunks into final file
func (d *Download) mergeChunks() error {
	// TODO: Implement chunk merging
	// This would combine all chunk files into the final file
	return errors.New("chunk merging not yet implemented")
}

// DownloadManager manages multiple downloads
type DownloadManager struct {
	downloads map[string]*Download
	mu        sync.RWMutex
}

// NewDownloadManager creates a new download manager
func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		downloads: make(map[string]*Download),
	}
}

// AddDownload adds a download to the manager
func (dm *DownloadManager) AddDownload(download *Download) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	if _, exists := dm.downloads[download.ID]; exists {
		return fmt.Errorf("download %s already exists", download.ID)
	}
	
	dm.downloads[download.ID] = download
	return nil
}

// RemoveDownload removes a download from the manager
func (dm *DownloadManager) RemoveDownload(id string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	
	if download, exists := dm.downloads[id]; exists {
		download.Pause() // Stop if running
		delete(dm.downloads, id)
		return nil
	}
	
	return fmt.Errorf("download %s not found", id)
}

// GetDownload retrieves a download by ID
func (dm *DownloadManager) GetDownload(id string) (*Download, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	
	download, exists := dm.downloads[id]
	if !exists {
		return nil, fmt.Errorf("download %s not found", id)
	}
	
	return download, nil
}

// ListDownloads returns all downloads
func (dm *DownloadManager) ListDownloads() []*Download {
	dm.mu.RLock()
	defer dm.mu.RUnlock()
	
	downloads := make([]*Download, 0, len(dm.downloads))
	for _, d := range dm.downloads {
		downloads = append(downloads, d)
	}
	
	return downloads
}
