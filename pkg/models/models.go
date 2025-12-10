package models

import "time"

// Transfer represents a file transfer (download or upload)
type Transfer struct {
	ID          string
	Type        TransferType
	URL         string
	FilePath    string
	FileName    string
	FileSize    int64
	Status      TransferStatus
	Progress    int64
	Speed       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	Error       string
}

// TransferType represents the type of transfer
type TransferType string

const (
	TypeDownload TransferType = "download"
	TypeUpload   TransferType = "upload"
)

// TransferStatus represents transfer status
type TransferStatus string

const (
	StatusQueued      TransferStatus = "queued"
	StatusDownloading TransferStatus = "downloading"
	StatusUploading   TransferStatus = "uploading"
	StatusPaused      TransferStatus = "paused"
	StatusCompleted   TransferStatus = "completed"
	StatusFailed      TransferStatus = "failed"
	StatusCanceled    TransferStatus = "canceled"
)

// Account represents a MEGA account
type Account struct {
	ID        string
	Email     string
	SessionID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// File represents a MEGA file
type File struct {
	ID         string
	Name       string
	Size       int64
	Type       string
	Key        string
	ParentID   string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Folder represents a MEGA folder
type Folder struct {
	ID        string
	Name      string
	ParentID  string
	CreatedAt time.Time
}
