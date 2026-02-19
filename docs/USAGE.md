# imgn Usage Guide

## Generate Command

```bash
imgn generate [flags] [prompt...]
imgn gen [flags] [prompt...]
imgn g [flags] [prompt...]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--model` | `-m` | `pro` | Model: `pro`, `flash` |
| `--aspect` | `-a` | `16:9` | Aspect ratio: `1:1`, `16:9`, `9:16`, `4:3`, `3:4` |
| `--size` | | `2k` | Resolution: `1k`, `2k`, `4k` (4K: Pro only) |
| `--output` | `-o` | auto | Output filename |
| `--output-dir` | | `.` | Output directory |
| `--count` | `-n` | `1` | Number of images |
| `--negative` | | | Things to avoid in the image |
| `--image` | `-i` | | Reference image file (repeatable) |
| `--prompt-file` | | | Read prompt from file |
| `--json` | | | Structured JSON output |
| `--verbose` | | | Verbose logging |
| `--quiet` | | | Suppress non-essential output |

### Prompt Sources

Prompts can come from multiple sources (concatenated in order):

1. **Command arguments**: `imgn g "a sunset" "over mountains"`
2. **File reference**: `imgn g @prompt.txt`
3. **--prompt-file flag**: `imgn g --prompt-file prompt.txt`
4. **stdin** (when no other prompt): `echo "a cat" | imgn g`

### Output Filenames

Default: `imgn-YYYYMMDD-HHMMSS.png`

With `--count 3`: `imgn-YYYYMMDD-HHMMSS-1.png`, `-2.png`, `-3.png`

With `--output hero.png --count 2`: `hero-1.png`, `hero-2.png`

### Reference Images

Pass one or more images to guide generation:

```bash
imgn g "make this look like a watercolor painting" -i photo.jpg
imgn g "combine these styles" -i style1.png -i style2.png
```

Supported formats: PNG, JPEG, WebP, GIF.

### JSON Output

```bash
imgn g "a logo" --json
```

```json
{
  "success": true,
  "model": "gemini-3-pro-image-preview",
  "prompt": "a logo",
  "aspect_ratio": "16:9",
  "size": "2k",
  "files": [
    {
      "path": "imgn-20260219-185500.png",
      "size_bytes": 1234567,
      "mime_type": "image/png",
      "index": 0
    }
  ],
  "token_usage": {
    "prompt_tokens": 10,
    "candidates_tokens": 100,
    "total_tokens": 110
  }
}
```

## Other Commands

### Models

```bash
imgn models
```

Lists all available models with their capabilities and supported sizes.

### Info

```bash
imgn info
```

Shows current configuration, API key status, and defaults.

### Version

```bash
imgn version
```
