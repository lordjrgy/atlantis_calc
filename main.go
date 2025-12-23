package main

import (
	"log"
	"pkd-bot/discord" // import your bot package
)

func main() {
	// Start the bot
	if err := discord.StartDiscordBot(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
