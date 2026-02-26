package domain

import (
	"fmt"
	"log"
	"time"
)

type Event interface {
	Validate() error
	ID() string
}

const (
	URGENCY_NORMAL   = 0
	URGENCY_HIGH     = 1
	URGENCY_CRITICAL = 2
)

// GroupNotificationInformation allows a source to provide
// relevant information for a notification that is sent instead of
// multiple events
type GroupNotificationInformation struct {
	Count int
	Types []string
}

// EventInformation has information about
type EventInformation struct {
	Username   string
	IP         string
	Timestamp  time.Time
	RawMessage string
}
type LogEvent struct {
	ID               string
	Source           string // ssh, ftp, etc
	Urgency          uint8  // 0=normal 1=high 2=critical
	ProcessedMessage string // (shorthand)
	EventInformation EventInformation
	Group            GroupNotificationInformation
}

func (s LogEvent) Validate() error {
	log.Printf("validating event from '%s' id=%s", s.Source, s.ID)

	if s.Source == "" {
		return fmt.Errorf("source is required")
	}

	// make an exception if group data is defined
	if s.Group.Count > 0 && len(s.Group.Types) > 0 {
		for _, t := range s.Group.Types {
			if t == "" {
				return fmt.Errorf("event type cannot be empty when sending grouped notification")
			}
		}
		return nil
	}

	if s.ID == "" {
		return fmt.Errorf("id is required")
	}

	if s.EventInformation.Username == "" {
		return fmt.Errorf("username is required")
	}

	if s.EventInformation.IP == "" {
		return fmt.Errorf("ip is required")
	}

	if s.EventInformation.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	if s.EventInformation.RawMessage == "" {
		return fmt.Errorf("raw message is required")
	}

	return nil
}

type EventMapChannel map[string]chan LogEvent

func (e EventMapChannel) SendEvent(ev LogEvent, keys []interface{}) {
	log.Printf("sending event from '%s' to '%s'\n", ev.Source, keys)
	for _, k := range keys {
		ks, ok := k.(string)
		if !ok {
			log.Println("warn(SendEvent): invalid key type, skipping:", k)
			continue
		}
		e[ks] <- ev
	}
}
