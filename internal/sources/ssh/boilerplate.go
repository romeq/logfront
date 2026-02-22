package ssh

import (
	"github.com/romeq/logfront/internal/domain"
	"github.com/romeq/logfront/internal/sources"
)

type Source struct {
	config Config
}

func NewSource(config sources.SourceConfigType) domain.Source {
	return &Source{config: NewConfig(config)}
}

func (s Source) Name() string {
	return ConfigName
}
