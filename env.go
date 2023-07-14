package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	// Debug settings
	VERBOSE bool = false

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
	JOIN_DURATION     int = 10
	GAME_DURATION_MIN int = 30
	GAME_DURATION_MAX int = 60
	TIMEOUT_DURATION  int = 30
	REWARD_BASE       int = 100
	COOLDOWN_DURATION int = 120
)

func loadEnv() {
	if verbose, err := strconv.ParseBool(os.Getenv("VERBOSE")); err == nil {
		VERBOSE = verbose
	}

	if CHANNEL = os.Getenv("GOTATO_CHANNEL"); CHANNEL == "" {
		log.Fatal("channel cannot be blank")
	}

	if USERNAME = os.Getenv("GOTATO_USERNAME"); USERNAME == "" {
		log.Fatal("username cannot be blank")
	}

	if CLIENT_ID = os.Getenv("GOTATO_CLIENT_ID"); CLIENT_ID == "" {
		log.Fatal("client ID cannot be blank")
	}

	if CLIENT_SECRET = os.Getenv("GOTATO_CLIENT_SECRET"); CLIENT_SECRET == "" {
		log.Fatal("client secret cannot be blank")
	}

	if duration, err := strconv.Atoi(os.Getenv("GOTATO_JOIN_DURATION")); err == nil {
		JOIN_DURATION = duration
	}

	if duration, err := strconv.Atoi(os.Getenv("GOTATO_GAME_DURATION_MIN")); err == nil {
		GAME_DURATION_MIN = duration
	}

	if duration, err := strconv.Atoi(os.Getenv("GOTATO_GAME_DURATION_MAX")); err == nil {
		GAME_DURATION_MAX = duration
	}

	if duration, err := strconv.Atoi(os.Getenv("GOTATO_TIMEOUT_DURATION")); err == nil {
		TIMEOUT_DURATION = duration
	}

	if reward, err := strconv.Atoi(os.Getenv("GOTATO_REWARD_BASE")); err == nil {
		REWARD_BASE = reward
	}

	if duration, err := strconv.Atoi(os.Getenv("GOTATO_COOLDOWN_DURATION")); err == nil {
		COOLDOWN_DURATION = duration
	}
}
