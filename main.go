package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gempir/go-twitch-irc/v4"
)

var (
	// Authorization settings
	CHANNEL       string
	USERNAME      string
	CLIENT_ID     string
	CLIENT_SECRET string
	ACCESS_TOKEN  string
	REFRESH_TOKEN string

	// Game settings
	TIMER_MIN int = 30
	TIMER_MAX int = 120
	TIMEOUT   int = 30
	REWARD    int = 100
	COOLDOWN  int = 300

	// Globally available Twitch IRC client
	CLIENT *twitch.Client
)

func init() {
	loadEnv()
	if err := authorize(); err != nil {
		log.Fatal("error authenticating:", err)
	}
}

func main() {
	// Summarize game settings
	log.Println("playing gotato with the following settings:")
	log.Println("  channel:", CHANNEL)
	log.Println("  minimum time:", TIMER_MIN)
	log.Println("  maximum time:", TIMER_MAX)
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
		log.Fatal("error received:", err)
	case sig := <-signals:
		log.Println("received", sig.String())

		// Clean up the channels and exit
		close(events)
		close(errors)
		os.Exit(0)
	}
}
