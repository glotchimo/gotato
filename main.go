package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
)

var (
	CHANNEL   string = ""
	USERNAME  string = ""
	PASSWORD  string = ""
	TIMER_MIN int    = 30
	TIMER_MAX int    = 120
	TIMEOUT   int    = 30
	REWARD    int    = 100
	COOLDOWN  int    = 300
)

func init() {
	if CHANNEL = os.Getenv("GOTATO_CHANNEL"); CHANNEL == "" {
		panic("channel cannot be blank")
	}

	if USERNAME = os.Getenv("GOTATO_USERNAME"); USERNAME == "" {
		panic("username cannot be blank")
	}

	if PASSWORD := os.Getenv("GOTATO_PASSWORD"); PASSWORD == "" {
		panic("password cannot be blank")
	}

	if timerMin, err := strconv.Atoi(os.Getenv("GOTATO_TIMER_MIN")); err == nil {
		TIMER_MIN = timerMin
	}

	if timerMax, err := strconv.Atoi(os.Getenv("GOTATO_TIMER_MAX")); err == nil {
		TIMER_MAX = timerMax
	}

	if timeout, err := strconv.Atoi(os.Getenv("GOTATO_TIMEOUT")); err == nil {
		TIMEOUT = timeout
	}

	if reward, err := strconv.Atoi(os.Getenv("GOTATO_REWARD")); err == nil {
		REWARD = reward
	}

	if cooldown, err := strconv.Atoi(os.Getenv("GOTATO_COOLDOWN")); err == nil {
		COOLDOWN = cooldown
	}
}

func main() {
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
