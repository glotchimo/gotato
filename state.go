package main

import "time"

type State struct {
	Timer     int
	Holder    string
	LastEvent time.Time
	Scores    map[string]int
}
