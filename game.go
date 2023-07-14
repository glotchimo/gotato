package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func game(events chan Event, errors chan error) {
	state := State{
		Timer:        rand.Intn(GAME_DURATION_MAX-GAME_DURATION_MIN+1) + GAME_DURATION_MIN,
		Holder:       "",
		LastUpdate:   time.Now(),
		Participants: []string{},
		Aliases:      map[string]string{},
		Scores:       map[string]int{},
		Bets:         map[string]int{},
		Reward:       REWARD_BASE,
	}

waitPhase:
	log.Println("in wait phase")

	state.Reset()
	CLIENT_IRC.Say(CHANNEL, "Type !gotato to start the game!")
	for e := range events {
		if e.Type != StartEvent {
			log.Println("non-start event received, skipping")
			continue
		}

		goto joinPhase
	}

joinPhase:
	log.Println("in join phase")

	joinTimer := time.NewTimer(time.Duration(JOIN_DURATION) * time.Second)
	CLIENT_IRC.Say(CHANNEL, "Someone started a game of hot potato! Type !join or !bet <number> to join.")
	for {
		select {
		// Watch the chat for join/bet commands
		case e := <-events:
			if e.Type == JoinEvent || e.Type == BetEvent {
				// Add user to participant list and alias map
				state.Participants = append(state.Participants, e.UserID)
				state.Aliases[e.UserID] = e.Username

				log.Println("added participant:", e.UserID)

				// Handle bets
				if e.Type == BetEvent {
					// Don't allow multiple bets
					for id := range state.Bets {
						if e.UserID == id {
							CLIENT_IRC.Say(CHANNEL, "You can't bet twice!")
							continue
						}
					}

					// Validate the bet command
					bet, ok := e.Data.(int)
					if !ok {
						continue
					}

					// Get the users existing points
					points, err := getPoints(e.UserID)
					if err != nil {
						errors <- fmt.Errorf("error getting points for bet: %w", err)
					}

					// If the user's trying to be more than they have, just use whatever's left
					if bet > points {
						bet = points
					}
					state.Reward += bet

					// Register the bet for execution at game start
					state.Bets[e.UserID] += bet

					CLIENT_IRC.Say(CHANNEL, fmt.Sprintf("Reward pool is at %d!", state.Reward))
				}
			}

		// Move forward to the game phase when the timer runs out
		case <-joinTimer.C:
			if len(state.Participants) < 2 {
				CLIENT_IRC.Say(CHANNEL, "Not enough participants 😔")
				goto waitPhase
			}

			goto gamePhase
		}
	}

gamePhase:
	log.Println("in game phase")

	// Subtract bets from totals
	for id, bet := range state.Bets {
		points, err := getPoints(id)
		if err != nil {
			errors <- fmt.Errorf("error getting points for bet: %w", err)
		}

		if err := setPoints(id, points-bet); err != nil {
			errors <- fmt.Errorf("error setting points after bet: %w", err)
		}
	}

	state.Pass()
	gameTimer := time.NewTimer(time.Duration(state.Timer) * time.Second)
	CLIENT_IRC.Say(CHANNEL, "The potato's hot, here it comes! 🥔")
	for {
		select {
		// Watch the chat for reset/pass commands
		case e := <-events:
			if e.Type == ResetEvent && e.Username == USERNAME {
				log.Println("reset received, initiating join phase")
				goto joinPhase
			} else if e.Type != "pass" || !state.IsParticipant(e.UserID) {
				log.Println("non-pass event received, skipping")
				continue
			}

			// Handle scoring and passing
			state.Scores[state.Holder] = int(time.Since(state.LastUpdate).Seconds())
			state.Pass()

		// Handle end game and start cooldown
		case <-gameTimer.C:
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
			points, err := getPoints(winner)
			if err != nil {
				errors <- fmt.Errorf("error getting points for reward: %w", err)
			}

			points += state.Reward
			if err := setPoints(winner, points); err != nil {
				errors <- fmt.Errorf("error rewarding winner: %w", err)
			}

			// Timeout the loser
			if err := timeout(state.Holder); err != nil {
				errors <- fmt.Errorf("error timing out loser: %w", err)
			}

			// Send end game message
			CLIENT_IRC.Say(CHANNEL, fmt.Sprintf(
				WIN_MSG+" | "+LOSS_MSG,
				state.Aliases[winner],
				(time.Duration(topScore)*time.Second).String(),
				state.Reward,
				state.Aliases[state.Holder],
				(time.Duration(TIMEOUT_DURATION)*time.Second).String(),
			))

			goto coolPhase
		}
	}

coolPhase:
	log.Println("in cool phase")

	coolTimer := time.NewTimer(time.Duration(COOLDOWN_DURATION) * time.Second)
	CLIENT_IRC.Say(CHANNEL, "The potato's cooling down. Use !points to check your spoils!")
	for {
		select {
		// Watch for point requests
		case e := <-events:
			if e.Type == PointsEvent {
				points, err := getPoints(e.UserID)
				if err != nil {
					errors <- err
				}

				CLIENT_IRC.Say(CHANNEL, fmt.Sprintf(POINTS_MSG, state.Aliases[e.UserID], points))
			}

		// Reset to the wait phase once the cooldown's done
		case <-coolTimer.C:
			goto waitPhase
		}
	}
}
