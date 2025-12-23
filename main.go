package main

import (
	"log"
	"atlantis_calc/discord" 
)

func main() {
	// Start the bot
	if err := discord.StartDiscordBot(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
