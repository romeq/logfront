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
	cfg := Config{}
	var ok bool
	cfg.Systemd, ok = config["systemd"].(bool)
	if !ok {
		cfg.Systemd = false
	}
	cfg.Logfile, ok = config["logfile"].(string)
	if !ok {
		cfg.Logfile = "/var/log/auth.log"
	}
	cfg.Consumers, ok = config["consumers"].([]interface{})
	if !ok {
		log.Fatalln("Sources must have consumers defined")
	}
	return cfg
}

type Source struct {
	config Config
}

func NewSource(config sources.SourceConfigType) domain.Source {
	return &Source{config: NewConfig(config)}
}

func (s Source) Name() string {
	return ConfigName
}

func (s Source) Start(ctx context.Context, out map[string]chan domain.FailedLoginEvent) error {
	// tail file
	// parse failed login

	// send the failed login to each consumer
	for _, consumer := range s.config.Consumers {
		out[consumer.(string)] <- domain.FailedLoginEvent{
			Source: ConfigName,
		}
	}

	return fmt.Errorf("Aborted SSH")
}
