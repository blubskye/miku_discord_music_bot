# 🎵 Miku Bot - Discord Music Bot 💙

> *"The world is mine!"* - Hatsune Miku ✨

A feature-rich Discord music bot written in Go with role-based permissions and SQLite persistence. Bringing the power of music to your Discord server! 🎶

## ✨ Features

### 🎼 Multi-Source Support
Stream music from anywhere! 🌐
- 📺 YouTube
- 🎧 SoundCloud
- 🎸 Bandcamp
- 📹 Vimeo
- 🎮 Twitch streams
- 💾 Local files
- 🔗 HTTP URLs

### 🎵 Supported Audio Formats
Crystal-clear audio in multiple formats! 💎
- 🎵 MP3
- 🎼 FLAC
- 🎹 WAV
- 📦 Matroska/WebM (AAC, Opus, Vorbis codecs)
- 📱 MP4/M4A (AAC codec)
- 🎶 OGG streams (Opus, Vorbis, FLAC codecs)
- 🎙️ AAC streams
- 📋 Stream playlists (M3U, PLS)

### 👥 Role-Based Permissions
Everyone has a part to play! 🎭
- 👤 **User**: Can add music to the queue
- 🎧 **DJ**: Can add music, skip tracks, remove tracks, move tracks to top, control playback (pause/resume), adjust volume
- 🛡️ **Moderator**: All DJ permissions + can stop playback and clear the entire queue
- 👑 **Admin**: All permissions + can configure bot settings and roles

### 💾 Database Persistence
Your music, your way, always saved! 💖
- 🗄️ SQLite database for persistent storage
- ⚙️ Guild-specific settings
- 🔄 Queue persistence across restarts
- 📊 Playback history tracking

### 🎵 Local Music Library
Play your own music collection! 📁
- 📂 Automatic folder scanning and indexing
- 🔍 Search files by name with fuzzy matching
- 📋 Browse by folder structure
- ⚡ Fast playback with FFmpeg direct encoding
- 🎼 Supports: MP3, FLAC, WAV, OGG, M4A, OPUS, AAC, WMA
- 🖼️ **Album art extraction** - Automatically extracts and displays album art from audio file metadata
- 🎵 **Metadata support** - Reads track title, artist, and album information from files

## 📋 Prerequisites

Before running Miku Bot, you need to install these essentials! 🔧

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

## 💿 Installation

Let's get Miku singing in your server! 🎤

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

5. **(Optional but Recommended)** Add API keys to avoid rate limiting: 🔑

   For production use or heavy usage, add service API keys to `.env`:

   ```bash
   # YouTube Data API v3 (recommended!)
   # Get from: https://console.cloud.google.com/
   YOUTUBE_API_KEY=your_youtube_api_key

   # SoundCloud (optional)
   # Get from: https://soundcloud.com/settings/applications
   SOUNDCLOUD_CLIENT_ID=your_client_id
   SOUNDCLOUD_AUTH_TOKEN=your_auth_token
   ```

   **Why add API keys?** 🤔
   - ✅ Avoid rate limiting on YouTube and other services
   - ✅ Better quality downloads
   - ✅ Faster extraction and playback
   - ✅ More reliable for high-traffic servers

   **Note:** The bot works without API keys for personal/small server use! They're optional but recommended for production.

6. Customize `configs/config.yaml` if needed:
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
     music_folder: "/path/to/your/music"  # Set this to enable local file playback

   sources:
     local: true  # Enable local file support
   ```

7. Build and run:
   ```bash
   go build -o miku_bot ./cmd/bot
   ./miku_bot
   ```

   Or run directly:
   ```bash
   go run ./cmd/bot
   ```

## 💾 Setting Up Local Music Library

Want to play your own music collection? Here's how! 🎵

1. **Organize your music folder:**
   ```
   /path/to/your/music/
   ├── Jazz/
   │   ├── Miles Davis - So What.mp3
   │   └── John Coltrane - Giant Steps.flac
   ├── Rock/
   │   ├── Led Zeppelin - Stairway to Heaven.mp3
   │   └── Pink Floyd - Comfortably Numb.wav
   ├── Classical/
   │   └── Beethoven - Symphony No 9.flac
   └── favorite_song.mp3  # Files in root appear in "root" folder
   ```

2. **Update your config:**

   Edit `configs/config.yaml`:
   ```yaml
   music:
     music_folder: "/path/to/your/music"  # Absolute path to your music directory

   sources:
     local: true  # Must be enabled
   ```

3. **Supported audio formats:**
   - 🎵 MP3
   - 🎼 FLAC (lossless)
   - 🎹 WAV (uncompressed)
   - 📦 OGG (Opus/Vorbis)
   - 📱 M4A/AAC
   - 🎶 OPUS
   - 🎙️ WMA

4. **Using local files:**
   ```
   !folders                  # See all your folders
   !files Rock               # List files in Rock folder
   !local Rock Stairway      # Play with partial name match
   !search beethoven         # Find files across all folders
   ```

**Tips:**
- 📂 The bot automatically scans subdirectories
- 🔄 Restart the bot to refresh the library after adding new files
- 🎯 Filename matching is case-insensitive and supports partial matches
- ⚡ Local files play faster than streaming (no download needed!)
- 🖼️ **Album art** is automatically extracted from MP3, FLAC, M4A, and other formats with embedded artwork
- 🎵 Metadata (title, artist, album) is read from file tags and displayed in "now playing"

## 🚀 Command Line Flags

Power up Miku with these special flags! ⚡

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

## 🎫 Bot Invite

Invite Miku to your server! 💌

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

## 🎮 Commands

Let the concert begin! 🎪

### 🎵 Music Commands

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

### 💾 Local File Commands

| Command | Description | Permission |
|---------|-------------|------------|
| `!folders` | List all music folders in your library | User+ |
| `!files <folder>` | List all files in a specific folder | User+ |
| `!local <folder> <filename>` / `!l <folder> <filename>` | Play a local file by folder and name | User+ |
| `!search <query>` | Search for local files by name | User+ |

### 🤖 Bot Commands

| Command | Description | Permission |
|---------|-------------|------------|
| `!join` | Join your voice channel | User+ |
| `!leave` / `!disconnect` | Leave voice channel | User+ |
| `!setrole <dj/mod> <@role>` | Set DJ or Moderator role | Admin |
| `!source` / `!info` | Show source code and creator info | User+ |
| `!help` | Show help message | User+ |

## 🎬 Usage Examples

Time to make some noise! 🔊

### 🎶 Playing Music

```
!play https://www.youtube.com/watch?v=dQw4w9WgXcQ
!play never gonna give you up
!p https://soundcloud.com/artist/track
```

### 📝 Managing Queue

```
!queue                  # Show current queue
!remove 3               # Remove song at position 3
!movetop 5              # Move song at position 5 to top
!clear                  # Clear entire queue (Mod only)
```

### ⏯️ Playback Control

```
!skip                   # Skip current song
!pause                  # Pause playback
!resume                 # Resume playback
!volume 75              # Set volume to 75%
!stop                   # Stop playback (Mod only)
```

### 💾 Local File Playback

```
!folders                        # List all folders in your music library
!files Jazz                     # Show all files in the "Jazz" folder
!local Jazz song.mp3            # Play song.mp3 from Jazz folder
!local Rock track               # Partial filename matching works!
!search beethoven               # Find all files with "beethoven" in the name
```

### ⚙️ Server Setup

```
!setrole dj @DJ         # Set DJ role
!setrole mod @Moderator # Set Moderator role
```

## 📁 Project Structure

Organized like a perfect setlist! 🎼

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
│   │   ├── player.go            # Music player with DCA encoding
│   │   └── library.go           # Local music library manager
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

## 🏗️ Architecture

Built for performance and reliability! 💪

### 🗄️ Database Schema

**guilds**
- Stores guild-specific settings (prefix, role IDs, volume)

**queue**
- Persistent queue storage with position tracking
- Links to guild and user information

**playback_history**
- Tracks all played songs for analytics

### 🎵 Music Playback Flow

The magic behind the music! ✨

**Online Sources (YouTube, SoundCloud, etc.):**
1. User issues `!play` command with URL or search query
2. Bot extracts video information using yt-dlp
3. Track is added to database and in-memory queue
4. If not already playing, bot starts playback
5. yt-dlp streams audio to FFmpeg
6. DCA encodes audio for Discord
7. Audio is sent to Discord voice channel
8. On completion, next track in queue starts automatically

**Local Files:**
1. User issues `!local <folder> <filename>` command
2. Bot searches local library for matching file
3. Track is added to queue with local file path
4. FFmpeg directly encodes the local file (faster!)
5. DCA encodes audio for Discord
6. Audio is sent to Discord voice channel
7. Next track plays automatically

### 🔐 Permission System

The bot uses a hierarchical permission system:
- Each guild can set custom DJ and Moderator roles
- Admins are determined by Discord server permissions
- Commands check user level before execution

## 🔧 Troubleshooting

Having trouble? Don't worry, we've got you covered! 💙

### ❌ Bot doesn't respond to commands
- Check that MESSAGE CONTENT INTENT is enabled
- Verify the bot has permission to read messages in the channel
- Ensure the correct command prefix is being used

### 🔇 Audio playback issues
- Verify FFmpeg is installed: `ffmpeg -version`
- Verify yt-dlp is installed: `yt-dlp --version`
- Check bot has permission to connect and speak in voice channel

### 💾 Database errors
- Ensure GCC/build tools are installed for sqlite3
- Check file permissions for database file
- Verify database path in config.yaml

### 📁 Local file playback issues
- Verify `music_folder` path is absolute (not relative)
- Ensure `sources.local` is set to `true` in config.yaml
- Check that the music folder and files have read permissions
- Supported formats: MP3, FLAC, WAV, OGG, M4A, OPUS, AAC, WMA
- Restart the bot after adding new files to refresh the library
- Use `!folders` to verify the library loaded correctly

### 🛠️ Build errors
```bash
# If you get CGO errors, ensure GCC is installed
# On Windows, you may need to set:
set CGO_ENABLED=1

# Clear module cache and rebuild
go clean -modcache
go mod download
go build -o miku_bot ./cmd/bot
```

## 💻 Development

Join the development! 🚀

### 🔨 Running in Development

```bash
# Run with hot reload (using air or similar)
go run ./cmd/bot

# Run with verbose logging
export LOG_LEVEL=debug
go run ./cmd/bot
```

### 📦 Building for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o miku_bot ./cmd/bot

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o miku_bot-linux ./cmd/bot
GOOS=windows GOARCH=amd64 go build -o miku_bot-windows.exe ./cmd/bot
GOOS=darwin GOARCH=arm64 go build -o miku_bot-macos ./cmd/bot
```

## 📚 Dependencies

Standing on the shoulders of giants! 🌟

- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite3 driver
- [dca](https://github.com/jonas747/dca) - Discord audio encoding
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parser
- [tag](https://github.com/dhowden/tag) - Audio metadata and album art extraction

## 📜 License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

Free as in freedom! 💝

**What this means:**
- You are free to use, modify, and distribute this software
- If you modify and deploy this bot (even as a network service), you must:
  - Make your source code available
  - License your modifications under AGPL-3.0
  - Provide a way for users to access the source (the `!source` command does this)

See the [LICENSE](LICENSE) file for the full license text.

**Network Use:** The AGPL-3.0 license requires that if you run a modified version of this bot as a network service (like a public Discord bot), you must make the complete source code of your modified version available to users. The `!source` command is included to help satisfy this requirement.

## 👨‍💻 Creator

Made with 💙 by a Miku fan!

- **GitHub:** [blubskye](https://github.com/blubskye) ⭐
- **Discord:** blubaustin 💬
- **Repository:** [miku_discord_music_bot](https://github.com/blubskye/miku_discord_music_bot) 🎵

## 💖 Credits

Built with Go and love for music! 🎶

Special thanks to Hatsune Miku for the inspiration! ✨💙

*"Music is a moral law. It gives soul to the universe, wings to the mind, flight to the imagination, and charm and gaiety to life and to everything."* 🎵

Copyright (C) 2025 blubskye

## 🆘 Support

Need help? We're here for you! 💪

For issues and feature requests, please open an issue on the [GitHub repository](https://github.com/blubskye/miku_discord_music_bot/issues).

---

<div align="center">

### ⭐ Star this repo if Miku brings music to your Discord! ⭐

Made with 💙 and lots of ☕

*Keep the music playing!* 🎵✨

</div>
