# imgn — AI Image Generation CLI

## What this tool does

Generates images from text prompts using Google's Gemini models (Flash and Pro). Supports reference images, multiple aspect ratios, up to 4K resolution, and structured JSON output for agent consumption.

## Installation

```bash
go install github.com/dimat/imgn/cmd/imgn@latest
```

## Quick Reference

```bash
# Basic generation
imgn generate "a sunset over mountains"
imgn g "a cute robot" --model flash

# With options
imgn g "abstract art" --aspect 1:1 --size 4k --count 3

# Reference image
imgn g "make this a watercolor" -i photo.jpg

# JSON output for agents
imgn g "a logo" --json

# Combine prompts (style preset + instruction)
imgn g @style.txt @composition.txt "a dragon on a throne"

# From stdin
echo "a cat" | imgn g
```

## Configuration

**Env vars** (highest priority):
- `IMGN_API_KEY` — API key (preferred)
- `GEMINI_API_KEY` — API key (fallback)

**Config file**: `~/.config/imgn/config.yaml`
```yaml
provider: google
providers:
  google:
    api_key: "..."
model: "pro"        # pro or flash
aspect: "16:9"      # 1:1, 16:9, 9:16, 4:3, 3:4
size: "2k"          # 1k, 2k, 4k (4k: pro only)
output_dir: "."
```

### Prompt Composition
Multiple `@file` args and text args are concatenated in order — use for reusable style/composition presets:
```bash
imgn g @style.txt @composition.txt "specific subject description"
```

## Commands

### `imgn generate [prompt...]` (aliases: `gen`, `g`)
Generate images from text.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--model` | `-m` | `pro` | Model: `pro`, `flash` |
| `--aspect` | `-a` | `16:9` | Aspect ratio |
| `--size` | | `2k` | Resolution: `1k`, `2k`, `4k` |
| `--output` | `-o` | auto | Output filename |
| `--output-dir` | | `.` | Output directory |
| `--count` | `-n` | `1` | Number of images |
| `--negative` | | | Negative prompt |
| `--image` | `-i` | | Reference image (repeatable) |
| `--prompt-file` | | | Read prompt from file |
| `--json` | | | JSON output |
| `--verbose` | | | Verbose mode |
| `--quiet` | | | Quiet mode |

### `imgn models`
List available models with capabilities.

### `imgn info`
Show current config, API key status.

### `imgn version`
Print version.

## AI Agent Usage

**JSON output**: `imgn g "prompt" --json` returns:
```json
{
  "success": true,
  "model": "gemini-3-pro-image-preview",
  "files": [{"path": "imgn-20260219-185500.png", "size_bytes": 45231, "mime_type": "image/png", "index": 0}],
  "token_usage": {"prompt_tokens": 5, "candidates_tokens": 200, "total_tokens": 205}
}
```

**Exit codes**: 0=success, 1=runtime error, 2=usage error

**stdout**: File paths (or JSON with `--json`). **stderr**: Messages/errors.

**stdin**: `echo "prompt" | imgn g`

**Capture path**: `FILE=$(imgn g "prompt" --quiet)`

## Troubleshooting

| Error | Fix |
|-------|-----|
| "API key not set" | Set `IMGN_API_KEY` env var, or add `providers.google.api_key` to `~/.config/imgn/config.yaml` |
| "size 4k not supported by flash" | Use `--model pro` for 4K |
| "no images generated" | Prompt may have been blocked by safety filters; rephrase |
| Timeout | API can take 15-30s for Pro; increase timeout or use Flash |
