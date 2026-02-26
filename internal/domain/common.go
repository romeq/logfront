package domain

import (
	"context"
)

// Consumer -interface is an interface that all the notification services are built to.
type Consumer interface {
	Name() string
	Consume(ctx context.Context, event LogEvent) error
}

// Source .- Interface contains the methods required for reading and analyzing sources
type Source interface {
	Name() string
	Start(ctx context.Context, out EventMapChannel) error
}
