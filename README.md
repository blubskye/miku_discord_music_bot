# Miku Bot - Discord Music Bot

A feature-rich Discord music bot written in Go with role-based permissions and SQLite persistence.

## Features

### Multi-Source Support
- YouTube
- SoundCloud
- Bandcamp
- Vimeo
- Twitch streams
- Local files
- HTTP URLs

### Supported Audio Formats
- MP3
- FLAC
- WAV
- Matroska/WebM (AAC, Opus, Vorbis codecs)
- MP4/M4A (AAC codec)
- OGG streams (Opus, Vorbis, FLAC codecs)
- AAC streams
- Stream playlists (M3U, PLS)

### Role-Based Permissions
- **User**: Can add music to the queue
- **DJ**: Can add music, skip tracks, remove tracks, move tracks to top, control playback (pause/resume), adjust volume
- **Moderator**: All DJ permissions + can stop playback and clear the entire queue
- **Admin**: All permissions + can configure bot settings and roles

### Database Persistence
- SQLite database for persistent storage
- Guild-specific settings
- Queue persistence across restarts
- Playback history tracking

## Prerequisites

Before running Miku Bot, you need to install:

1. **Go** (1.21 or higher)
   ```bash
   # Download from https://go.dev/dl/
   ```

2. **FFmpeg**
   ```bash
   # Ubuntu/Debian
   sudo apt install ffmpeg

   # macOS
   brew install ffmpeg

   # Windows
   # Download from https://ffmpeg.org/download.html
   ```

3. **yt-dlp**
   ```bash
   # Using pip
   pip install yt-dlp

   # Or download binary from https://github.com/yt-dlp/yt-dlp/releases
   ```

4. **GCC** (required for sqlite3)
   ```bash
   # Ubuntu/Debian
   sudo apt install build-essential

   # macOS
   xcode-select --install

   # Windows
   # Install MinGW-w64 or TDM-GCC
   ```

## Installation

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd miku_bot_go
   ```

2. Install dependencies:
   ```bash
   go work init
   go work use .
   go mod download
   ```

3. Create a Discord bot:
   - Go to [Discord Developer Portal](https://discord.com/developers/applications)
   - Create a new application
   - Go to the "Bot" section and create a bot
   - Enable the following Privileged Gateway Intents:
     - SERVER MEMBERS INTENT
     - MESSAGE CONTENT INTENT
   - Copy the bot token

4. Configure the bot:
   ```bash
   cp .env.example .env
   ```

   Edit `.env` and add your Discord bot token:
   ```
   DISCORD_TOKEN=your_discord_bot_token_here
   ```

5. Customize `configs/config.yaml` if needed:
   ```yaml
   bot:
     prefix: "!"
     activity: "Music | !help"
     status: "online"

   database:
     path: "miku_bot.db"

   roles:
     dj: "DJ"
     mod: "Moderator"

   music:
     max_queue_size: 100
     default_volume: 50
     timeout: 300
   ```

6. Build and run:
   ```bash
   go build -o miku_bot ./cmd/bot
   ./miku_bot
   ```

   Or run directly:
   ```bash
   go run ./cmd/bot
   ```

## Command Line Flags

The bot supports the following command-line flags:

```bash
./miku_bot --help        # Show help message
./miku_bot --version     # Show version, author, and license info
./miku_bot --trace       # Enable full stack tracing for debugging
```

Examples:
```bash
# Run with stack tracing enabled
./miku_bot --trace

# Show version information
./miku_bot --version

# Normal operation
./miku_bot
```

## Bot Invite

To invite the bot to your server, create an invite link with the following permissions:

**Required Permissions:**
- View Channels
- Send Messages
- Embed Links
- Read Message History
- Connect (to voice channels)
- Speak (in voice channels)
- Use Voice Activity

**OAuth2 URL Generator:**
```
https://discord.com/api/oauth2/authorize?client_id=YOUR_CLIENT_ID&permissions=3165184&scope=bot
```

Replace `YOUR_CLIENT_ID` with your bot's client ID from the Discord Developer Portal.

## Commands

### Music Commands

| Command | Description | Permission |
|---------|-------------|------------|
| `!play <url/query>` | Play a song from URL or search query | User+ |
| `!skip` / `!s` | Skip the current song | DJ+ |
| `!stop` | Stop playback and clear queue | Mod+ |
| `!pause` | Pause playback | DJ+ |
| `!resume` | Resume playback | DJ+ |
| `!queue` / `!q` | Display the current queue | User+ |
| `!nowplaying` / `!np` | Show currently playing song | User+ |
| `!remove <position>` / `!rm <position>` | Remove song at position | DJ+ |
| `!clear` | Clear the entire queue | Mod+ |
| `!movetop <position>` / `!mt <position>` | Move song to top of queue | DJ+ |
| `!volume <0-100>` / `!vol <0-100>` | Set playback volume | DJ+ |

### Bot Commands

| Command | Description | Permission |
|---------|-------------|------------|
| `!join` | Join your voice channel | User+ |
| `!leave` / `!disconnect` | Leave voice channel | User+ |
| `!setrole <dj/mod> <@role>` | Set DJ or Moderator role | Admin |
| `!source` / `!info` | Show source code and creator info | User+ |
| `!help` | Show help message | User+ |

## Usage Examples

### Playing Music

```
!play https://www.youtube.com/watch?v=dQw4w9WgXcQ
!play never gonna give you up
!p https://soundcloud.com/artist/track
```

### Managing Queue

```
!queue                  # Show current queue
!remove 3               # Remove song at position 3
!movetop 5              # Move song at position 5 to top
!clear                  # Clear entire queue (Mod only)
```

### Playback Control

```
!skip                   # Skip current song
!pause                  # Pause playback
!resume                 # Resume playback
!volume 75              # Set volume to 75%
!stop                   # Stop playback (Mod only)
```

### Server Setup

```
!setrole dj @DJ         # Set DJ role
!setrole mod @Moderator # Set Moderator role
```

## Project Structure

```
miku_bot_go/
├── cmd/
│   └── bot/
│       └── main.go              # Entry point
├── internal/
│   ├── bot/
│   │   ├── bot.go               # Bot core logic
│   │   └── config.go            # Configuration loader
│   ├── commands/
│   │   └── commands.go          # Command handlers
│   ├── database/
│   │   └── database.go          # SQLite database layer
│   ├── music/
│   │   └── player.go            # Music player with DCA encoding
│   ├── permissions/
│   │   └── permissions.go       # Role-based permission system
│   └── queue/
│       └── queue.go             # Queue manager
├── configs/
│   └── config.yaml              # Bot configuration
├── .env.example                 # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Architecture

### Database Schema

**guilds**
- Stores guild-specific settings (prefix, role IDs, volume)

**queue**
- Persistent queue storage with position tracking
- Links to guild and user information

**playback_history**
- Tracks all played songs for analytics

### Music Playback Flow

1. User issues `!play` command with URL or search query
2. Bot extracts video information using yt-dlp
3. Track is added to database and in-memory queue
4. If not already playing, bot starts playback
5. yt-dlp streams audio to FFmpeg
6. DCA encodes audio for Discord
7. Audio is sent to Discord voice channel
8. On completion, next track in queue starts automatically

### Permission System

The bot uses a hierarchical permission system:
- Each guild can set custom DJ and Moderator roles
- Admins are determined by Discord server permissions
- Commands check user level before execution

## Troubleshooting

### Bot doesn't respond to commands
- Check that MESSAGE CONTENT INTENT is enabled
- Verify the bot has permission to read messages in the channel
- Ensure the correct command prefix is being used

### Audio playback issues
- Verify FFmpeg is installed: `ffmpeg -version`
- Verify yt-dlp is installed: `yt-dlp --version`
- Check bot has permission to connect and speak in voice channel

### Database errors
- Ensure GCC/build tools are installed for sqlite3
- Check file permissions for database file
- Verify database path in config.yaml

### Build errors
```bash
# If you get CGO errors, ensure GCC is installed
# On Windows, you may need to set:
set CGO_ENABLED=1

# Clear module cache and rebuild
go clean -modcache
go mod download
go build -o miku_bot ./cmd/bot
```

## Development

### Running in Development

```bash
# Run with hot reload (using air or similar)
go run ./cmd/bot

# Run with verbose logging
export LOG_LEVEL=debug
go run ./cmd/bot
```

### Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o miku_bot ./cmd/bot

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o miku_bot-linux ./cmd/bot
GOOS=windows GOARCH=amd64 go build -o miku_bot-windows.exe ./cmd/bot
GOOS=darwin GOARCH=arm64 go build -o miku_bot-macos ./cmd/bot
```

## Dependencies

- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite3 driver
- [dca](https://github.com/jonas747/dca) - Discord audio encoding
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parser

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

**What this means:**
- You are free to use, modify, and distribute this software
- If you modify and deploy this bot (even as a network service), you must:
  - Make your source code available
  - License your modifications under AGPL-3.0
  - Provide a way for users to access the source (the `!source` command does this)

See the [LICENSE](LICENSE) file for the full license text.

**Network Use:** The AGPL-3.0 license requires that if you run a modified version of this bot as a network service (like a public Discord bot), you must make the complete source code of your modified version available to users. The `!source` command is included to help satisfy this requirement.

## Creator

- **GitHub:** [blubskye](https://github.com/blubskye)
- **Discord:** blubaustin
- **Repository:** [miku_discord_music_bot](https://github.com/blubskye/miku_discord_music_bot)

## Credits

Built with Go and love for music.

Copyright (C) 2025 blubskye

## Support

For issues and feature requests, please open an issue on the [GitHub repository](https://github.com/blubskye/miku_discord_music_bot/issues).
