/*
 * Miku Discord Music Bot
 * Copyright (C) 2025 blubskye (https://github.com/blubskye)
 * Discord: blubaustin
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Source: https://github.com/blubskye/miku_discord_music_bot
 */

package music

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type Track struct {
	Title     string
	URL       string
	Duration  int
	Thumbnail string
	Requester string
	IsLocal   bool
}

type Player struct {
	guildID       string
	voiceConn     *discordgo.VoiceConnection
	encoding      *dca.EncodeSession
	streaming     *dca.StreamingSession
	queue         []*Track
	nowPlaying    *Track
	volume        int
	mu            sync.RWMutex
	stopChan      chan bool
	isPlaying     bool
	isPaused      bool
}

func NewPlayer(guildID string) *Player {
	return &Player{
		guildID:   guildID,
		queue:     make([]*Track, 0),
		volume:    50,
		stopChan:  make(chan bool),
		isPlaying: false,
		isPaused:  false,
	}
}

type VideoInfo struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func ExtractInfo(url string) (*VideoInfo, error) {
	// Check if it's a local file path
	if isLocalFile(url) {
		return extractLocalFileInfo(url)
	}

	args := []string{
		"--dump-json",
		"--no-playlist",
		"--format", "bestaudio",
	}

	// Add API keys if available (helps avoid rate limiting)
	if youtubeKey := os.Getenv("YOUTUBE_API_KEY"); youtubeKey != "" {
		args = append(args, "--username", "oauth2", "--password", "")
	}

	if soundcloudAuth := os.Getenv("SOUNDCLOUD_AUTH_TOKEN"); soundcloudAuth != "" {
		args = append(args, "--add-header", "Authorization:OAuth "+soundcloudAuth)
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to extract info: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %w", err)
	}

	return &info, nil
}

func isLocalFile(path string) bool {
	// Check if it's an absolute path
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

func extractLocalFileInfo(path string) (*VideoInfo, error) {
	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Get file name without extension for title
	fileName := filepath.Base(path)
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	return &VideoInfo{
		Title:     title,
		URL:       path,
		Duration:  0, // Duration extraction can be added later
		Thumbnail: "",
	}, nil
}

func (p *Player) Connect(s *discordgo.Session, channelID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.voiceConn != nil {
		return nil
	}

	vc, err := s.ChannelVoiceJoin(p.guildID, channelID, false, true)
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %w", err)
	}

	p.voiceConn = vc
	return nil
}

func (p *Player) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.voiceConn != nil {
		if err := p.voiceConn.Disconnect(); err != nil {
			return err
		}
		p.voiceConn = nil
	}

	p.Stop()
	return nil
}

func (p *Player) AddTrack(track *Track) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.queue = append(p.queue, track)
}

func (p *Player) GetQueue() []*Track {
	p.mu.RLock()
	defer p.mu.RUnlock()

	queueCopy := make([]*Track, len(p.queue))
	copy(queueCopy, p.queue)
	return queueCopy
}

func (p *Player) RemoveTrack(position int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if position < 0 || position >= len(p.queue) {
		return errors.New("invalid position")
	}

	p.queue = append(p.queue[:position], p.queue[position+1:]...)
	return nil
}

func (p *Player) MoveToTop(position int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if position < 0 || position >= len(p.queue) {
		return errors.New("invalid position")
	}

	track := p.queue[position]
	p.queue = append(p.queue[:position], p.queue[position+1:]...)
	p.queue = append([]*Track{track}, p.queue...)

	return nil
}

func (p *Player) ClearQueue() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.queue = make([]*Track, 0)
}

func (p *Player) Play() error {
	p.mu.Lock()
	if p.isPlaying {
		p.mu.Unlock()
		return nil
	}

	if len(p.queue) == 0 {
		p.mu.Unlock()
		return errors.New("queue is empty")
	}

	if p.voiceConn == nil {
		p.mu.Unlock()
		return errors.New("not connected to voice channel")
	}

	p.isPlaying = true
	p.mu.Unlock()

	go p.playLoop()
	return nil
}

func (p *Player) playLoop() {
	for {
		p.mu.Lock()
		if len(p.queue) == 0 {
			p.isPlaying = false
			p.nowPlaying = nil
			p.mu.Unlock()
			return
		}

		track := p.queue[0]
		p.queue = p.queue[1:]
		p.nowPlaying = track
		p.mu.Unlock()

		if err := p.playTrack(track); err != nil {
			fmt.Printf("Error playing track: %v\n", err)
		}

		select {
		case <-p.stopChan:
			p.mu.Lock()
			p.isPlaying = false
			p.nowPlaying = nil
			p.mu.Unlock()
			return
		default:
		}
	}
}

func (p *Player) playTrack(track *Track) error {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 128
	options.Application = "audio"
	options.Volume = p.volume

	var cmd *exec.Cmd
	var stdout io.ReadCloser
	var err error

	// Handle local files differently
	if track.IsLocal {
		// Use ffmpeg directly for local files
		encodeSession, err := dca.EncodeFile(track.URL, options)
		if err != nil {
			return fmt.Errorf("failed to encode local file: %w", err)
		}
		defer encodeSession.Cleanup()

		p.mu.Lock()
		p.encoding = encodeSession
		done := make(chan error)
		streamSession := dca.NewStream(encodeSession, p.voiceConn, done)
		p.streaming = streamSession
		p.mu.Unlock()

		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				return fmt.Errorf("streaming error: %w", err)
			}
		case <-p.stopChan:
			p.mu.Lock()
			if p.streaming != nil {
				p.streaming.SetPaused(true)
			}
			if p.encoding != nil {
				p.encoding.Cleanup()
			}
			p.mu.Unlock()
			return nil
		}

		return nil
	}

	// Handle online URLs with yt-dlp
	args := []string{
		"--format", "bestaudio",
		"--output", "-",
		"--no-playlist",
	}

	// Add API keys if available (helps avoid rate limiting)
	if youtubeKey := os.Getenv("YOUTUBE_API_KEY"); youtubeKey != "" {
		args = append(args, "--username", "oauth2", "--password", "")
	}

	if soundcloudAuth := os.Getenv("SOUNDCLOUD_AUTH_TOKEN"); soundcloudAuth != "" {
		args = append(args, "--add-header", "Authorization:OAuth "+soundcloudAuth)
	}

	args = append(args, track.URL)

	cmd = exec.Command("yt-dlp", args...)

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	encodeSession, err := dca.EncodeMem(stdout, options)
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to encode audio: %w", err)
	}
	defer encodeSession.Cleanup()

	p.mu.Lock()
	p.encoding = encodeSession
	done := make(chan error)
	streamSession := dca.NewStream(encodeSession, p.voiceConn, done)
	p.streaming = streamSession
	p.mu.Unlock()

	select {
	case err := <-done:
		if err != nil && err != io.EOF {
			return fmt.Errorf("streaming error: %w", err)
		}
	case <-p.stopChan:
		p.mu.Lock()
		if p.streaming != nil {
			p.streaming.SetPaused(true)
		}
		if p.encoding != nil {
			p.encoding.Cleanup()
		}
		p.mu.Unlock()
		return nil
	}

	cmd.Wait()
	return nil
}

func (p *Player) Skip() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isPlaying {
		return errors.New("nothing is playing")
	}

	if p.streaming != nil {
		p.streaming.SetPaused(true)
	}

	if p.encoding != nil {
		p.encoding.Cleanup()
	}

	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isPlaying {
		select {
		case p.stopChan <- true:
		default:
		}
	}

	if p.streaming != nil {
		p.streaming.SetPaused(true)
	}

	if p.encoding != nil {
		p.encoding.Cleanup()
	}

	p.isPlaying = false
	p.nowPlaying = nil
}

func (p *Player) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isPlaying {
		return errors.New("nothing is playing")
	}

	if p.streaming != nil {
		p.streaming.SetPaused(true)
		p.isPaused = true
	}

	return nil
}

func (p *Player) Resume() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isPaused {
		return errors.New("player is not paused")
	}

	if p.streaming != nil {
		p.streaming.SetPaused(false)
		p.isPaused = false
	}

	return nil
}

func (p *Player) SetVolume(volume int) error {
	if volume < 0 || volume > 100 {
		return errors.New("volume must be between 0 and 100")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.volume = volume
	if p.streaming != nil {
		p.streaming.SetPaused(true)
		p.streaming.SetPaused(false)
	}

	return nil
}

func (p *Player) NowPlaying() *Track {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.nowPlaying
}

func (p *Player) IsPlaying() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.isPlaying
}

func (p *Player) IsPaused() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.isPaused
}
