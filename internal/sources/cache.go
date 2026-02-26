package sources

import (
	"sync"

	"github.com/romeq/logfront/internal/domain"
)

type Cache struct {
	Events map[string][]domain.LogEvent
	lock   sync.Locker
}

func NewCache() Cache {
	return Cache{
		Events: map[string][]domain.LogEvent{},
		lock:   &sync.Mutex{},
	}
}

func (c *Cache) Add(event domain.LogEvent) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Events[event.ID] = append(c.Events[event.ID], event)
}

func (c *Cache) Get() map[string][]domain.LogEvent {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Events
}

func (c *Cache) Exists(id string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	for key := range c.Events {
		if key == id {
			return true
		}
	}
	return false
}

func (c *Cache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.Events)
}

func (c *Cache) Flush() map[string][]domain.LogEvent {
	c.lock.Lock()
	defer c.lock.Unlock()
	events := c.Events
	c.Events = map[string][]domain.LogEvent{}
	return events
}
