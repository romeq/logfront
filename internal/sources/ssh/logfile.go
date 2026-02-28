package ssh

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	"github.com/romeq/logfront/internal/domain"
	"github.com/romeq/logfront/internal/sources"
)

func (s Source) StartWithLogfile(ctx context.Context, out domain.EventMapChannel, logfile string) error {
	lf := logfile
	if lf == "" {
		lf = s.config.Logfile
	}

	return s.StartLogfileWorker(ctx, out, lf)
}

func (s Source) StartLogfileWorker(_ context.Context, out domain.EventMapChannel, logfile string) error {
	cache := sources.NewCache()

	for {
		fp, err := os.Open(logfile)
		if err != nil {
			return err
		}

		const maxCapacity int = 512 // maximum line length
		buf := make([]byte, maxCapacity)

		scanner := bufio.NewScanner(fp)
		scanner.Buffer(buf, maxCapacity)

		uncachedEventsById := make(map[string][]domain.LogEvent)
		for scanner.Scan() {
			lineText := scanner.Text()

			id, err := idForLine(lineText)
			if err != nil {
				log.Println("couldn't find id for line: ", err)
				log.Println("scanner error: ", scanner.Err())
				log.Println("if line is partial, it may be caused by corrupted logfile. \n\tlineText:", lineText)
				break
			}
			if cache.Exists(id) {
				continue
			}

			ev, processed := eventForLine(lineText)
			if !processed {
				continue
			}

			uncachedEventsById[ev.ID] = append(uncachedEventsById[ev.ID], ev)
		}

		if err := scanner.Err(); err != nil {
			_ = fp.Close()
			return err
		}

		var uncachedEvents []domain.LogEvent
		for _, events := range uncachedEventsById {
			var nonurgentEvents []domain.LogEvent
			for _, event := range events {
				// don't combine urgent and nonurgent events
				if event.Urgency > domain.URGENCY_NORMAL {
					uncachedEvents = append(uncachedEvents, event)
				} else {
					nonurgentEvents = append(nonurgentEvents, event)
				}
			}
			if len(nonurgentEvents) <= 1 {
				uncachedEvents = append(uncachedEvents, nonurgentEvents...)
				continue
			}

			// len(events) > 1
			newEvent, err := multilineEvents(nonurgentEvents)
			if err != nil {
				return err
			}

			uncachedEvents = append(uncachedEvents, newEvent)
		}

		// send new events and save to cache
		s.sendNewEvents(uncachedEvents, out, cache)

		_ = fp.Close()
		time.Sleep(time.Second * 5)
	}
}

func (s Source) sendNewEvents(uncachedEvents []domain.LogEvent, out domain.EventMapChannel, cache sources.Cache) {
	if len(uncachedEvents) == 0 {
		// nothing to do
		return
	}

	if len(uncachedEvents) == 1 {
		ev := uncachedEvents[0]
		cache.Add(ev)
		out.SendEvent(ev, s.config.Consumers)
		return
	}

	var nonurgentEvents []string
	urgentEvents := map[string]domain.LogEvent{}
	for _, ev := range uncachedEvents {
		switch ev.Urgency {
		case domain.URGENCY_NORMAL:
			nonurgentEvents = append(nonurgentEvents, ev.ProcessedMessage)

		case domain.URGENCY_HIGH,
			domain.URGENCY_CRITICAL:
			urgentEvents[ev.ProcessedMessage] = ev
		}

		cache.Add(ev)
	}

	for _, urgentEvent := range urgentEvents {
		out.SendEvent(urgentEvent, s.config.Consumers)
	}

	uniqueNonUrgentEventTypes := unique(nonurgentEvents)
	if len(nonurgentEvents) > 0 {
		out.SendEvent(domain.LogEvent{
			ID:     "grouped-notification",
			Source: ConfigName,
			Group: domain.GroupNotificationInformation{
				Count: len(nonurgentEvents),
				Types: uniqueNonUrgentEventTypes,
			},
		}, s.config.Consumers)
	}
}

func unique(s []string) []string {
	var unique []string
	m := map[string]bool{}

	for _, v := range s {
		if !m[v] {
			m[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}
