package main

import (
	"os"
	"strconv"
)

func loadEnv() {
	if CHANNEL = os.Getenv("GOTATO_CHANNEL"); CHANNEL == "" {
		panic("channel cannot be blank")
	}

	if USERNAME = os.Getenv("GOTATO_USERNAME"); USERNAME == "" {
		panic("username cannot be blank")
	}

	if CLIENT_ID = os.Getenv("GOTATO_CLIENT_ID"); CLIENT_ID == "" {
		panic("client ID cannot be blank")
	}

	if CLIENT_SECRET = os.Getenv("GOTATO_CLIENT_SECRET"); CLIENT_SECRET == "" {
		panic("client secret cannot be blank")
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
