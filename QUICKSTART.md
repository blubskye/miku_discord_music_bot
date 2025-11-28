# Quick Start Guide

Get Miku Bot running in 5 minutes!

## Prerequisites Check

Before starting, verify you have:

```bash
# Check Go
go version
# Should show: go version go1.21 or higher

# Check FFmpeg
ffmpeg -version
# Should show FFmpeg version info

# Check yt-dlp
yt-dlp --version
# Should show version number
```

If any are missing, install them first:

### Install FFmpeg
- **Ubuntu/Debian**: `sudo apt install ffmpeg`
- **macOS**: `brew install ffmpeg`
- **Windows**: Download from https://ffmpeg.org/download.html

### Install yt-dlp
```bash
pip install yt-dlp
# or
pip3 install yt-dlp
```

## Setup Steps

### 1. Get Your Discord Bot Token

1. Go to https://discord.com/developers/applications
2. Click "New Application" and give it a name
3. Go to the "Bot" section in the left sidebar
4. Click "Add Bot"
5. Enable these Privileged Gateway Intents:
   - **Message Content Intent** ✓
   - **Server Members Intent** ✓
6. Click "Reset Token" and copy your bot token
7. **IMPORTANT**: Keep this token secret!

### 2. Invite Bot to Your Server

1. In Developer Portal, go to "OAuth2" > "URL Generator"
2. Select scopes:
   - `bot`
3. Select bot permissions:
   - View Channels
   - Send Messages
   - Embed Links
   - Read Message History
   - Connect
   - Speak
   - Use Voice Activity
4. Copy the generated URL
5. Open URL in browser and invite bot to your server

### 3. Configure the Bot

```bash
# Navigate to project directory
cd miku_bot_go

# Create environment file
cp .env.example .env

# Edit .env and add your token
nano .env
# or
vim .env
# or use any text editor
```

**Add your token to .env:**
```
DISCORD_TOKEN=your_bot_token_here
```

### 4. Build and Run

```bash
# Install dependencies
make deps

# Build the bot
make build

# Run the bot
./bin/miku_bot
```

Or run directly without building:
```bash
make dev
```

You should see:
```
Logged in as: YourBotName#1234
Bot is now running. Press CTRL-C to exit.
```

## First Commands

Join a voice channel in your Discord server, then try:

```
!join           # Bot joins your voice channel
!play https://www.youtube.com/watch?v=dQw4w9WgXcQ
!queue          # Show queue
!nowplaying     # Show current song
!help           # Show all commands
```

## Setting Up Roles

By default, all users can use basic commands. To set up role-based permissions:

1. Create roles in your Discord server:
   - Create a "DJ" role
   - Create a "Moderator" role

2. Assign the roles to users

3. Configure the bot (requires Admin permissions):
```
!setrole dj @DJ
!setrole mod @Moderator
```

Now:
- **Regular users** can add songs
- **DJ role** can skip, remove songs, control playback
- **Moderator role** can clear queue and stop playback
- **Server admins** can configure everything

## Common Issues

### "Bot doesn't respond"
- Make sure Message Content Intent is enabled
- Check bot has permission to read/send messages
- Verify bot is online (green status)

### "Can't play audio"
- Verify you're in a voice channel
- Check bot has Connect and Speak permissions
- Ensure FFmpeg and yt-dlp are installed

### "Build errors"
- Make sure GCC/build tools are installed
- On Windows: Install MinGW or TDM-GCC
- Run `make clean` then `make build`

## Next Steps

- Read the full [README.md](README.md) for all features
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
- Customize [configs/config.yaml](configs/config.yaml) for your needs

## Support

If you encounter issues:
1. Check logs for error messages
2. Verify all prerequisites are installed
3. Make sure bot token is correct
4. Ensure bot has required permissions in Discord

## Example Session

```bash
# Build and run
$ make build
Building Miku Bot...
Build complete: bin/miku_bot

$ ./bin/miku_bot
Logged in as: MikuBot#1234
Bot is now running. Press CTRL-C to exit.
```

In Discord:
```
You: !join
Bot: Joined voice channel!

You: !play never gonna give you up
Bot: Fetching track information...
Bot: Added to queue: Rick Astley - Never Gonna Give You Up (Official Music Video)

You: !nowplaying
Bot: [Embed showing current song with thumbnail]

You: !volume 75
Bot: Volume set to 75%

You: !queue
Bot: [Embed showing queue with all songs]
```

Enjoy your music bot!
