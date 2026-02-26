package ssh

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/romeq/logfront/internal/domain"
)

func idForLine(l string) string {
	sliced := strings.Split(l, " ")
	id := strings.Trim(sliced[4], "shd-eion[]:")
	return id
}

// TODO: return a more detailed report
func getSingleMessageForEvents(events []domain.LogEvent) (string, error) {
	return events[0].ProcessedMessage, nil
}

// multilineEvents returns a populated, single log event when the event is multiple lines
func multilineEvents(events []domain.LogEvent) (domain.LogEvent, error) {
	var lastId string
	for _, event := range events {
		if lastId != "" && lastId != event.ID {
			return domain.LogEvent{}, fmt.Errorf("event ids don't match: %s != %s", lastId, event.ID)
		}
		lastId = event.ID
	}

	// get an accurate shorthand
	shorthand, err := getSingleMessageForEvents(events)
	if err != nil {
		return domain.LogEvent{}, err
	}

	var event domain.LogEvent
	event.ID = lastId
	event.Source = ConfigName
	event.EventInformation = domain.EventInformation{}
	event.EventInformation.IP = events[0].EventInformation.IP
	event.EventInformation.Username = events[0].EventInformation.Username
	event.EventInformation.Timestamp = events[0].EventInformation.Timestamp
	event.ProcessedMessage = shorthand

	return event, nil
}

// eventForLine returns a populated struct and a bool indicating whether the struct was populated
func eventForLine(l string) (domain.LogEvent, bool) {
	sliced := strings.Split(l, " ")
	text := sliced[5:]
	id := idForLine(l)

	timestamp := strings.Join(sliced[:3], " ")
	parsedTimestamp, err := time.Parse("Jan 2 15:04:05", timestamp)
	if err != nil {
		log.Println(err)
	}

	if !strings.Contains(l, "sshd-session") {
		return domain.LogEvent{}, false
	}

	if strings.Contains(l, "Connection closed by authenticating user") {
		username := text[5]
		ip := text[6]
		shorthand := fmt.Sprintf("invalid login (%s)", username)

		return NewSSHLogEvent(id, username, ip, shorthand, text, parsedTimestamp), true
	}

	if strings.Contains(l, "Connection closed by invalid user") {
		username := text[5]
		ip := text[6]
		shorthand := fmt.Sprintf("invalid user (%s)", username)

		return NewSSHLogEvent(id, username, ip, shorthand, text, parsedTimestamp), true
	}

	if strings.Contains(l, "Failed publickey for") {
		username := text[3]
		ip := text[5]
		shorthand := fmt.Sprintf("failed key (%s)", username)

		return NewSSHLogEvent(id, username, ip, shorthand, text, parsedTimestamp), true
	}

	if strings.Contains(l, "Unable to negotiate with") && strings.Contains(l, "no matching key exchange method found. Their offer: diffie-hellman-group1-sha1,diffie-hellman-group14-sha1,diffie-hellman-group14-sha256,diffie-hellman-group16-sha512,diffie-hellman-group-exchange-sha1,diffie-hellman-group-exchange-sha256 [preauth]") {
		ip := text[5]
		shorthand := "potential recon (key exchange)"
		return NewSSHLogEvent(id, "NOUID-recon", ip, shorthand, text, parsedTimestamp), true
	}

	// TODO: support other config errors as well
	if strings.Contains(l, "error: AuthorizedKeysCommand path is not absolute") {
		shorthand := "invalid config (AuthorizedKeysCommand)"

		sshEvent := NewSSHLogEvent(id, "NOUID-invalid-configuration", "localhost", shorthand, text, parsedTimestamp)
		sshEvent.Urgency = domain.URGENCY_HIGH
		return sshEvent, true
	}

	return domain.LogEvent{}, false
}

func NewSSHLogEvent(id, username, ip, shorthand string, text []string, parsedTimestamp time.Time) domain.LogEvent {
	return domain.LogEvent{
		ID:               id,
		Source:           ConfigName,
		ProcessedMessage: shorthand,
		EventInformation: domain.EventInformation{
			Username:   username,
			IP:         ip,
			Timestamp:  parsedTimestamp,
			RawMessage: strings.Join(text, " "),
		},
	}
}
