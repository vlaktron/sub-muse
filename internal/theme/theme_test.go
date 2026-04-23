package theme

import (
	"testing"
)

func TestLoadOrDefault(t *testing.T) {
	colors := LoadOrDefault()
	if colors.Accent == "" {
		t.Error("Expected Accent to be non-empty")
	}
	if colors.Foreground == "" {
		t.Error("Expected Foreground to be non-empty")
	}
	if colors.Background == "" {
		t.Error("Expected Background to be non-empty")
	}
}

func TestDefaultColors(t *testing.T) {
	defaults := Default
	if defaults.Accent != "#82FB9C" {
		t.Errorf("Expected Accent to be #82FB9C, got %s", defaults.Accent)
	}
	if defaults.Foreground != "#ddf7ff" {
		t.Errorf("Expected Foreground to be #ddf7ff, got %s", defaults.Foreground)
	}
	if defaults.Background != "#0B0C16" {
		t.Errorf("Expected Background to be #0B0C16, got %s", defaults.Background)
	}
}
