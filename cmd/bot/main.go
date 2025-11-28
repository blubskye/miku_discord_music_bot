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

package main

import (
	"flag"
	"log"
	"os"
	"runtime/debug"

	"miku_bot/internal/bot"

	"github.com/joho/godotenv"
)

var (
	traceFlag   = flag.Bool("trace", false, "Enable full stack tracing")
	versionFlag = flag.Bool("version", false, "Show version information")
)

const (
	Version = "1.0.0"
	Author  = "blubskye (blubaustin)"
	Repo    = "https://github.com/blubskye/miku_discord_music_bot"
)

func main() {
	flag.Parse()

	if *versionFlag {
		log.Printf("Miku Discord Music Bot v%s\n", Version)
		log.Printf("Author: %s\n", Author)
		log.Printf("Source: %s\n", Repo)
		log.Printf("License: AGPL-3.0\n")
		os.Exit(0)
	}

	if *traceFlag {
		log.Println("Stack tracing enabled")
		debug.SetTraceback("all")
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN environment variable is required")
	}

	mikuBot, err := bot.New(token, "configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Miku Bot v%s starting...", Version)
	if err := mikuBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
