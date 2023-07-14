package main

import (
	"log"
	"math/rand"
	"time"
)

type State struct {
	Timer        int
	Holder       string
	LastUpdate   time.Time
	Participants []string
	Aliases      map[string]string
	Scores       map[string]int
	Bets         map[string]int
	Reward       int
}

func (s State) IsParticipant(id string) bool {
	for _, p := range s.Participants {
		if id == p {
			return true
		}
	}

	return false
}

func (s *State) Pass() {
	var pool []string
	for _, p := range s.Participants {
		if p != s.Holder {
			pool = append(pool, p)
		}
	}

	var selection int
	if len(pool) < 2 {
		selection = 0
	} else {
		selection = rand.Intn(len(pool))
	}

	s.Holder = pool[selection]
	s.LastUpdate = time.Now()
	CLIENT_IRC.Say(CHANNEL, s.Aliases[s.Holder]+" has the potato!")
	log.Println("passed potato:", s.Holder)
}

func (s *State) Reset() {
	s = &State{
		Timer:        rand.Intn(GAME_DURATION_MAX-GAME_DURATION_MIN+1) + GAME_DURATION_MIN,
		Holder:       "",
		LastUpdate:   time.Now(),
		Participants: []string{},
		Reward:       REWARD_BASE,
		Scores:       map[string]int{},
	}
}
