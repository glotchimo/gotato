package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v4"
)

func listen(events chan string, errors chan error) {
	// Initialize client with a callback that listens for calls and sends IDs to game loop
	client := twitch.NewClient(USERNAME, "oauth:"+ACCESS_TOKEN)

	// Watch for calls (PrivateMessage = any message in the channel)
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		if m.Message == "!gotato" {
			events <- m.User.ID
		}
	})

	// Watch for notices (i.e. login failures)
	client.OnNoticeMessage(func(m twitch.NoticeMessage) {
		errors <- fmt.Errorf(m.Message)
	})

	// Join channel and connect client (blocking)
	client.Join(CHANNEL)
	client.Say(CHANNEL, "gotato connected! Use !gotato to start a game")
	if err := client.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}
