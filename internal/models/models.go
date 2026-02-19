// Package models defines the available Gemini image generation models and their capabilities.
package models

import "fmt"

// Model represents a Gemini image generation model.
type Model struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Alias       string   `json:"alias"`
	Description string   `json:"description"`
	MaxSize     string   `json:"max_size"` // max resolution: "2k" or "4k"
	Aspects     []string `json:"aspects"`  // supported aspect ratios
	Sizes       []string `json:"sizes"`    // supported sizes
}

var (
	Flash = Model{
		ID:          "gemini-2.0-flash-preview-image-generation",
		Name:        "Gemini Flash (Nano Banana)",
		Alias:       "flash",
		Description: "Fast image generation, up to 2K resolution",
		MaxSize:     "2k",
		Aspects:     []string{"1:1", "16:9", "9:16", "4:3", "3:4"},
		Sizes:       []string{"1k", "2k"},
	}

	Pro = Model{
		ID:          "gemini-3-pro-image-preview",
		Name:        "Gemini Pro (Nano Banana Pro)",
		Alias:       "pro",
		Description: "High-quality image generation, up to 4K resolution",
		MaxSize:     "4k",
		Aspects:     []string{"1:1", "16:9", "9:16", "4:3", "3:4"},
		Sizes:       []string{"1k", "2k", "4k"},
	}

	All = []Model{Pro, Flash}
)

// DefaultModel is the default model used when none is specified.
const DefaultModel = "pro"

// Resolve returns the Model for a given alias or full model ID.
func Resolve(name string) (Model, error) {
	switch name {
	case "pro", Pro.ID:
		return Pro, nil
	case "flash", Flash.ID:
		return Flash, nil
	default:
		return Model{}, fmt.Errorf("unknown model %q (valid: pro, flash)", name)
	}
}

// ValidateAspect checks if the aspect ratio is valid.
func ValidateAspect(aspect string) error {
	for _, a := range Pro.Aspects {
		if a == aspect {
			return nil
		}
	}
	return fmt.Errorf("invalid aspect ratio %q (valid: 1:1, 16:9, 9:16, 4:3, 3:4)", aspect)
}

// ValidateSize checks if the size is valid for the given model.
func ValidateSize(size string, model Model) error {
	for _, s := range model.Sizes {
		if s == size {
			return nil
		}
	}
	return fmt.Errorf("size %q not supported by %s (valid: %v)", size, model.Alias, model.Sizes)
}
