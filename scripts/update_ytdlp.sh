#!/bin/bash
# Update yt-dlp to latest version

echo "Updating yt-dlp..."

if command -v pip3 &> /dev/null; then
    pip3 install --upgrade yt-dlp
elif command -v pip &> /dev/null; then
    pip install --upgrade yt-dlp
else
    echo "pip not found. Please install Python pip first."
    exit 1
fi

echo "yt-dlp updated successfully!"
yt-dlp --version
