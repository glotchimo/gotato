package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
)

var (
	// Debug settings
	VERBOSE bool

	// Authorization settings
	CHANNEL        string
	USERNAME       string
	CLIENT_ID      string
	CLIENT_SECRET  string
	ACCESS_TOKEN   string
	REFRESH_TOKEN  string
	BROADCASTER_ID string

	// Game settings
	JOIN_TIMER     int = 10
	GAME_TIMER_MIN int = 30
	GAME_TIMER_MAX int = 60
	TIMEOUT        int = 30
	REWARD         int = 100
	COOLDOWN       int = 120

	// Globally available Twitch IRC/API clients
	CLIENT_IRC *twitch.Client
	CLIENT_API *helix.Client

	// Message templates
	WIN_MSG  string = "%s held the potato for %s and wins EZ Clap +%d"
	LOSS_MSG string = "%s lost to potato OMEGALUL -%s"
)

func init() {
	loadEnv()

	if err := authIRC(); err != nil {
		log.Fatal("error authenticating irc: ", err)
	}

	if err := authAPI(); err != nil {
		log.Fatal("error authenticating api: ", err)
	}
}

func main() {
	// Summarize game settings
	log.Println("playing gotato with the following settings:")
	log.Println("  channel:", CHANNEL)
	log.Println("  join time:", JOIN_TIMER)
	log.Println("  minimum game time:", GAME_TIMER_MIN)
	log.Println("  maximum game time:", GAME_TIMER_MAX)
	log.Println("  loss timeout:", TIMEOUT)
	log.Println("  win reward:", REWARD)
	log.Println("  cooldown between games:", COOLDOWN)
	log.Println()

	// Initialize event and error channels
	events := make(chan string)
	errors := make(chan error)

	// Launch game loop and listener concurrently
	log.Println("launching game loop and listener")
	go loop(events, errors)
	go listen(events, errors)

	// Send a no-op to verify loop aliveness
	log.Println("sending no-op")
	events <- ""

	// Wait for errors or interrupt signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	select {
	case err := <-errors:
		log.Fatal("error received: ", err)
	case sig := <-signals:
		log.Println("received", sig.String())

		CLIENT_IRC.Say(CHANNEL, "gotato disconnected")
		time.Sleep(1 * time.Second)

		// Clean up the channels and exit
		close(events)
		close(errors)
		os.Exit(0)
	}
}
