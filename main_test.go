package main

import (
	"testing"
)

func TestExtractVideoURL(t *testing.T) {
	scriptContent := `,\"video_url\":\"https:\\\/\\\/example.com\\\/video\",`
	expectedURL := "https://example.com/video"

	videoURL, err := extractVideoURL(scriptContent)
	if err != nil {
		t.Error(err)
	}

	if videoURL != expectedURL {
		t.Errorf("expected: %s, got %s", expectedURL, videoURL)
	}
}
