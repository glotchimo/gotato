package main

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v4"
)

func listen(events chan string, errors chan error) {
	// Initialize client with callbacks for game calls and connection issues
	CLIENT = twitch.NewClient(USERNAME, "oauth:"+ACCESS_TOKEN)

	// Watch messages (PrivateMessage = any message in the channel)
	CLIENT.OnPrivateMessage(func(m twitch.PrivateMessage) {
		if m.Message == "!gotato" {
			events <- "call:" + m.User.ID
		} else if m.Message == "!reset" {
			events <- "reset"
		} else if m.FirstMessage {
			events <- "join:" + m.User.ID
		}
	})

	// Watch for notices (i.e. login failures)
	CLIENT.OnNoticeMessage(func(m twitch.NoticeMessage) {
		errors <- fmt.Errorf(m.Message)
	})

	// Join channel and connect client (blocking)
	CLIENT.Join(CHANNEL)
	CLIENT.Say(CHANNEL, "gotato connected! Use !gotato to start a game")
	if err := CLIENT.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}
