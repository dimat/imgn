package main

import (
	"fmt"
	"strings"

	"github.com/dimat/imgn/internal/models"
	"github.com/spf13/cobra"
)

func newModelsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "models",
		Short: "List available models and their capabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, m := range models.All {
				fmt.Printf("%-8s %s\n", m.Alias, m.ID)
				fmt.Printf("         %s\n", m.Description)
				fmt.Printf("         Sizes: %s | Aspects: %s\n\n",
					strings.Join(m.Sizes, ", "),
					strings.Join(m.Aspects, ", "))
			}
			return nil
		},
	}
}
