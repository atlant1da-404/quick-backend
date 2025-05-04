package usecase

import (
	"context"

	"github.com/atlant1da-404/internal/model"
)

func (u usecase) CreateNote(ctx context.Context, dto *model.NoteCreate) (*model.Note, error) {
	note, err := model.NewNote(dto.Title)
	if err != nil {
		return nil, err
	}

	err = u.cache.CreateNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (u usecase) GetNote(ctx context.Context, filter model.NoteGet) (*model.Note, error) {
	return u.cache.GetNote(ctx, filter)
}
