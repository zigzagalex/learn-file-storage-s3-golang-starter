package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

type ffprobeOutput struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func getVideoAspectRatio(filePath string) (string, error) {
	// Prepare the ffprobe command with separate args
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)

	// Capture stdout into a buffer
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe failed: %w", err)
	}

	// Parse the JSON
	var probe ffprobeOutput
	if err := json.Unmarshal(outBuf.Bytes(), &probe); err != nil {
		return "", fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	if len(probe.Streams) == 0 {
		return "", fmt.Errorf("no streams found")
	}

	width := probe.Streams[0].Width
	height := probe.Streams[0].Height
	if width == 0 || height == 0 {
		return "", fmt.Errorf("invalid width/height in video metadata")
	}

	// Compare with tolerance
	ratio := float64(width) / float64(height)
	if almostEqual(ratio, 16.0/9.0, 0.01) {
		return "16:9", nil
	}
	if almostEqual(ratio, 9.0/16.0, 0.01) {
		return "9:16", nil
	}
	return "other", nil
}
