package usecase

import (
	"github.com/atlant1da-404/internal/model"
)

func (u usecase) CreateNotesBatch(notes []*model.NoteCreate) {
	if len(notes) == 0 {
		return
	}

	for _, n := range notes {
		u.cache.CreateNote(&model.Note{
			Id:    n.Id,
			Title: n.Title,
		})
	}
}
