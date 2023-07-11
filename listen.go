package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func listen(events chan string, errors chan error) {
	// Initialize client with callbacks for game calls and connection issues
	CLIENT_IRC = twitch.NewClient(USERNAME, "oauth:"+ACCESS_TOKEN)

	// Watch messages (PrivateMessage = any message in the channel)
	CLIENT_IRC.OnPrivateMessage(func(m twitch.PrivateMessage) {
		if VERBOSE {
			log.Println(m)
		}

		switch strings.TrimSpace(m.Message) {
		// Commands
		case "!gotato":
			events <- "start"
		case "!pass":
			fallthrough
		case "!toss":
			events <- "pass:" + m.User.ID + ":" + m.User.Name
		case "!join":
			events <- "join:" + m.User.ID + ":" + m.User.Name
		case "!points":
			events <- "points:" + m.User.ID + ":" + m.User.Name
		case "!reset":
			events <- "reset:" + m.User.ID + ":" + m.User.Name

		// Get broadcaster ID for API authentication from join message
		case "gotato connected":
			if m.User.Name == USERNAME {
				BROADCASTER_ID = m.User.ID
			}
		}
	})

	// Watch for notices (i.e. login failures)
	CLIENT_IRC.OnNoticeMessage(func(m twitch.NoticeMessage) {
		if VERBOSE {
			log.Println(m)
		}

		errors <- fmt.Errorf("error in notice callback: %s", m.Message)
	})

	// Join channel and connect client (blocking)
	CLIENT_IRC.Join(CHANNEL)
	CLIENT_IRC.Say(CHANNEL, "gotato connected")
	if err := CLIENT_IRC.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}
