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

package queue

import (
	"fmt"
	"sync"

	"miku_bot/internal/database"
	"miku_bot/internal/music"
)

type Manager struct {
	db      *database.Database
	players map[string]*music.Player
	mu      sync.RWMutex
}

func NewManager(db *database.Database) *Manager {
	return &Manager{
		db:      db,
		players: make(map[string]*music.Player),
	}
}

func (m *Manager) GetPlayer(guildID string) *music.Player {
	m.mu.Lock()
	defer m.mu.Unlock()

	if player, exists := m.players[guildID]; exists {
		return player
	}

	player := music.NewPlayer(guildID)
	m.players[guildID] = player
	return player
}

func (m *Manager) RemovePlayer(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if player, exists := m.players[guildID]; exists {
		player.Stop()
		player.Disconnect()
		delete(m.players, guildID)
	}
}

func (m *Manager) AddTrack(guildID, channelID, userID string, track *music.Track) error {
	dbItem := &database.QueueItem{
		GuildID:   guildID,
		ChannelID: channelID,
		UserID:    userID,
		Title:     track.Title,
		URL:       track.URL,
		Duration:  track.Duration,
		Thumbnail: track.Thumbnail,
	}

	if err := m.db.AddToQueue(dbItem); err != nil {
		return fmt.Errorf("failed to add track to database: %w", err)
	}

	player := m.GetPlayer(guildID)
	player.AddTrack(track)

	return nil
}

func (m *Manager) RemoveTrack(guildID string, position int) error {
	if err := m.db.RemoveFromQueue(guildID, position); err != nil {
		return fmt.Errorf("failed to remove track from database: %w", err)
	}

	player := m.GetPlayer(guildID)
	if err := player.RemoveTrack(position); err != nil {
		return fmt.Errorf("failed to remove track from player: %w", err)
	}

	return nil
}

func (m *Manager) MoveToTop(guildID string, position int) error {
	if err := m.db.MoveToTop(guildID, position); err != nil {
		return fmt.Errorf("failed to move track in database: %w", err)
	}

	player := m.GetPlayer(guildID)
	if err := player.MoveToTop(position); err != nil {
		return fmt.Errorf("failed to move track in player: %w", err)
	}

	return nil
}

func (m *Manager) ClearQueue(guildID string) error {
	if err := m.db.ClearQueue(guildID); err != nil {
		return fmt.Errorf("failed to clear queue in database: %w", err)
	}

	player := m.GetPlayer(guildID)
	player.ClearQueue()

	return nil
}

func (m *Manager) LoadQueue(guildID string) error {
	items, err := m.db.GetQueue(guildID)
	if err != nil {
		return fmt.Errorf("failed to load queue from database: %w", err)
	}

	player := m.GetPlayer(guildID)
	player.ClearQueue()

	for _, item := range items {
		track := &music.Track{
			Title:     item.Title,
			URL:       item.URL,
			Duration:  item.Duration,
			Thumbnail: item.Thumbnail,
			Requester: item.UserID,
		}
		player.AddTrack(track)
	}

	return nil
}
