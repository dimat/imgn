package main

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "imgn",
		Short: "AI image generation CLI powered by Google Gemini",
		Long: `imgn — Generate images from text prompts using Google's Gemini models.

Supports both Gemini Flash (fast) and Gemini Pro (high quality, up to 4K).
Designed to be used by humans and AI agents alike.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(newGenerateCmd())
	root.AddCommand(newModelsCmd())
	root.AddCommand(newInfoCmd())
	root.AddCommand(newVersionCmd())

	return root
}
