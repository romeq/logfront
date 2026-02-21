package ntfy_sh

import (
	"context"
	"fmt"

	"github.com/romeq/logfront/internal/domain"
)

const ConfigName = "ntfy_sh"

type NtfyshConsumer struct {
	Webhook string
}

func (s NtfyshConsumer) Name() string {
	return ConfigName
}

func (s NtfyshConsumer) Consume(ctx context.Context, e domain.FailedLoginEvent) error {
	// format + send webhook
	fmt.Println("NtfyshConsumer got event", e)
	return nil
}
