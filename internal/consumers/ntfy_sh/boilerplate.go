package ntfy_sh

import (
	"log"

	"github.com/romeq/logfront/internal/consumers"
	"github.com/romeq/logfront/internal/domain"
)

type NtfyshConsumer struct {
	config Config
}

func NewConsumer(config consumers.ConsumerConfigType) domain.Consumer {
	cfg, err := newConfig(config)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return &NtfyshConsumer{config: cfg}
}

func (s NtfyshConsumer) Name() string {
	return ConfigName
}
