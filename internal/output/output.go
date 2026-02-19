// Package output handles writing generated images to disk and formatting output.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileResult represents a written output file.
type FileResult struct {
	Path     string `json:"path"`
	Size     int    `json:"size_bytes"`
	MimeType string `json:"mime_type"`
	Index    int    `json:"index"`
}

// JSONOutput is the structured output for --json mode.
type JSONOutput struct {
	Success      bool         `json:"success"`
	Model        string       `json:"model"`
	Prompt       string       `json:"prompt"`
	AspectRatio  string       `json:"aspect_ratio"`
	Size         string       `json:"size"`
	Files        []FileResult `json:"files"`
	TokenUsage   *TokenUsage  `json:"token_usage,omitempty"`
	ErrorMessage string       `json:"error,omitempty"`
}

// TokenUsage holds token consumption info.
type TokenUsage struct {
	Prompt     int `json:"prompt_tokens"`
	Candidates int `json:"candidates_tokens"`
	Total      int `json:"total_tokens"`
}

// GenerateFilename creates a default output filename with timestamp.
func GenerateFilename(outputDir, ext string, index, count int) string {
	ts := time.Now().Format("20060102-150405")
	name := fmt.Sprintf("imgn-%s", ts)
	if count > 1 {
		name = fmt.Sprintf("%s-%d", name, index+1)
	}
	return filepath.Join(outputDir, name+ext)
}

// ResolveFilename resolves the output path, handling user-provided names and multi-image indexing.
func ResolveFilename(userOutput, outputDir string, index, count int, mimeType string) string {
	ext := ExtForMIME(mimeType)

	if userOutput == "" {
		return GenerateFilename(outputDir, ext, index, count)
	}

	if count > 1 {
		base := strings.TrimSuffix(userOutput, filepath.Ext(userOutput))
		origExt := filepath.Ext(userOutput)
		if origExt == "" {
			origExt = ext
		}
		return filepath.Join(outputDir, fmt.Sprintf("%s-%d%s", base, index+1, origExt))
	}

	if filepath.Ext(userOutput) == "" {
		userOutput += ext
	}
	if !filepath.IsAbs(userOutput) {
		userOutput = filepath.Join(outputDir, userOutput)
	}
	return userOutput
}

// ExtForMIME returns the file extension for a MIME type.
func ExtForMIME(mime string) string {
	switch mime {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".png"
	}
}

// WriteFile writes image data to a file, creating directories as needed.
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// PrintJSON writes structured JSON output to stdout.
func PrintJSON(out JSONOutput) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// PrintJSONError writes a JSON error to stdout.
func PrintJSONError(errMsg string) {
	out := JSONOutput{
		Success:      false,
		ErrorMessage: errMsg,
	}
	_ = PrintJSON(out)
}
