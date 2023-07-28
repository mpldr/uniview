// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Event struct {
	Action    EventType `json:"action"`
	Timestamp *int      `json:"timestamp,omitempty"`
}

type EventSubmission struct {
	Action    EventType `json:"action"`
	Timestamp *int      `json:"timestamp,omitempty"`
}

type EventType string

const (
	EventTypeJoined EventType = "JOINED"
	EventTypeLeft   EventType = "LEFT"
	EventTypePlay   EventType = "PLAY"
	EventTypePause  EventType = "PAUSE"
	EventTypeJump   EventType = "JUMP"
)

var AllEventType = []EventType{
	EventTypeJoined,
	EventTypeLeft,
	EventTypePlay,
	EventTypePause,
	EventTypeJump,
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeJoined, EventTypeLeft, EventTypePlay, EventTypePause, EventTypeJump:
		return true
	}
	return false
}

func (e EventType) String() string {
	return string(e)
}

func (e *EventType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid eventType", str)
	}
	return nil
}

func (e EventType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}