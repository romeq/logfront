package sources

import "github.com/romeq/logfront/internal/domain"

var registry = make(map[string]func(config SourceConfigType) domain.Source)

func Register(name string, factory func(config SourceConfigType) domain.Source) {
	registry[name] = factory
}

func Create(name string, config SourceConfigType) (domain.Source, bool) {
	f, ok := registry[name]
	if !ok {
		return nil, false
	}
	return f(config), true
}
