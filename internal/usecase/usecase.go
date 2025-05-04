package usecase

import (
	"context"

	"github.com/atlant1da-404/internal/model"
)

type (
	CacheRepository interface {
		CreateNote(ctx context.Context, note *model.Note) error
		GetNote(ctx context.Context, noteID model.NoteGet) (*model.Note, error)
	}
)

type usecase struct {
	cache CacheRepository
}

func NewUsecase(cache CacheRepository) *usecase {
	return &usecase{
		cache: cache,
	}
}
