package app

import (
"testing"
)

func TestParseMegaURL(t *testing.T) {
app := NewApp()

tests := []struct {
name       string
url        string
wantFileID string
wantKey    string
wantErr    bool
}{
{
name:       "Valid MEGA URL",
url:        "https://mega.nz/file/UlVzWKwY#KAMYD5AnqV5kmioRv6P0hQ3KdQjDWLAsszmo_SizLn0",
wantFileID: "UlVzWKwY",
wantKey:    "KAMYD5AnqV5kmioRv6P0hQ3KdQjDWLAsszmo_SizLn0",
wantErr:    false,
},
{
name:    "Invalid URL - not MEGA",
url:     "https://example.com/file",
wantErr: true,
},
{
name:    "Invalid URL - missing key",
url:     "https://mega.nz/file/UlVzWKwY",
wantErr: true,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
fileID, key, err := app.ParseMegaURL(tt.url)

if tt.wantErr {
if err == nil {
t.Errorf("ParseMegaURL() expected error, got nil")
}
return
}

if err != nil {
t.Errorf("ParseMegaURL() unexpected error: %v", err)
return
}

if fileID != tt.wantFileID {
t.Errorf("ParseMegaURL() fileID = %v, want %v", fileID, tt.wantFileID)
}

if key != tt.wantKey {
t.Errorf("ParseMegaURL() key = %v, want %v", key, tt.wantKey)
}
})
}
}

func TestDownloadStruct(t *testing.T) {
download := &Download{
ID:         "test123",
FileName:   "test.zip",
FileSize:   1024 * 1024, // 1MB
Downloaded: 512 * 1024,  // 512KB
Status:     "downloading",
}

if download.FileSize > 0 {
download.Progress = float64(download.Downloaded) / float64(download.FileSize) * 100
}

expectedProgress := 50.0
if download.Progress != expectedProgress {
t.Errorf("Progress calculation = %v, want %v", download.Progress, expectedProgress)
}
}
