package usecase

import (
	"github.com/atlant1da-404/internal/model"
)

type (
	CacheRepository interface {
		CreateNote(note *model.Note)
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
