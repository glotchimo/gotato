package main

type EventType string

const (
	StartEvent  EventType = "start"
	JoinEvent             = "join"
	BetEvent              = "bet"
	PassEvent             = "pass"
	PointsEvent           = "points"
	ResetEvent            = "reset"
)

type Event struct {
	Type     EventType
	UserID   string
	Username string
	Data     any
}
