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
	"log"
	"os"
	"os/signal"
	"syscall"

	"miku_bot/internal/commands"
	"miku_bot/internal/database"
	"miku_bot/internal/queue"

	"github.com/bwmarrin/discordgo"
)

func displayASCIIArt() {
	asciiArt, err := os.ReadFile("ascii.txt")
	if err != nil {
		log.Println("Could not load ASCII art file")
		return
	}
	fmt.Println(string(asciiArt))
}

type Bot struct {
	Session  *discordgo.Session
	Config   *Config
	DB       *database.Database
	QueueMgr *queue.Manager
	Commands *commands.Handler
}

func New(token string, configPath string) (*Bot, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	db, err := database.New(config.Database.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	queueMgr := queue.NewManager(db)
	commandHandler := commands.NewHandler(db, queueMgr, config.Bot.Prefix)

	bot := &Bot{
		Session:  session,
		Config:   config,
		DB:       db,
		QueueMgr: queueMgr,
		Commands: commandHandler,
	}

	session.AddHandler(bot.ready)
	session.AddHandler(commandHandler.HandleMessage)

	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsMessageContent

	return bot, nil
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	displayASCIIArt()
	fmt.Println()
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	log.Printf("Bot is ready and serving %d guilds!", len(event.Guilds))
	fmt.Println()

	s.UpdateGameStatus(0, b.Config.Bot.Activity)
}

func (b *Bot) Start() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return b.Stop()
}

func (b *Bot) Stop() error {
	log.Println("Shutting down...")

	if err := b.Session.Close(); err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}

	if err := b.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}
