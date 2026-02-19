# imgn — AI Image Generation CLI

Generate images from text prompts using Google's Gemini models, right from your terminal.

## Features

- **Two models**: Gemini Flash (fast) and Gemini Pro (high quality, up to 4K)
- **Reference images**: Pass existing images to guide generation
- **Flexible prompts**: From args, files (`@prompt.txt`), or stdin
- **AI agent friendly**: `--json` output, predictable exit codes, stderr/stdout separation
- **Configurable**: Env vars, config file, or flags

## Installation

```bash
go install github.com/dimat/imgn/cmd/imgn@latest
```

Or build from source:

```bash
git clone https://github.com/dimat/imgn.git
cd imgn
go build -o imgn ./cmd/imgn
```

## Quick Start

```bash
# Set your API key
export GEMINI_API_KEY="your-key-here"

# Generate an image
imgn generate "a sunset over mountains in watercolor style"

# Use the flash model for faster results
imgn g "a cute robot" --model flash --aspect 1:1

# Generate multiple images
imgn g "abstract art" --count 3 --size 4k

# Use a reference image
imgn g "make this image look like a painting" -i photo.jpg

# Read prompt from file
imgn g @prompt.txt

# Combine multiple prompt sources (style + composition + instruction)
imgn g @style.txt @composition.txt "a dragon sitting on a throne"

# Pipe from stdin
echo "a cat wearing a top hat" | imgn generate

# JSON output for scripts/agents
imgn g "a logo" --json
```

## Configuration

Set up a config file at `~/.config/imgn/config.yaml`:

```yaml
provider: google       # active provider (google for now, more coming)

providers:
  google:
    api_key: "your-key"  # or use IMGN_API_KEY / GEMINI_API_KEY env vars

model: "pro"           # default model (pro or flash)
aspect: "16:9"         # default aspect ratio
size: "2k"             # default resolution
output_dir: "."        # where to save images
```

Priority: flags > env vars > config file > defaults.

### Prompt Composition

Multiple prompt sources are concatenated in order — great for reusable style/composition presets:

```bash
# style.txt: "hyper-realistic digital painting, cinematic lighting, 8k detail"
# composition.txt: "centered subject, rule of thirds, shallow depth of field"
imgn g @style.txt @composition.txt "a knight standing in a misty forest"
```

You can also use `--prompt-file` alongside args:

```bash
imgn g --prompt-file style.txt "a portrait of a cat"
```

## Commands

| Command | Description |
|---------|-------------|
| `imgn generate` (alias: `gen`, `g`) | Generate images from a prompt |
| `imgn models` | List available models |
| `imgn info` | Show current configuration |
| `imgn version` | Print version info |

## Documentation

- [Detailed Usage Guide](docs/USAGE.md)
- [Model Comparison](docs/MODELS.md)
- [AI Agent Guide](docs/AI-AGENT-GUIDE.md)

## License

MIT
