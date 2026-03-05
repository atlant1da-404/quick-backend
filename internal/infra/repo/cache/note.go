package cache

import (
	"github.com/atlant1da-404/internal/model"
	"unsafe"
)

const cost = 1

func (r *Repository) CreateNote(note *model.Note) {
	var keyBytes [19]byte
	copy(keyBytes[:7], "storage")
	copy(keyBytes[7:], note.Id[:])

	key := unsafe.String(&keyBytes[0], len(keyBytes))

	noteCopy := *note
	r.r.Set(key, &noteCopy, cost)
}
