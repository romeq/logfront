package consumers

import "github.com/romeq/logfront/internal/domain"

var registry = make(map[string]func(config ConsumerConfigType) domain.Consumer)

func Register(name string, factory func(config ConsumerConfigType) domain.Consumer) {
	registry[name] = factory
}

func Create(name string, config ConsumerConfigType) (domain.Consumer, bool) {
	f, ok := registry[name]
	if !ok {
		return nil, false
	}
	return f(config), true
}
