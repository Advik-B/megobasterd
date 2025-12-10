package utils

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		result := FormatBytes(tt.input)
		if result != tt.expected {
			t.Errorf("FormatBytes(%d) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{512.0, "512 B/s"},
		{1024.0, "1.00 KB/s"},
		{1572864.0, "1.50 MB/s"},
		{10485760.0, "10.00 MB/s"},
	}

	for _, tt := range tests {
		result := FormatSpeed(tt.input)
		if result != tt.expected {
			t.Errorf("FormatSpeed(%.0f) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{30, "30s"},
		{90, "1m 30s"},
		{3600, "1h 0m"},
		{7200, "2h 0m"},
		{86400, "1d 0h"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDuration(%d) = %s; want %s", tt.input, result, tt.expected)
		}
	}
}

func TestCalculateETA(t *testing.T) {
	tests := []struct {
		downloaded int64
		total      int64
		speed      float64
		expected   int64
	}{
		{0, 1000, 100.0, 10},
		{500, 1000, 100.0, 5},
		{1000, 1000, 100.0, 0},
		{500, 1000, 0.0, 0},
	}

	for _, tt := range tests {
		result := CalculateETA(tt.downloaded, tt.total, tt.speed)
		if result != tt.expected {
			t.Errorf("CalculateETA(%d, %d, %.1f) = %d; want %d", 
				tt.downloaded, tt.total, tt.speed, result, tt.expected)
		}
	}
}

func TestClampInt(t *testing.T) {
	tests := []struct {
		value    int
		min      int
		max      int
		expected int
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
	}

	for _, tt := range tests {
		result := ClampInt(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("ClampInt(%d, %d, %d) = %d; want %d", 
				tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}
