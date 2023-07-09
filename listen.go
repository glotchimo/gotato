package main

import (
	"fmt"
	"github.com/gempir/go-twitch-irc/v4"
	"log"
)

func listen(events chan string, errors chan error) {
	// Initialize client with a callback that listens for calls and sends IDs to game loop
	client := twitch.NewClient(USERNAME, "oauth:"+ACCESS_TOKEN)
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Println("message received:", message)

		if message.Message == "!gotato" {
			events <- message.User.ID
		}
	})

	log.Println("Username: " + USERNAME + " Channel: " + CHANNEL)

	// Join channel and connect client
	client.Join(CHANNEL)
	log.Println("joined channel " + CHANNEL)
	if err := client.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
	log.Println("connection established")
}
