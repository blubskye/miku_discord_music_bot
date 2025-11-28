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

package permissions

import (
	"github.com/bwmarrin/discordgo"
)

type Level int

const (
	LevelUser Level = iota
	LevelDJ
	LevelMod
	LevelAdmin
)

type Permission struct {
	djRoleID  string
	modRoleID string
}

func New(djRoleID, modRoleID string) *Permission {
	return &Permission{
		djRoleID:  djRoleID,
		modRoleID: modRoleID,
	}
}

func (p *Permission) GetUserLevel(s *discordgo.Session, guildID, userID string) (Level, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return LevelUser, err
	}

	perms, err := s.UserChannelPermissions(userID, guildID)
	if err == nil && (perms&discordgo.PermissionAdministrator != 0 || perms&discordgo.PermissionManageServer != 0) {
		return LevelAdmin, nil
	}

	for _, roleID := range member.Roles {
		if p.modRoleID != "" && roleID == p.modRoleID {
			return LevelMod, nil
		}
		if p.djRoleID != "" && roleID == p.djRoleID {
			return LevelDJ, nil
		}
	}

	return LevelUser, nil
}

func (p *Permission) HasPermission(userLevel, requiredLevel Level) bool {
	return userLevel >= requiredLevel
}

func (p *Permission) CanAddMusic(level Level) bool {
	return level >= LevelUser
}

func (p *Permission) CanRemoveMusic(level Level) bool {
	return level >= LevelDJ
}

func (p *Permission) CanMoveToTop(level Level) bool {
	return level >= LevelDJ
}

func (p *Permission) CanSkip(level Level) bool {
	return level >= LevelDJ
}

func (p *Permission) CanManageQueue(level Level) bool {
	return level >= LevelMod
}

func (p *Permission) CanChangeSettings(level Level) bool {
	return level >= LevelAdmin
}

func (p *Permission) UpdateRoles(djRoleID, modRoleID string) {
	p.djRoleID = djRoleID
	p.modRoleID = modRoleID
}
