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

package bot

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		Prefix   string `yaml:"prefix"`
		Activity string `yaml:"activity"`
		Status   string `yaml:"status"`
	} `yaml:"bot"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`

	Roles struct {
		DJ  string `yaml:"dj"`
		Mod string `yaml:"mod"`
	} `yaml:"roles"`

	Music struct {
		MaxQueueSize  int    `yaml:"max_queue_size"`
		DefaultVolume int    `yaml:"default_volume"`
		Timeout       int    `yaml:"timeout"`
		MusicFolder   string `yaml:"music_folder"`
	} `yaml:"music"`

	Sources struct {
		YouTube    bool `yaml:"youtube"`
		SoundCloud bool `yaml:"soundcloud"`
		Bandcamp   bool `yaml:"bandcamp"`
		Vimeo      bool `yaml:"vimeo"`
		Twitch     bool `yaml:"twitch"`
		Local      bool `yaml:"local"`
		HTTP       bool `yaml:"http"`
	} `yaml:"sources"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
