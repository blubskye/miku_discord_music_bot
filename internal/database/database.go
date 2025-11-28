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

package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func New(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &Database{DB: db}
	if err := database.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return database, nil
}

func (d *Database) Initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS guilds (
		id TEXT PRIMARY KEY,
		prefix TEXT DEFAULT '!',
		dj_role_id TEXT,
		mod_role_id TEXT,
		volume INTEGER DEFAULT 50,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS queue (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guild_id TEXT NOT NULL,
		channel_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		duration INTEGER,
		thumbnail TEXT,
		position INTEGER NOT NULL,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS playback_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guild_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		title TEXT NOT NULL,
		url TEXT NOT NULL,
		played_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (guild_id) REFERENCES guilds(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_queue_guild_position ON queue(guild_id, position);
	CREATE INDEX IF NOT EXISTS idx_history_guild ON playback_history(guild_id);
	`

	_, err := d.DB.Exec(schema)
	return err
}

func (d *Database) Close() error {
	return d.DB.Close()
}

type Guild struct {
	ID         string
	Prefix     string
	DJRoleID   sql.NullString
	ModRoleID  sql.NullString
	Volume     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (d *Database) GetGuild(guildID string) (*Guild, error) {
	query := `SELECT id, prefix, dj_role_id, mod_role_id, volume, created_at, updated_at FROM guilds WHERE id = ?`

	var guild Guild
	err := d.DB.QueryRow(query, guildID).Scan(
		&guild.ID,
		&guild.Prefix,
		&guild.DJRoleID,
		&guild.ModRoleID,
		&guild.Volume,
		&guild.CreatedAt,
		&guild.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return d.CreateGuild(guildID)
	}

	if err != nil {
		return nil, err
	}

	return &guild, nil
}

func (d *Database) CreateGuild(guildID string) (*Guild, error) {
	query := `INSERT INTO guilds (id) VALUES (?) RETURNING id, prefix, dj_role_id, mod_role_id, volume, created_at, updated_at`

	var guild Guild
	err := d.DB.QueryRow(query, guildID).Scan(
		&guild.ID,
		&guild.Prefix,
		&guild.DJRoleID,
		&guild.ModRoleID,
		&guild.Volume,
		&guild.CreatedAt,
		&guild.UpdatedAt,
	)

	return &guild, err
}

func (d *Database) UpdateGuildRoles(guildID, djRoleID, modRoleID string) error {
	query := `UPDATE guilds SET dj_role_id = ?, mod_role_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := d.DB.Exec(query, djRoleID, modRoleID, guildID)
	return err
}

func (d *Database) UpdateGuildVolume(guildID string, volume int) error {
	query := `UPDATE guilds SET volume = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := d.DB.Exec(query, volume, guildID)
	return err
}

type QueueItem struct {
	ID        int
	GuildID   string
	ChannelID string
	UserID    string
	Title     string
	URL       string
	Duration  int
	Thumbnail string
	Position  int
	AddedAt   time.Time
}

func (d *Database) AddToQueue(item *QueueItem) error {
	query := `
		INSERT INTO queue (guild_id, channel_id, user_id, title, url, duration, thumbnail, position)
		VALUES (?, ?, ?, ?, ?, ?, ?,
			COALESCE((SELECT MAX(position) + 1 FROM queue WHERE guild_id = ?), 0)
		)
	`
	_, err := d.DB.Exec(query, item.GuildID, item.ChannelID, item.UserID, item.Title, item.URL, item.Duration, item.Thumbnail, item.GuildID)
	return err
}

func (d *Database) GetQueue(guildID string) ([]*QueueItem, error) {
	query := `SELECT id, guild_id, channel_id, user_id, title, url, duration, thumbnail, position, added_at FROM queue WHERE guild_id = ? ORDER BY position ASC`

	rows, err := d.DB.Query(query, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*QueueItem
	for rows.Next() {
		var item QueueItem
		err := rows.Scan(&item.ID, &item.GuildID, &item.ChannelID, &item.UserID, &item.Title, &item.URL, &item.Duration, &item.Thumbnail, &item.Position, &item.AddedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func (d *Database) RemoveFromQueue(guildID string, position int) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM queue WHERE guild_id = ? AND position = ?`, guildID, position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE queue SET position = position - 1 WHERE guild_id = ? AND position > ?`, guildID, position)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) ClearQueue(guildID string) error {
	_, err := d.DB.Exec(`DELETE FROM queue WHERE guild_id = ?`, guildID)
	return err
}

func (d *Database) MoveToTop(guildID string, position int) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE queue SET position = -1 WHERE guild_id = ? AND position = ?`, guildID, position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE queue SET position = position + 1 WHERE guild_id = ? AND position >= 0 AND position < ?`, guildID, position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE queue SET position = 0 WHERE guild_id = ? AND position = -1`, guildID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *Database) AddToHistory(guildID, userID, title, url string) error {
	query := `INSERT INTO playback_history (guild_id, user_id, title, url) VALUES (?, ?, ?, ?)`
	_, err := d.DB.Exec(query, guildID, userID, title, url)
	return err
}
