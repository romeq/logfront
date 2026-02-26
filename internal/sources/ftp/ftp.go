package ftp

import (
	"context"
	"fmt"
	"log"

	"github.com/romeq/logfront/internal/domain"
	"github.com/romeq/logfront/internal/sources"
)

const ConfigName = "ftp"

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
	return fmt.Errorf("ftp not implemented")
}
