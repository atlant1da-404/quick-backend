package usecase

import (
	"github.com/atlant1da-404/internal/model"
	"sync"
)

func (u usecase) CreateNotesBatch(notes []*model.NoteCreate) {
	if len(notes) == 0 {
		return
	}

	workers := 4
	if len(notes) < workers {
		workers = len(notes)
	}

	chunkSize := (len(notes) + workers - 1) / workers

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if start >= len(notes) {
			wg.Done() // нічого немає для цього worker
			continue
		}
		if end > len(notes) {
			end = len(notes)
		}

		go func(subNotes []*model.NoteCreate) {
			defer wg.Done()
			for _, n := range subNotes {
				u.cache.CreateNote(&model.Note{Id: n.Id, Title: n.Title})
			}
		}(notes[start:end])
	}

	wg.Wait()
}
