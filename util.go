package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
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

	if timer, err := strconv.Atoi(os.Getenv("GOTATO_JOIN_TIMER")); err == nil {
		JOIN_TIMER = timer
	}

	if timerMin, err := strconv.Atoi(os.Getenv("GOTATO_GAME_TIMER_MIN")); err == nil {
		GAME_TIMER_MIN = timerMin
	}

	if timerMax, err := strconv.Atoi(os.Getenv("GOTATO_GAME_TIMER_MAX")); err == nil {
		GAME_TIMER_MAX = timerMax
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

func deslug(s string) (string, string, error) {
	var event, value string
	split := strings.Split(s, ":")
	if len(split) != 2 {
		return event, value, errors.New("invalid slug")
	}

	event = split[0]
	value = split[1]
	return event, value, nil
}

func timer(t int, done chan bool) {
	for {
		time.Sleep(1 * time.Second)
		t--
		if t == 0 {
			done <- true
		}
	}
}
