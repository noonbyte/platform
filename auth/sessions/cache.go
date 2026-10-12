package sessions

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type SessionCache struct {
	cache *cache.Cache

	ttl time.Duration
}

func NewSessionCache(ttl time.Duration) *SessionCache {
	return &SessionCache{
		cache: cache.New(ttl, ttl*2),
		ttl:   ttl,
	}
}

func (c *SessionCache) Get(id string) (*Session, bool) {
	value, ok := c.cache.Get(id)
	if !ok {
		return nil, false
	}

	session, ok := value.(*Session)
	if !ok {
		c.cache.Delete(id)
		return nil, false
	}

	return session, true
}

func (c *SessionCache) Set(session *Session) {
	c.cache.Set(session.ID, session, c.ttl*2)
}

func (c *SessionCache) Delete(id string) {
	c.cache.Delete(id)
}
