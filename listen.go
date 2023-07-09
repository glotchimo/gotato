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
			events <- "start"
		} else if m.Message == "!pass" {
			events <- "pass:" + m.User.ID
		} else if m.Message == "!join" {
			events <- "join:" + m.User.ID
		} else if m.Message == "!reset" {
			events <- "reset:" + m.User.Name
		}
	})

	// Watch for notices (i.e. login failures)
	CLIENT.OnNoticeMessage(func(m twitch.NoticeMessage) {
		errors <- fmt.Errorf(m.Message)
	})

	// Watch for timeouts
	CLIENT.OnClearChatMessage(func(m twitch.ClearChatMessage) {
		return
	})

	// Join channel and connect client (blocking)
	CLIENT.Join(CHANNEL)
	CLIENT.Say(CHANNEL, "gotato connected")
	if err := CLIENT.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}
