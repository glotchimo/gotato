package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
	"go.etcd.io/bbolt"
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
	TOKEN_TTL      time.Duration = 3 * time.Hour
	BROADCASTER_ID string

	// Game settings
	JOIN_TIMER     int = 10
	GAME_TIMER_MIN int = 30
	GAME_TIMER_MAX int = 60
	TIMEOUT        int = 30
	REWARD         int = 100
	COOLDOWN       int = 120

	// Twitch IRC/API clients
	CLIENT_IRC *twitch.Client
	CLIENT_API *helix.Client

	// Points database
	POINTS_DB *bbolt.DB

	// Message templates
	POINTS_MSG string = "%s has %d points"
	WIN_MSG    string = "%s held the potato for %s and wins EZ Clap +%d (now has %d)"
	LOSS_MSG   string = "%s lost to potato OMEGALUL -%s"
)

func init() {
	loadEnv()
}

func main() {
	// Open points database
	if db, err := bbolt.Open("points.db", 0666, nil); err != nil {
		log.Fatal("error opening points database: ", err)
	} else {
		POINTS_DB = db
	}
	defer POINTS_DB.Close()
	log.Println("opened points database")

	// Run OAuth flow and build IRC client
	if err := authenticate(); err != nil {
		log.Fatal("error authenticating irc: ", err)
	}
	log.Println("authenticated with Twitch")

	// Build API client
	if err := createAPIClient(); err != nil {
		log.Fatal("error authenticating api: ", err)
	}
	log.Println("created API client")

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

	// Set up token refresh timer
	authTimer := time.NewTimer(TOKEN_TTL)

	// Set up interrupt signals channel
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		// Watch for token refresh signals
		case <-authTimer.C:
			if err := refreshToken(); err != nil {
				log.Fatal("error authenticating irc: ", err)
			}
			authTimer.Reset(TOKEN_TTL)
			log.Println("token refreshed")

		// Watch for process errors
		case err := <-errors:
			log.Fatal("error received: ", err)

		// Watch for manual interrupt signals
		case sig := <-signals:
			log.Println("received", sig.String())

			CLIENT_IRC.Say(CHANNEL, "gotato disconnected")
			CLIENT_IRC.Depart(CHANNEL)
			time.Sleep(1 * time.Second)

			// Clean up and exit
			if err := CLIENT_IRC.Disconnect(); err != nil {
				log.Println(err)
			}
			close(events)
			close(errors)
			os.Exit(0)
		}
	}
}
