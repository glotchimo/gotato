package main

import (
	"log"
	"math/rand"
	"time"
)

func loop(events chan string, errors chan error) {
	// Initialize state with a random timer
	state := State{
		Timer:     rand.Intn(TIMER_MAX-TIMER_MIN+1) + TIMER_MIN,
		Holder:    "",
		LastEvent: time.Now(),
		Scores:    map[string]int{},
	}

	// Consume events and update state
	for event := range events {
		if event == "" {
			log.Println("no-op received")
			continue
		}

		// Update score
		if event != state.Holder && state.Holder != "" {
			state.Scores[state.Holder] = int(time.Since(state.LastEvent).Seconds())
		}

		// Update holder
		state.Holder = event

		state.LastEvent = time.Now()
		log.Println("state updated:", state)
	}
}
