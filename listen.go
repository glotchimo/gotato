package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v4"
)

func listen(events chan Event, errors chan error) {
	// Initialize client with callbacks for game calls and connection issues
	CLIENT_IRC = twitch.NewClient(USERNAME, "oauth:"+ACCESS_TOKEN)

	// Watch messages (PrivateMessage = any message in the channel)
	CLIENT_IRC.OnPrivateMessage(func(m twitch.PrivateMessage) {
		if VERBOSE {
			log.Println(m)
		}

		components := strings.Split(m.Message, " ")
		cmd := strings.TrimSpace(components[0])
		switch cmd {
		case "!enable":
			if m.User.Name == USERNAME {
				BROADCASTER_ID = m.User.ID
				log.Println("enabled timeouts")
				return
			}

		case "!gotato":
			events <- Event{Type: StartEvent}

		case "!join":
			events <- Event{Type: JoinEvent, UserID: m.User.ID, Username: m.User.Name}

		case "!bet", "!wager":
			if len(components) != 2 {
				return
			}
			value, err := strconv.Atoi(components[1])
			if err != nil {
				return
			}
			events <- Event{Type: BetEvent, UserID: m.User.ID, Username: m.User.Name, Data: value}

		case "!pass", "!toss":
			events <- Event{Type: PassEvent, UserID: m.User.ID, Username: m.User.Name}

		case "!points":
			events <- Event{Type: PointsEvent, UserID: m.User.ID, Username: m.User.Name}

		case "!reset":
			events <- Event{Type: ResetEvent, UserID: m.User.ID, Username: m.User.Name}
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
	if err := CLIENT_IRC.Connect(); err != nil {
		errors <- fmt.Errorf("error connecting to channel: %w", err)
	}
}
