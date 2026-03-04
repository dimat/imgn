package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dimat/imgn/internal/client"
	"github.com/dimat/imgn/internal/config"
	"github.com/dimat/imgn/internal/models"
	"github.com/dimat/imgn/internal/output"
	"github.com/spf13/cobra"
)

type generateFlags struct {
	model      string
	aspect     string
	size       string
	outputFile string
	outputDir  string
	count      int
	negative   string
	images     []string
	promptFile string
	verbose    bool
	quiet      bool
	jsonOutput bool
}

func newGenerateCmd() *cobra.Command {
	f := &generateFlags{}

	cmd := &cobra.Command{
		Use:     "generate [prompt...]",
		Aliases: []string{"gen", "g"},
		Short:   "Generate images from a text prompt",
		Long: `Generate one or more images from a text prompt using Google Gemini.

Prompts can be provided as arguments, from a file (--prompt-file or @filename),
or piped via stdin. Multiple sources are concatenated.

Examples:
  imgn generate "a sunset over mountains"
  imgn g "a cat wearing a hat" --model flash --aspect 1:1
  imgn generate @prompt.txt --count 3
  echo "a robot painting" | imgn generate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(cmd.Context(), f, args)
		},
	}

	cmd.Flags().StringVarP(&f.model, "model", "m", "", "Model to use: flash2, pro, flash (default from config)")
	cmd.Flags().StringVarP(&f.aspect, "aspect", "a", "", "Aspect ratio: 1:1, 16:9, 9:16, 4:3, 3:4")
	cmd.Flags().StringVar(&f.size, "size", "", "Resolution: 512px, 1k, 2k, 4k")
	cmd.Flags().StringVarP(&f.outputFile, "output", "o", "", "Output filename")
	cmd.Flags().StringVar(&f.outputDir, "output-dir", "", "Output directory")
	cmd.Flags().IntVarP(&f.count, "count", "n", 1, "Number of images to generate")
	cmd.Flags().StringVar(&f.negative, "negative", "", "Negative prompt (things to avoid)")
	cmd.Flags().StringArrayVarP(&f.images, "image", "i", nil, "Reference image file(s)")
	cmd.Flags().StringVar(&f.promptFile, "prompt-file", "", "Read prompt from file")
	cmd.Flags().BoolVar(&f.verbose, "verbose", false, "Verbose output")
	cmd.Flags().BoolVar(&f.quiet, "quiet", false, "Suppress non-essential output")
	cmd.Flags().BoolVar(&f.jsonOutput, "json", false, "Output structured JSON")

	return cmd
}

func runGenerate(ctx context.Context, f *generateFlags, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if cfg.APIKey() == "" {
		msg := "API key not set. Set IMGN_API_KEY or GEMINI_API_KEY env var, or add providers.google.api_key to config file."
		if f.jsonOutput {
			output.PrintJSONError(msg)
		}
		return fmt.Errorf("%s", msg)
	}

	// Resolve model
	modelName := cfg.Model
	if f.model != "" {
		modelName = f.model
	}
	model, err := models.Resolve(modelName)
	if err != nil {
		if f.jsonOutput {
			output.PrintJSONError(err.Error())
			return err
		}
		return err
	}

	// Resolve aspect ratio
	aspect := cfg.Aspect
	if f.aspect != "" {
		aspect = f.aspect
	}
	if err := models.ValidateAspect(aspect, model); err != nil {
		if f.jsonOutput {
			output.PrintJSONError(err.Error())
			return err
		}
		return err
	}

	// Resolve size
	size := cfg.Size
	if f.size != "" {
		size = f.size
	}
	if err := models.ValidateSize(size, model); err != nil {
		if f.jsonOutput {
			output.PrintJSONError(err.Error())
			return err
		}
		return err
	}

	// Resolve output dir
	outputDir := cfg.OutputDir
	if f.outputDir != "" {
		outputDir = f.outputDir
	}

	// Build prompt from args, files, and stdin
	prompt, err := buildPrompt(args, f.promptFile)
	if err != nil {
		if f.jsonOutput {
			output.PrintJSONError(err.Error())
			return err
		}
		return err
	}
	if prompt == "" {
		msg := "no prompt provided"
		if f.jsonOutput {
			output.PrintJSONError(msg)
		}
		return fmt.Errorf("%s", msg)
	}

	// Load reference images
	var imageInputs []client.ImageInput
	for _, imgPath := range f.images {
		data, err := os.ReadFile(imgPath)
		if err != nil {
			return fmt.Errorf("read image %s: %w", imgPath, err)
		}
		mime := client.DetectMIMEType(data)
		imageInputs = append(imageInputs, client.ImageInput{MimeType: mime, Data: data})
	}

	// Log
	if f.verbose {
		slog.Info("generating image",
			"model", model.ID,
			"aspect", aspect,
			"size", size,
			"count", f.count,
			"prompt_len", len(prompt),
			"images", len(imageInputs),
		)
	}
	if !f.quiet && !f.jsonOutput {
		fmt.Fprintf(os.Stderr, "Generating with %s (%s, %s)...\n", model.Alias, aspect, size)
	}

	// Map size to API format (512px passes through as-is, others uppercase)
	apiSize := size
	if size != "512px" {
		apiSize = strings.ToUpper(size)
	}

	// Generate images (one call per image since API returns 1 image per call)
	c := client.New(cfg.APIKey())
	var allFiles []output.FileResult
	var lastResp *client.GenerateResponse

	for i := 0; i < f.count; i++ {
		resp, err := c.Generate(ctx, client.GenerateRequest{
			Model:       model.ID,
			Prompt:      prompt,
			Negative:    f.negative,
			AspectRatio: aspect,
			ImageSize:   apiSize,
			Images:      imageInputs,
		})
		if err != nil {
			if f.jsonOutput {
				output.PrintJSONError(err.Error())
				return err
			}
			return fmt.Errorf("generate image: %w", err)
		}
		lastResp = resp

		for _, result := range resp.Results {
			filePath := output.ResolveFilename(f.outputFile, outputDir, i, f.count, result.MimeType)
			if err := output.WriteFile(filePath, result.ImageData); err != nil {
				return fmt.Errorf("write output: %w", err)
			}

			allFiles = append(allFiles, output.FileResult{
				Path:     filePath,
				Size:     len(result.ImageData),
				MimeType: result.MimeType,
				Index:    i,
			})

			if !f.quiet && !f.jsonOutput {
				fmt.Fprintf(os.Stderr, "Saved: %s (%d bytes)\n", filePath, len(result.ImageData))
			}
		}
	}

	if f.jsonOutput {
		jsonOut := output.JSONOutput{
			Success:     true,
			Model:       model.ID,
			Prompt:      prompt,
			AspectRatio: aspect,
			Size:        size,
			Files:       allFiles,
		}
		if lastResp != nil && lastResp.TotalTokens > 0 {
			jsonOut.TokenUsage = &output.TokenUsage{
				Prompt:     lastResp.PromptTokens,
				Candidates: lastResp.CandidatesTokens,
				Total:      lastResp.TotalTokens,
			}
		}
		return output.PrintJSON(jsonOut)
	}

	// Print file paths to stdout for easy piping
	for _, f := range allFiles {
		fmt.Println(f.Path)
	}

	return nil
}

func buildPrompt(args []string, promptFile string) (string, error) {
	var parts []string

	// From args (handle @filename syntax)
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			data, err := os.ReadFile(strings.TrimPrefix(arg, "@"))
			if err != nil {
				return "", fmt.Errorf("read prompt file %s: %w", arg, err)
			}
			parts = append(parts, strings.TrimSpace(string(data)))
		} else {
			parts = append(parts, arg)
		}
	}

	// From --prompt-file
	if promptFile != "" {
		data, err := os.ReadFile(promptFile)
		if err != nil {
			return "", fmt.Errorf("read prompt file: %w", err)
		}
		parts = append(parts, strings.TrimSpace(string(data)))
	}

	// From stdin (only if not a TTY and no other prompt sources)
	if len(parts) == 0 {
		stat, _ := os.Stdin.Stat()
		if stat != nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if len(lines) > 0 {
				parts = append(parts, strings.Join(lines, "\n"))
			}
		}
	}

	return strings.Join(parts, "\n"), nil
}
