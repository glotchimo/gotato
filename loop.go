package main

import (
	"fmt"
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
		Aliases:      map[string]string{},
		Scores:       map[string]int{},
	}

waitPhase:
	log.Println("in wait phase")

	// Make sure state is clean
	state.Reset()

	CLIENT_IRC.Say(CHANNEL, "!gotato to start the game")
	for event := range events {
		if event != "start" {
			log.Println("non-start event received, skipping")
			continue
		}

		goto joinPhase
	}

joinPhase:
	log.Println("in join phase")

	joinPhaseDone := make(chan bool, 1)
	go wait(JOIN_TIMER, joinPhaseDone)

	CLIENT_IRC.Say(CHANNEL, "!join to join hot potato")
	for {
		select {
		// Watch the chat for join commands
		case event := <-events:
			// Split out event type and value (invalid = no-op)
			t, id, name, err := deslug(event)
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
			state.Participants = append(state.Participants, id)
			state.Aliases[id] = name
			log.Println("added participant:", id)

		// Move forward to the game phase when the timer runs out
		case done := <-joinPhaseDone:
			if done {
				if len(state.Participants) < 2 {
					CLIENT_IRC.Say(CHANNEL, "not enough participants :(")
					goto waitPhase
				}

				close(joinPhaseDone)
				goto gamePhase
			}
		}
	}

gamePhase:
	log.Println("in game phase")

	// Start the game by passing to a random player and starting the timer
	state.Pass()

	gamePhaseDone := make(chan bool, 1)
	go wait(state.Timer, gamePhaseDone)

	CLIENT_IRC.Say(CHANNEL, "The potato's hot, here it comes!")
	for {
		select {
		// Watch the chat for pass commands
		case event := <-events:
			// Split out event type and value (invalid = no-op)
			t, id, name, err := deslug(event)
			if err != nil {
				log.Println("no-op received")
				log.Println()
				continue
			}

			// Handle resets, skip anything else that isn't a pass
			if t == "reset" && name == USERNAME {
				log.Println("reset received, initiating join phase")
				goto joinPhase
			} else if t != "pass" || !state.IsParticipant(id) {
				log.Println("non-pass event received, skipping")
				continue
			}

			// Handle scoring and passing
			state.Scores[state.Holder] = int(time.Since(state.LastUpdate).Seconds())
			state.Pass()
			state.LastUpdate = time.Now()

		// Handle end game and start cooldown
		case done := <-gamePhaseDone:
			if done {
				// Get highest score/winner ID
				var topScore int
				var winner string
				for id, score := range state.Scores {
					if score > topScore && id != state.Holder {
						winner = id
						topScore = score
					}
				}

				// Reward the winner
				points, err := reward(winner)
				if err != nil {
					errors <- err
				}

				// Timeout the loser
				if err := timeout(state.Holder); err != nil {
					errors <- err
				}

				// Send end game message
				CLIENT_IRC.Say(CHANNEL, fmt.Sprintf(
					WIN_MSG+" | "+LOSS_MSG,
					state.Aliases[winner],
					(time.Duration(topScore)*time.Second).String(),
					REWARD,
					points,
					state.Aliases[state.Holder],
					(time.Duration(TIMEOUT)*time.Second).String(),
				))

				goto coolPhase
			}
		}
	}

coolPhase:
	log.Println("in cool phase")

	coolPhaseDone := make(chan bool, 1)
	go wait(COOLDOWN, coolPhaseDone)

	// Wait for the duration of the cooldown setting
	for {
		select {
		// Watch for point requests
		case event := <-events:
			// Handle points commands
			t, id, _, err := deslug(event)
			if err == nil && t == "points" {
				points, err := getPoints(id)
				if err != nil {
					errors <- err
				}

				CLIENT_IRC.Say(CHANNEL, fmt.Sprintf(POINTS_MSG, state.Aliases[id], points))
			}

		// Reset to the wait phase once the cooldown's done
		case done := <-coolPhaseDone:
			if done {
				goto waitPhase
			}
		}
	}
}
