package ntfy_sh

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
	cfg, err := consumers.ParseConfig[Config](rawConfig)
	if err != nil {
		return cfg, err
	}

	if len(cfg.Urls) == 0 {
		return cfg, fmt.Errorf("no urls configured")
	}

	return cfg, err
}

func (s NtfyshConsumer) Consume(ctx context.Context, e domain.LogEvent) error {
	if err := e.Validate(); err != nil {
		return fmt.Errorf("failed to validate event: %w", err)
	}

	for _, url := range s.config.Urls {
		if url == nil {
			continue
		}
		if _, ok := url.(string); !ok {
			return fmt.Errorf("invalid url type: %T", url)
		}

		// validations ok, can convert type
		var url = url.(string)
		if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
			url = "https://" + url
		}

		var message string
		if e.Group.Count > 0 && len(e.Group.Types) > 0 {
			joinedEventTypes := strings.Join(e.Group.Types, ", ")
			if len(e.Group.Types) > 3 {
				joinedEventTypes = strings.Join(e.Group.Types[:3], ", ") + " and more"
			}

			message = fmt.Sprintf("%d new events were found: %s", e.Group.Count, joinedEventTypes)
		} else {
			message = fmt.Sprintf("new event from %s: %s", e.EventInformation.IP, e.ProcessedMessage)
		}

		// format + send webhook
		req, _ := http.NewRequest("POST", url, strings.NewReader(message))
		if ctx != nil {
			req = req.WithContext(ctx)
		}
		req.Header.Set("Title", fmt.Sprintf("New notification from %s", e.Source))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("unexpected status code %d", res.StatusCode)
		}
	}

	return nil
}
