package pipeline

import (
	"context"
	"log"

	"github.com/romeq/logfront/internal/domain"
)

type Dispatcher struct {
	consumers []domain.Consumer
}

func NewDispatcher(consumers []domain.Consumer) *Dispatcher {
	return &Dispatcher{consumers: consumers}
}

func (d *Dispatcher) Run(ctx context.Context, in map[string]chan domain.LogEvent) {
	for _, inChannel := range in {
		select {
		case <-ctx.Done():
			return
		case event := <-inChannel:
			for _, c := range d.consumers {
				go func() {
					err := c.Consume(ctx, event)
					if err != nil {
						log.Printf("consumer '%s' returned error when consuming '%s': %s\n", c.Name(), event.Source, err)
						log.Println(event)
					}
				}()
			}
		}
	}
}
