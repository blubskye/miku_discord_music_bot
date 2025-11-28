# Deployment Guide

This guide covers deploying Miku Bot to a production server.

## Server Requirements

- **OS**: Linux (Ubuntu 20.04+ recommended)
- **RAM**: Minimum 512MB, recommended 1GB+
- **CPU**: 1+ cores
- **Storage**: 1GB+ free space
- **Network**: Stable internet connection

## Installation Steps

### 1. Install Prerequisites

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y git build-essential ffmpeg python3-pip

# Install Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install yt-dlp
sudo pip3 install yt-dlp

# Verify installations
go version
ffmpeg -version
yt-dlp --version
```

### 2. Create Bot User

```bash
# Create dedicated user for the bot
sudo useradd -r -m -s /bin/bash miku
sudo mkdir -p /opt/miku_bot
sudo chown miku:miku /opt/miku_bot
```

### 3. Deploy Bot Files

```bash
# Clone or copy your bot to the server
cd /opt/miku_bot
sudo -u miku git clone <your-repo-url> .

# Or upload files via SCP
scp -r miku_bot_go/* user@server:/opt/miku_bot/
```

### 4. Build the Bot

```bash
cd /opt/miku_bot
sudo -u miku make deps
sudo -u miku make build
```

### 5. Configure Environment

```bash
# Create .env file
sudo -u miku nano /opt/miku_bot/.env

# Add your Discord token
DISCORD_TOKEN=your_actual_token_here
```

```bash
# Edit config if needed
sudo -u miku nano /opt/miku_bot/configs/config.yaml
```

### 6. Set Up Systemd Service

```bash
# Copy service file
sudo cp /opt/miku_bot/miku_bot.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable auto-start on boot
sudo systemctl enable miku_bot

# Start the bot
sudo systemctl start miku_bot

# Check status
sudo systemctl status miku_bot
```

## Managing the Bot

### Start/Stop/Restart

```bash
# Start the bot
sudo systemctl start miku_bot

# Stop the bot
sudo systemctl stop miku_bot

# Restart the bot
sudo systemctl restart miku_bot

# Check status
sudo systemctl status miku_bot
```

### View Logs

```bash
# Follow logs in real-time
sudo journalctl -u miku_bot -f

# View last 100 lines
sudo journalctl -u miku_bot -n 100

# View logs from today
sudo journalctl -u miku_bot --since today

# View logs with timestamps
sudo journalctl -u miku_bot -o short-precise
```

### Update the Bot

```bash
# Stop the bot
sudo systemctl stop miku_bot

# Pull latest changes
cd /opt/miku_bot
sudo -u miku git pull

# Rebuild
sudo -u miku make clean
sudo -u miku make build

# Restart
sudo systemctl start miku_bot
```

## Firewall Configuration

If you're using UFW firewall:

```bash
# Allow SSH (if not already allowed)
sudo ufw allow 22/tcp

# The bot only makes outbound connections, no inbound ports needed
# Just ensure outbound connections are allowed (default)
```

## Backup

### Database Backup

```bash
# Create backup script
sudo nano /opt/miku_bot/backup.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/opt/miku_bot/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR
cp /opt/miku_bot/miku_bot.db $BACKUP_DIR/miku_bot_$DATE.db

# Keep only last 7 days of backups
find $BACKUP_DIR -name "miku_bot_*.db" -mtime +7 -delete
```

```bash
# Make executable
sudo chmod +x /opt/miku_bot/backup.sh

# Add to crontab (daily backup at 3 AM)
sudo -u miku crontab -e
```

Add this line:
```
0 3 * * * /opt/miku_bot/backup.sh
```

## Monitoring

### Resource Usage

```bash
# Check memory usage
ps aux | grep miku_bot

# Check CPU usage
top -p $(pgrep miku_bot)

# Detailed system stats
htop
```

### Health Check Script

Create `/opt/miku_bot/health_check.sh`:

```bash
#!/bin/bash

if ! systemctl is-active --quiet miku_bot; then
    echo "Miku Bot is not running! Attempting to restart..."
    systemctl start miku_bot
    # Optional: Send notification (email, Discord webhook, etc.)
fi
```

```bash
sudo chmod +x /opt/miku_bot/health_check.sh

# Add to crontab (check every 5 minutes)
sudo crontab -e
```

Add:
```
*/5 * * * * /opt/miku_bot/health_check.sh
```

## Docker Deployment (Alternative)

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o miku_bot ./cmd/bot

FROM alpine:latest

RUN apk add --no-cache ffmpeg python3 py3-pip
RUN pip3 install yt-dlp

WORKDIR /app
COPY --from=builder /app/miku_bot .
COPY configs/ ./configs/
COPY .env .env

CMD ["./miku_bot"]
```

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  miku_bot:
    build: .
    container_name: miku_bot
    restart: unless-stopped
    volumes:
      - ./miku_bot.db:/app/miku_bot.db
      - ./configs:/app/configs
    env_file:
      - .env
```

Deploy with Docker:

```bash
# Build and run
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Update
docker-compose pull
docker-compose up -d --build
```

## Performance Tuning

### For High-Traffic Bots

Edit `/etc/systemd/system/miku_bot.service`:

```ini
[Service]
# Increase file descriptor limits
LimitNOFILE=65536

# Set nice priority
Nice=-10

# Memory limits (optional)
MemoryMax=2G
```

### Go Runtime Tuning

Set environment variables in service file:

```ini
[Service]
Environment="GOMAXPROCS=2"
Environment="GOGC=100"
```

## Troubleshooting

### Bot Won't Start

```bash
# Check logs
sudo journalctl -u miku_bot -n 50 --no-pager

# Check permissions
ls -la /opt/miku_bot

# Verify token
sudo -u miku cat /opt/miku_bot/.env

# Test manually
cd /opt/miku_bot
sudo -u miku ./bin/miku_bot
```

### Audio Issues

```bash
# Verify FFmpeg
ffmpeg -version

# Verify yt-dlp
yt-dlp --version

# Update yt-dlp
sudo pip3 install --upgrade yt-dlp
```

### Database Corruption

```bash
# Stop bot
sudo systemctl stop miku_bot

# Check database integrity
sqlite3 /opt/miku_bot/miku_bot.db "PRAGMA integrity_check;"

# Restore from backup
cp /opt/miku_bot/backups/miku_bot_YYYYMMDD_HHMMSS.db /opt/miku_bot/miku_bot.db

# Start bot
sudo systemctl start miku_bot
```

## Security Best Practices

1. **Never commit .env or tokens to Git**
2. **Use a dedicated user with limited privileges**
3. **Keep system and dependencies updated**
4. **Monitor logs for suspicious activity**
5. **Use firewall to restrict access**
6. **Regular backups**
7. **Use strong passwords for server access**
8. **Enable SSH key authentication**
9. **Disable root SSH login**
10. **Keep Discord bot token secure**

## Support

For issues during deployment, check:
- System logs: `sudo journalctl -xe`
- Bot logs: `sudo journalctl -u miku_bot -f`
- Disk space: `df -h`
- Memory: `free -h`
- Process list: `ps aux | grep miku`
