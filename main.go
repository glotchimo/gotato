package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
	"go.etcd.io/bbolt"
)

var (
	// Twitch IRC/API clients
	CLIENT_IRC *twitch.Client
	CLIENT_API *helix.Client

	// Points database
	POINTS_DB *bbolt.DB

	// Message templates
	POINTS_MSG string = "%s has %d points 💸"
	WIN_MSG    string = "%s held the potato for %s and wins 😎 +%d (now has %d)"
	LOSS_MSG   string = "%s lost to potato 💀 -%s"
)

func init() {
	loadEnv()
}

func main() {
	// Print out greeting
	fmt.Println("🥔 Welcome to gotato, hot potato for Twitch chat.")
	fmt.Println()
	fmt.Println("🖥️  Please complete the authentication flow in your browser.")
	fmt.Println()

	// Run OAuth flow and build IRC client
	if err := authenticate(); err != nil {
		log.Fatal("error authenticating with Twitch: ", err)
	}

	fmt.Println("🚪 I'm in!")
	fmt.Println()
	fmt.Println("🚧 Now I just need to set up a few more things...")
	fmt.Println()

	// Open points database
	if db, err := bbolt.Open("points.db", 0666, nil); err != nil {
		log.Fatal("error opening points database: ", err)
	} else {
		POINTS_DB = db
	}
	defer POINTS_DB.Close()

	fmt.Println("✅ Points database loaded")

	// Build API client
	if err := createAPIClient(); err != nil {
		log.Fatal("error authenticating api: ", err)
	}

	fmt.Println("✅ API client created")

	// Initialize event and error channels
	events := make(chan Event)
	errors := make(chan error)

	// Set up token refresh timer
	authTimer := time.NewTimer(TOKEN_TTL)

	fmt.Println("✅ Auth token refresh timer set")

	// Set up interrupt signals channel
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Println("✅ Interrupt signal channel opened")

	// Launch game loop and listener concurrently
	fmt.Println()
	fmt.Println("👍 All set, see you in chat!")
	fmt.Println()

	go loop(events, errors)
	go listen(events, errors)

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
