// Package service provides shared business logic used by CLI commands.
// It extracts reusable operations from cmd/ to reduce file size and
// enable testing without cobra dependencies.
package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFile downloads a URL to a local file path.
func DownloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, data, 0644)
}

// ExtractExt returns the file extension from a URL, defaulting to ".mp4".
func ExtractExt(rawURL string) string {
	if idx := strings.Index(rawURL, "?"); idx != -1 {
		rawURL = rawURL[:idx]
	}
	ext := filepath.Ext(rawURL)
	if ext == "" {
		return ".mp4"
	}
	return ext
}

// SavePrompt writes the generation prompt alongside result files.
func SavePrompt(outputDir, taskID, prompt string) {
	if prompt == "" {
		return
	}
	filename := filepath.Join(outputDir, fmt.Sprintf("image_%s.md", taskID))
	content := fmt.Sprintf("# %s\n\n%s\n", taskID, prompt)
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save prompt file: %v\n", err)
		return
	}
	fmt.Printf("Prompt saved: %s\n", filename)
}
