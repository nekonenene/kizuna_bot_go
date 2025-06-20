# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go implementation of Kizuna Bot, a Discord bot that provides comprehensive features including weather reports, news articles, restaurant search, image search, video search, translation, user ranking, and conversational responses. The bot is designed to be deployed as a single binary on a server and offers full feature parity with the original Ruby version.

## Development Commands

```bash
# Install dependencies
make deps

# Build the application
make build

# Build for server deployment (Linux)
make build-linux

# Run in development mode
make dev

# Run without building
make run

# Run tests
make test

# Clean build artifacts
make clean

# Create .env file from example
make init

# Show help
make help
```

## Project Structure

```
kizuna_bot_go/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Makefile               # Build and development commands
├── .env.example           # Environment variables template
├── internal/
│   ├── bot/               # Discord bot implementation
│   │   ├── bot.go         # Main bot structure and setup
│   │   └── handlers.go    # Command and message handlers
│   ├── api/               # External API integrations
│   │   ├── client.go      # HTTP client wrapper
│   │   ├── weather.go     # Weather API integration
│   │   ├── news.go        # News API integration
│   │   └── gourmet.go     # Restaurant search API
│   └── config/
│       └── config.go      # Configuration management
└── build/                 # Build output directory
```

## Architecture

The bot follows a modular architecture:

- **main.go**: Entry point that initializes the bot and handles graceful shutdown
- **bot package**: Contains Discord bot logic, command handlers, and message processing
- **api package**: Handles all external API integrations with proper error handling
- **config package**: Manages environment variables and configuration

## Environment Configuration

Required environment variables (copy from `.env.example` to `.env`):

```
BOT_CLIENT_ID=              # Discord bot client ID
BOT_TOKEN=                  # Discord bot token
RSS2JSON_API_KEY=          # For news fetching
RECRUIT_API_KEY=           # For restaurant search
CUSTOM_SEARCH_ENGINE_ID=   # For image search
CUSTOM_SEARCH_API_KEY=     # For image search
YOUTUBE_DATA_API_KEY=      # For video search
```

## Bot Features

### Implemented Commands
- `/ping` - Response time test
- `/help` - Show available commands
- `/weather` - Tokyo weather forecast
- `/news` - Random news from Hatena
- `/dice [max]` - Roll a dice (default 6-sided)
- `/gourmet <address> [keyword]` - Restaurant search
- `/image`, `/img <query>` - Image search (Google Custom Search API)
- `/video`, `/youtube <query>` - YouTube video search
- `/vtuber [query]` - VTuber video search
- `/eng <text>` - English translation
- `/jpn`, `/jap <text>` - Japanese translation
- `/rank` - User activity ranking in channel

### Implemented Responses
- Mention responses with conversational patterns (weather, news, translation, ranking)
- Weather pattern matching ("天気は？")
- Various greeting and emotional responses
- "ゆーま？" pattern for specific channel video search
- Translation with quotes: "英語で「テキスト」" / "日本語で「テキスト」"

## Development Notes

- Uses discordgo library for Discord API integration
- HTTP client with 30-second timeout for external APIs
- Graceful shutdown handling with signal catching
- Structured logging for debugging
- Cross-platform build support (Linux, Windows, macOS)

## Building for Production

For server deployment, use:

```bash
make build-linux
```

This creates a statically linked binary in `build/kizuna_bot_go-linux` that can be deployed directly to a Linux server.

## Testing

Run tests with:

```bash
make test
```

## Coding Standards

- All files must end with a newline character
- Use Go standard formatting (gofmt)
- Follow Go naming conventions
- Include proper error handling for all external API calls
- Use structured logging for debugging information
- Write comments in Japanese to maintain consistency with the original Ruby version
- Add comprehensive comments for functions, types, and complex logic
- Document API integrations and their expected responses

### EditorConfig Rules

Follow the `.editorconfig` settings for consistent code formatting:

- **General files**: 2 spaces indentation, UTF-8 encoding, LF line endings, trim trailing whitespace
- **Go files** (`.go`): Tab indentation with 4-space width
- **Makefiles**: Tab indentation with 4-space width  
- **Markdown files** (`.md`): 4 spaces indentation, do not trim trailing whitespace (for proper line breaks)
- **All files**: Must end with a newline character (`insert_final_newline = true`)

When editing code, ensure your editor respects these `.editorconfig` settings to maintain consistency across the codebase.

## Common Issues

1. **Missing API Keys**: Ensure all required environment variables are set in `.env`
2. **Discord Permissions**: Bot needs message read/send permissions and message content intent
3. **API Rate Limits**: Some APIs have usage limits (e.g., image search is limited to 100 requests/day)
