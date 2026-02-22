package ntfy_sh

import (
	"context"
	"fmt"

	"github.com/romeq/logfront/internal/consumers"
	"github.com/romeq/logfront/internal/domain"
)

const (
	ConfigName = "ntfy_sh"
)

type Config struct {
	Urls []interface{} `yaml:"urls"`
}

// newConfig should initialize and verify the config.
func newConfig(rawConfig consumers.ConsumerConfigType) (Config, error) {
	return consumers.ParseConfig[Config](rawConfig)
}

func (s NtfyshConsumer) Consume(ctx context.Context, e domain.FailedLoginEvent) error {
	// format + send webhook
	fmt.Println("NtfyshConsumer got event", e)
	return nil
}
