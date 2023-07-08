package main

import (
	"fmt"
	"log"

	"github.com/gempir/go-twitch-irc/v4"
)

func listen(events chan string, errors chan error) {
	// Initialize client with a callback that listens for calls and sends IDs to game loop
	client := twitch.NewClient(USERNAME, PASSWORD)
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.Message == "!gotato" {
			events <- message.User.ID
		}
	})

	// Join channel and connect client
	client.Join(CHANNEL)
	if err := client.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
	log.Println("connection established")
}
