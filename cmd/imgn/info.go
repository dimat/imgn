package main

import (
	"fmt"

	"github.com/dimat/imgn/internal/config"
	"github.com/dimat/imgn/internal/models"
	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show current configuration and status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			keyStatus := "not set"
			if cfg.APIKey() != "" {
				keyStatus = fmt.Sprintf("set (%d chars)", len(cfg.APIKey()))
			}

			model, _ := models.Resolve(cfg.Model)

			fmt.Printf("imgn %s\n\n", version)
			fmt.Printf("Config file:  %s\n", config.ConfigFilePath())
			fmt.Printf("API key:      %s\n", keyStatus)
			fmt.Printf("Model:        %s (%s)\n", cfg.Model, model.ID)
			fmt.Printf("Aspect ratio: %s\n", cfg.Aspect)
			fmt.Printf("Size:         %s\n", cfg.Size)
			fmt.Printf("Output dir:   %s\n", cfg.OutputDir)

			return nil
		},
	}
}
