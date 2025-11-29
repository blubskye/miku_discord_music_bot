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

package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"miku_bot/internal/database"
	"miku_bot/internal/music"
	"miku_bot/internal/permissions"
	"miku_bot/internal/queue"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	db          *database.Database
	queueMgr    *queue.Manager
	permissions map[string]*permissions.Permission
	prefix      string
	library     *music.Library
}

func NewHandler(db *database.Database, queueMgr *queue.Manager, prefix string, library *music.Library) *Handler {
	return &Handler{
		db:          db,
		queueMgr:    queueMgr,
		permissions: make(map[string]*permissions.Permission),
		prefix:      prefix,
		library:     library,
	}
}

func (h *Handler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, h.prefix) {
		return
	}

	content := strings.TrimPrefix(m.Content, h.prefix)
	args := strings.Fields(content)

	if len(args) == 0 {
		return
	}

	command := strings.ToLower(args[0])
	args = args[1:]

	switch command {
	case "play", "p":
		h.handlePlay(s, m, args)
	case "skip", "s":
		h.handleSkip(s, m)
	case "stop":
		h.handleStop(s, m)
	case "pause":
		h.handlePause(s, m)
	case "resume":
		h.handleResume(s, m)
	case "queue", "q":
		h.handleQueue(s, m)
	case "nowplaying", "np":
		h.handleNowPlaying(s, m)
	case "remove", "rm":
		h.handleRemove(s, m, args)
	case "clear":
		h.handleClear(s, m)
	case "movetop", "mt":
		h.handleMoveTop(s, m, args)
	case "volume", "vol":
		h.handleVolume(s, m, args)
	case "join":
		h.handleJoin(s, m)
	case "leave", "disconnect":
		h.handleLeave(s, m)
	case "help":
		h.handleHelp(s, m)
	case "source", "src", "info":
		h.handleSource(s, m)
	case "setrole":
		h.handleSetRole(s, m, args)
	case "folders":
		h.handleFolders(s, m)
	case "files":
		h.handleFiles(s, m, args)
	case "local", "l":
		h.handleLocalPlay(s, m, args)
	case "search":
		h.handleSearch(s, m, args)
	}
}

func (h *Handler) getPermission(guildID string) *permissions.Permission {
	if perm, exists := h.permissions[guildID]; exists {
		return perm
	}

	guild, err := h.db.GetGuild(guildID)
	if err != nil {
		return permissions.New("", "")
	}

	djRole := ""
	modRole := ""
	if guild.DJRoleID.Valid {
		djRole = guild.DJRoleID.String
	}
	if guild.ModRoleID.Valid {
		modRole = guild.ModRoleID.String
	}

	perm := permissions.New(djRole, modRole)
	h.permissions[guildID] = perm
	return perm
}

func (h *Handler) getUserVoiceChannel(s *discordgo.Session, guildID, userID string) (string, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return "", err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs.ChannelID, nil
		}
	}

	return "", fmt.Errorf("user not in voice channel")
}

func (h *Handler) handlePlay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a URL or search query!")
		return
	}

	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanAddMusic(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to add music!")
		return
	}

	voiceChannel, err := h.getUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "You must be in a voice channel!")
		return
	}

	url := strings.Join(args, " ")

	if !strings.HasPrefix(url, "http") {
		url = "ytsearch:" + url
	}

	msg, _ := s.ChannelMessageSend(m.ChannelID, "Fetching track information...")

	info, err := music.ExtractInfo(url)
	if err != nil {
		s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("Error: %v", err))
		return
	}

	track := &music.Track{
		Title:     info.Title,
		URL:       info.URL,
		Duration:  info.Duration,
		Thumbnail: info.Thumbnail,
		Requester: m.Author.ID,
	}

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Connect(s, voiceChannel); err != nil {
		s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("Error connecting to voice channel: %v", err))
		return
	}

	if err := h.queueMgr.AddTrack(m.GuildID, voiceChannel, m.Author.ID, track); err != nil {
		s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("Error adding track: %v", err))
		return
	}

	s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("Added to queue: **%s**", info.Title))

	if !player.IsPlaying() {
		player.Play()
	}
}

func (h *Handler) handleSkip(s *discordgo.Session, m *discordgo.MessageCreate) {
	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanSkip(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to skip!")
		return
	}

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Skip(); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Skipped current track!")
}

func (h *Handler) handleStop(s *discordgo.Session, m *discordgo.MessageCreate) {
	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanManageQueue(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to stop playback!")
		return
	}

	player := h.queueMgr.GetPlayer(m.GuildID)
	player.Stop()
	h.queueMgr.ClearQueue(m.GuildID)

	s.ChannelMessageSend(m.ChannelID, "Stopped playback and cleared queue!")
}

func (h *Handler) handlePause(s *discordgo.Session, m *discordgo.MessageCreate) {
	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanSkip(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to pause!")
		return
	}

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Pause(); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Paused playback!")
}

func (h *Handler) handleResume(s *discordgo.Session, m *discordgo.MessageCreate) {
	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanSkip(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to resume!")
		return
	}

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Resume(); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Resumed playback!")
}

func (h *Handler) handleQueue(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := h.queueMgr.GetPlayer(m.GuildID)
	queue := player.GetQueue()

	if len(queue) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Queue is empty!")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Music Queue",
		Color: 0x9B59B6,
	}

	nowPlaying := player.NowPlaying()
	if nowPlaying != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Now Playing",
			Value:  fmt.Sprintf("**%s**\nRequested by <@%s>", nowPlaying.Title, nowPlaying.Requester),
			Inline: false,
		})
	}

	queueText := ""
	for i, track := range queue {
		if i >= 10 {
			queueText += fmt.Sprintf("\n...and %d more tracks", len(queue)-10)
			break
		}
		queueText += fmt.Sprintf("%d. **%s**\n   Requested by <@%s>\n", i+1, track.Title, track.Requester)
	}

	if queueText != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Up Next",
			Value:  queueText,
			Inline: false,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h *Handler) handleNowPlaying(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := h.queueMgr.GetPlayer(m.GuildID)
	nowPlaying := player.NowPlaying()

	if nowPlaying == nil {
		s.ChannelMessageSend(m.ChannelID, "Nothing is currently playing!")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Now Playing",
		Description: fmt.Sprintf("**%s**", nowPlaying.Title),
		Color:       0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Requested by",
				Value:  fmt.Sprintf("<@%s>", nowPlaying.Requester),
				Inline: true,
			},
		},
	}

	// Handle album art for local files
	if nowPlaying.IsLocal && strings.HasPrefix(nowPlaying.Thumbnail, "attachment://") {
		// Extract the actual file path from the library
		// We need to get the original file info to access the album art path
		if h.library != nil {
			// Find the file in the library to get the actual album art path
			allFiles := h.library.GetAllFiles()
			for _, file := range allFiles {
				if file.Path == nowPlaying.URL {
					if file.AlbumArt != "" {
						// Read the album art file
						artData, err := os.ReadFile(file.AlbumArt)
						if err == nil {
							// Send as message with embed and file attachment
							fileName := filepath.Base(file.AlbumArt)
							embed.Image = &discordgo.MessageEmbedImage{
								URL: "attachment://" + fileName,
							}

							s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
								Embed: embed,
								Files: []*discordgo.File{
									{
										Name:   fileName,
										Reader: bytes.NewReader(artData),
									},
								},
							})
							return
						}
					}
					break
				}
			}
		}
	} else if nowPlaying.Thumbnail != "" {
		// For online sources, use thumbnail URL
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: nowPlaying.Thumbnail,
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h *Handler) handleRemove(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please specify a position to remove!")
		return
	}

	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanRemoveMusic(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to remove tracks!")
		return
	}

	var position int
	fmt.Sscanf(args[0], "%d", &position)
	position--

	if err := h.queueMgr.RemoveTrack(m.GuildID, position); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Removed track at position %d", position+1))
}

func (h *Handler) handleClear(s *discordgo.Session, m *discordgo.MessageCreate) {
	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanManageQueue(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to clear the queue!")
		return
	}

	if err := h.queueMgr.ClearQueue(m.GuildID); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Queue cleared!")
}

func (h *Handler) handleMoveTop(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please specify a position to move!")
		return
	}

	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanMoveToTop(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to move tracks!")
		return
	}

	var position int
	fmt.Sscanf(args[0], "%d", &position)
	position--

	if err := h.queueMgr.MoveToTop(m.GuildID, position); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Moved track at position %d to top of queue", position+1))
}

func (h *Handler) handleVolume(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please specify a volume (0-100)!")
		return
	}

	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanSkip(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to change volume!")
		return
	}

	var volume int
	fmt.Sscanf(args[0], "%d", &volume)

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.SetVolume(volume); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	h.db.UpdateGuildVolume(m.GuildID, volume)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Volume set to %d%%", volume))
}

func (h *Handler) handleJoin(s *discordgo.Session, m *discordgo.MessageCreate) {
	voiceChannel, err := h.getUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "You must be in a voice channel!")
		return
	}

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Connect(s, voiceChannel); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Joined voice channel!")
}

func (h *Handler) handleLeave(s *discordgo.Session, m *discordgo.MessageCreate) {
	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Disconnect(); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	h.queueMgr.RemovePlayer(m.GuildID)

	s.ChannelMessageSend(m.ChannelID, "Left voice channel!")
}

func (h *Handler) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Miku Bot Help",
		Description: "Music bot with role-based permissions",
		Color:       0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Music Commands",
				Value: "`!play <url/query>` - Play a song\n" +
					"`!skip` - Skip current song (DJ+)\n" +
					"`!stop` - Stop playback (Mod+)\n" +
					"`!pause` - Pause playback (DJ+)\n" +
					"`!resume` - Resume playback (DJ+)\n" +
					"`!queue` - Show queue\n" +
					"`!nowplaying` - Show current song\n" +
					"`!remove <position>` - Remove song (DJ+)\n" +
					"`!clear` - Clear queue (Mod+)\n" +
					"`!movetop <position>` - Move song to top (DJ+)\n" +
					"`!volume <0-100>` - Set volume (DJ+)",
				Inline: false,
			},
			{
				Name: "Local Files",
				Value: "`!folders` - List all music folders\n" +
					"`!files <folder>` - List files in a folder\n" +
					"`!local <folder> <filename>` - Play local file\n" +
					"`!search <query>` - Search for files by name",
				Inline: false,
			},
			{
				Name: "Bot Commands",
				Value: "`!join` - Join voice channel\n" +
					"`!leave` - Leave voice channel\n" +
					"`!setrole <dj/mod> <@role>` - Set roles (Admin)\n" +
					"`!source` - Show source code and creator info\n" +
					"`!help` - Show this message",
				Inline: false,
			},
			{
				Name: "Supported Sources",
				Value: "YouTube, SoundCloud, Bandcamp, Vimeo, Twitch, Local files, HTTP URLs",
				Inline: false,
			},
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h *Handler) handleSource(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Miku Discord Music Bot",
		Description: "An open-source music bot for Discord",
		Color:       0x9B59B6,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Source Code",
				Value:  "[GitHub Repository](https://github.com/blubskye/miku_discord_music_bot)",
				Inline: false,
			},
			{
				Name:   "Creator",
				Value:  "**GitHub:** [blubskye](https://github.com/blubskye)\n**Discord:** blubaustin",
				Inline: false,
			},
			{
				Name:   "License",
				Value:  "GNU Affero General Public License v3.0 (AGPL-3.0)",
				Inline: false,
			},
			{
				Name:   "Version",
				Value:  "1.0.0",
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Free and Open Source Software",
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h *Handler) handleSetRole(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !setrole <dj/mod> <@role>")
		return
	}

	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanChangeSettings(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to change settings!")
		return
	}

	roleType := strings.ToLower(args[0])
	roleID := strings.Trim(args[1], "<@&>")

	guild, err := h.db.GetGuild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error getting guild settings!")
		return
	}

	djRoleID := ""
	modRoleID := ""

	if guild.DJRoleID.Valid {
		djRoleID = guild.DJRoleID.String
	}
	if guild.ModRoleID.Valid {
		modRoleID = guild.ModRoleID.String
	}

	switch roleType {
	case "dj":
		djRoleID = roleID
	case "mod", "moderator":
		modRoleID = roleID
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid role type! Use 'dj' or 'mod'")
		return
	}

	if err := h.db.UpdateGuildRoles(m.GuildID, djRoleID, modRoleID); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error updating roles!")
		return
	}

	perm.UpdateRoles(djRoleID, modRoleID)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Updated %s role to <@&%s>", roleType, roleID))
}

func (h *Handler) handleFolders(s *discordgo.Session, m *discordgo.MessageCreate) {
	if h.library == nil {
		s.ChannelMessageSend(m.ChannelID, "Local library is not configured!")
		return
	}

	folders := h.library.GetFolders()

	if len(folders) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No folders found in local library!")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Local Music Folders",
		Color: 0x9B59B6,
	}

	foldersList := ""
	for i, folder := range folders {
		files := h.library.GetFiles(folder)
		foldersList += fmt.Sprintf("%d. **%s** (%d files)\n", i+1, folder, len(files))
	}

	embed.Description = foldersList
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total: %d files", h.library.GetTotalFiles()),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h *Handler) handleFiles(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if h.library == nil {
		s.ChannelMessageSend(m.ChannelID, "Local library is not configured!")
		return
	}

	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please specify a folder name! Use `!folders` to see available folders.")
		return
	}

	folder := strings.Join(args, " ")
	files := h.library.GetFiles(folder)

	if len(files) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No files found in folder: %s", folder))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Files in: %s", folder),
		Color: 0x9B59B6,
	}

	filesList := ""
	for i, file := range files {
		if i >= 20 {
			filesList += fmt.Sprintf("\n...and %d more files", len(files)-20)
			break
		}
		filesList += fmt.Sprintf("%d. %s\n", i+1, file.Name)
	}

	embed.Description = filesList
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total: %d files", len(files)),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func (h *Handler) handleLocalPlay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if h.library == nil {
		s.ChannelMessageSend(m.ChannelID, "Local library is not configured!")
		return
	}

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `!local <folder> <filename>`")
		return
	}

	perm := h.getPermission(m.GuildID)
	userLevel, err := perm.GetUserLevel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking permissions!")
		return
	}

	if !perm.CanAddMusic(userLevel) {
		s.ChannelMessageSend(m.ChannelID, "You don't have permission to add music!")
		return
	}

	voiceChannel, err := h.getUserVoiceChannel(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "You must be in a voice channel!")
		return
	}

	folder := args[0]
	fileName := strings.Join(args[1:], " ")

	file, err := h.library.GetFileByFolderAndName(folder, fileName)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	// Use metadata title if available, otherwise use filename
	trackTitle := file.Title
	if trackTitle == "" {
		trackTitle = file.Name
	}

	// Use album art if available
	thumbnail := file.AlbumArt
	if thumbnail != "" {
		// Convert local file path to file:// URL for Discord
		thumbnail = "attachment://" + filepath.Base(thumbnail)
	}

	track := &music.Track{
		Title:     trackTitle,
		URL:       file.Path,
		Duration:  file.Duration,
		Thumbnail: thumbnail,
		Requester: m.Author.ID,
		IsLocal:   true,
	}

	player := h.queueMgr.GetPlayer(m.GuildID)

	if err := player.Connect(s, voiceChannel); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error connecting to voice channel: %v", err))
		return
	}

	if err := h.queueMgr.AddTrack(m.GuildID, voiceChannel, m.Author.ID, track); err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error adding track: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added to queue: **%s** (from %s)", file.Name, folder))

	if !player.IsPlaying() {
		player.Play()
	}
}

func (h *Handler) handleSearch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if h.library == nil {
		s.ChannelMessageSend(m.ChannelID, "Local library is not configured!")
		return
	}

	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a search query!")
		return
	}

	query := strings.Join(args, " ")
	results := h.library.SearchByName(query)

	if len(results) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No files found matching: %s", query))
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Search results for: %s", query),
		Color: 0x9B59B6,
	}

	resultsList := ""
	for i, file := range results {
		if i >= 15 {
			resultsList += fmt.Sprintf("\n...and %d more results", len(results)-15)
			break
		}
		resultsList += fmt.Sprintf("%d. **%s** (in %s)\n", i+1, file.Name, file.Folder)
	}

	embed.Description = resultsList
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: fmt.Sprintf("Total: %d results", len(results)),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
