package cache

import (
	"context"
	"time"

	"github.com/bytedance/sonic"

	"github.com/atlant1da-404/internal/model"
)

const (
	// DefaultCacheTimeout is the default timeout for cache operations.
	defaultCacheTimeout = 1 * time.Second
	// CacheExpiration is the default cache expiration time.
	cacheExpiration = time.Hour * 999
)

func (r *Repository) CreateNote(ctx context.Context, note *model.Note) error {
	data, err := sonic.Marshal(note)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, defaultCacheTimeout)
	defer cancel()

	return r.dfClient.Set(ctx, note.Id, data, cacheExpiration)
}

func (r *Repository) GetNote(ctx context.Context, filter model.NoteGet) (*model.Note, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultCacheTimeout)
	defer cancel()

	data, err := r.dfClient.GetBytes(ctx, filter.Id)
	if err != nil {
		return nil, err
	}

	var note model.Note
	err = sonic.Unmarshal(data, &note)
	if err != nil {
		return nil, err
	}

	return &note, nil
}
