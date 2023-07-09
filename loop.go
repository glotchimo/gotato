package main

import (
	"log"
	"math/rand"
	"time"
)

func loop(events chan string, errors chan error) {
	state := State{
		Timer:        rand.Intn(GAME_TIMER_MAX-GAME_TIMER_MIN+1) + GAME_TIMER_MIN,
		Holder:       "",
		LastUpdate:   time.Now(),
		Participants: []string{},
		Scores:       map[string]int{},
	}

waitPhase:
	log.Println("in wait phase")
	state.Reset()

	CLIENT.Say(CHANNEL, "!gotato to start the game")
	for event := range events {
		if event != "start" {
			log.Println("non-start event received, skipping")
			continue
		}

		goto joinPhase
	}

joinPhase:
	log.Println("in join phase")
	state.Reset()

	joinPhaseDone := make(chan bool, 1)
	joinTimer := JOIN_TIMER
	go timer(joinTimer, joinPhaseDone)

	CLIENT.Say(CHANNEL, "!join to join hot potato")
	for {
		select {
		// Watch the chat for join commands
		case event := <-events:
			// Split out event type and value (invalid = no-op)
			t, v, err := deslug(event)
			if err != nil {
				log.Println("no-op received")
				continue
			}

			// Skip anything that isn't a join command
			if t != "join" {
				log.Println("non-join event received, skipping")
				continue
			}

			// Register the issuer
			state.Participants = append(state.Participants, v)
			log.Println("added participant:", v)

		// Move forward to the game phase when the timer runs out
		case done := <-joinPhaseDone:
			if done {
				if len(state.Participants) < 2 {
					CLIENT.Say(CHANNEL, "not enough participants :(")
					goto waitPhase
				}

				close(joinPhaseDone)
				goto gamePhase
			}
		}
	}

gamePhase:
	log.Println("in game phase")

	// Start the game by passing to a random player
	state.Pass()

	// Watch for subsequent resets or passes
	for event := range events {
		// Split out event type and value (invalid = no-op)
		t, v, err := deslug(event)
		if err != nil {
			log.Println("no-op received")
			log.Println()
			continue
		}

		// Handle resets, skip anything else that isn't a pass
		if t == "reset" && v == USERNAME {
			log.Println("reset received, initiating join phase")
			goto joinPhase
		} else if t != "pass" {
			log.Println("non-pass event received, skipping")
			continue
		}

		// Handle scoring and passing
		state.Scores[state.Holder] = int(time.Since(state.LastUpdate).Seconds())
		state.Pass()

		// Send update messages
		state.LastUpdate = time.Now()
		log.Println("state updated:", state)
	}

	goto waitPhase
}
