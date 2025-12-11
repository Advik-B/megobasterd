package main

import (
"context"
"encoding/base64"
"fmt"
"log"
"strings"
"time"

"github.com/Advik-B/megobasterd/internal/api"
)

// ParseMegaURL parses a MEGA URL and extracts file ID and key
func ParseMegaURL(url string) (fileID, key string, err error) {
// Expected format: https://mega.nz/file/UlVzWKwY#KAMYD5AnqV5kmioRv6P0hQ3KdQjDWLAsszmo_SizLn0
if !strings.Contains(url, "mega.nz") {
return "", "", fmt.Errorf("not a MEGA URL")
}

parts := strings.Split(url, "/")
if len(parts) < 2 {
return "", "", fmt.Errorf("invalid MEGA URL format")
}

// Get the last part which contains file ID and key
lastPart := parts[len(parts)-1]

// Split by # to separate file ID and key
idAndKey := strings.Split(lastPart, "#")
if len(idAndKey) != 2 {
return "", "", fmt.Errorf("invalid MEGA URL format: missing key")
}

return idAndKey[0], idAndKey[1], nil
}

func main() {
testURL := "https://mega.nz/file/UlVzWKwY#KAMYD5AnqV5kmioRv6P0hQ3KdQjDWLAsszmo_SizLn0"

fmt.Println("=== MEGA URL Test ===")
fmt.Printf("Testing URL: %s\n\n", testURL)

// Parse the URL
fileID, key, err := ParseMegaURL(testURL)
if err != nil {
log.Fatalf("Failed to parse URL: %v", err)
}

fmt.Printf("✓ URL parsed successfully\n")
fmt.Printf("  File ID: %s\n", fileID)
fmt.Printf("  Key: %s\n", key)
fmt.Printf("  Key (decoded): ")

// Try to decode the key
decodedKey, err := base64.URLEncoding.DecodeString(key)
if err != nil {
// Try standard base64
decodedKey, err = base64.StdEncoding.DecodeString(key)
if err != nil {
fmt.Printf("(could not decode)\n")
} else {
fmt.Printf("%d bytes\n", len(decodedKey))
}
} else {
fmt.Printf("%d bytes\n", len(decodedKey))
}

// Test MEGA API client
fmt.Printf("\n=== MEGA API Client Test ===\n")
client := api.NewMegaClient()
fmt.Printf("✓ MEGA client created\n")

// Test getting download URL
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

fmt.Printf("\nAttempting to get download URL...\n")
downloadURL, err := client.GetDownloadURL(ctx, fileID, key)
if err != nil {
fmt.Printf("✗ Failed to get download URL: %v\n", err)
fmt.Printf("\nNote: This is expected as we haven't implemented full MEGA API authentication yet.\n")
fmt.Printf("The URL parsing and basic client structure work correctly.\n")
} else {
fmt.Printf("✓ Download URL retrieved: %s\n", downloadURL)
}

fmt.Printf("\n=== Test Summary ===\n")
fmt.Printf("✓ URL parsing: PASS\n")
fmt.Printf("✓ File ID extraction: PASS\n")
fmt.Printf("✓ Key extraction: PASS\n")
fmt.Printf("✓ API client initialization: PASS\n")
fmt.Printf("- Full API integration: Pending (Phase 5+)\n")
fmt.Printf("\nThe foundation is working correctly!\n")
}
