package ssh

import (
	"context"
	"fmt"
	"log"

	"github.com/romeq/logfront/internal/domain"
	"github.com/romeq/logfront/internal/sources"
)

const ConfigName = "ssh"

type Config struct {
	Systemd   bool          `yaml:"systemd"`
	Logfile   string        `yaml:"logfile"`
	Consumers []interface{} `yaml:"consumers"`
}

func NewConfig(config sources.SourceConfigType) Config {
	cfg, err := sources.ParseConfig[Config](config)
	if err != nil {
		log.Fatalln(err)
	}
	if len(cfg.Consumers) == 0 {
		log.Fatalln("No consumers defined for", ConfigName)
	}
	return cfg
}

func (s Source) Start(ctx context.Context, out domain.EventMapChannel) error {
	if s.config.Systemd {
		return fmt.Errorf("systemd not supported yet")
	}
	if s.config.Logfile != "" {
		if err := s.StartWithLogfile(ctx, out, ""); err != nil {
			return err
		}
	}
	return nil
}
