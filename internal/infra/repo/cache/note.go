package cache

import (
	"github.com/atlant1da-404/internal/model"
)

const cost = 10

func (r *Repository) CreateNote(note *model.Note) {
	r.r.Set("storage"+note.Id, note, cost)
}
