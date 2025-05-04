package cache

import (
	"context"
	"time"
)

type (
	DragonFlyClient interface {
		Set(ctx context.Context, key string, value any, expiration time.Duration) error
		GetBytes(ctx context.Context, key string) ([]byte, error)
	}
)

type Repository struct {
	dfClient DragonFlyClient
}

// New creates a new Dragonfly repository.
func New(dfClient DragonFlyClient) *Repository {
	return &Repository{
		dfClient: dfClient,
	}
}
