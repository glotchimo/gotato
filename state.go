package main

import (
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
	CLIENT_IRC.Say(CHANNEL, s.Aliases[s.Holder]+" has the potato!")
}

func (s *State) Reset() {
	s = &State{
		Timer:        rand.Intn(GAME_TIMER_MAX-GAME_TIMER_MIN+1) + GAME_TIMER_MIN,
		Holder:       "",
		LastUpdate:   time.Now(),
		Participants: []string{},
		Scores:       map[string]int{},
	}
}
