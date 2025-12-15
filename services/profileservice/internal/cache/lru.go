package cache

import (
	"bioly/profileservice/internal/types"
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

type LruProfileCache struct {
	c *lru.Cache[string, *types.Profile]
}

func NewLruProfileCache(size int) ProfileCache {
	if size <= 0 {
		size = 1
	}

	c, err := lru.New[string, *types.Profile](size)
	if err != nil {
		panic(fmt.Errorf("failed to create LRU cache: %w", err))
	}

	return &LruProfileCache{c: c}
}

func (c *LruProfileCache) GetProfile(username string) (*types.Profile, bool) {
	if username == "" {
		return nil, false
	}
	v, ok := c.c.Get(username)
	if !ok || v == nil {
		return nil, false
	}
	return v, true
}

func (c *LruProfileCache) AddProfile(profile types.Profile) error {
	if profile.Username == "" {
		return fmt.Errorf("empty username")
	}

	p := profile
	c.c.Add(p.Username, &p)

	return nil
}
