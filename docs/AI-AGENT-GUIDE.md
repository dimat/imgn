# imgn — AI Agent Integration Guide

## Overview

`imgn` is designed for programmatic use by AI coding agents. Use `--json` for structured output.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Runtime error (API failure, network error) |
| 2 | Usage error (invalid flags, missing prompt) |

## Output Conventions

- **stdout**: File paths (default) or JSON (`--json`)
- **stderr**: Progress messages, warnings, errors
- Non-interactive by default (no spinners, no confirmations)

## JSON Output Format

```bash
imgn generate "a blue circle" --json
```

Success:
```json
{
  "success": true,
  "model": "gemini-3-pro-image-preview",
  "prompt": "a blue circle",
  "aspect_ratio": "16:9",
  "size": "2k",
  "files": [
    {
      "path": "imgn-20260219-185500.png",
      "size_bytes": 45231,
      "mime_type": "image/png",
      "index": 0
    }
  ],
  "token_usage": {
    "prompt_tokens": 5,
    "candidates_tokens": 200,
    "total_tokens": 205
  }
}
```

Error:
```json
{
  "success": false,
  "error": "gemini API error (429): Resource exhausted"
}
```

## Common Patterns

### Generate and capture path
```bash
PATH=$(imgn g "icon" --quiet)
echo "Image saved to: $PATH"
```

### Batch generation
```bash
for prompt in "a sun" "a moon" "a star"; do
  imgn g "$prompt" --json --output-dir ./batch/
done
```

### Pipe prompt from another tool
```bash
echo "Generate an icon for a weather app" | imgn g --json
```

### Check configuration
```bash
imgn info  # verify API key is set before running
```

## Error Handling

Always check exit code. With `--json`, parse the `success` field:

```bash
result=$(imgn g "test" --json 2>/dev/null)
if echo "$result" | jq -e '.success' > /dev/null 2>&1; then
  file=$(echo "$result" | jq -r '.files[0].path')
else
  error=$(echo "$result" | jq -r '.error')
fi
```
