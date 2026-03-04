package rest

import (
	"context"
	"github.com/atlant1da-404/internal/model"
	"github.com/dgraph-io/ristretto"
	"sync"
	"time"
)

type (
	Usecase interface {
		CreateNotesBatch(notes []*model.NoteCreate)
	}
)

type apiHandler struct {
	uc        Usecase
	cache     *ristretto.Cache
	flushChan chan *model.NoteCreate
	maxBatch  int
	flushDur  time.Duration
}

func NewAPIHandler(
	ctx context.Context,
	uc Usecase,
	maxBatch int,
	flushDur time.Duration,
	wg *sync.WaitGroup,
) *apiHandler {

	h := &apiHandler{
		uc:        uc,
		flushChan: make(chan *model.NoteCreate, 200_000),
		maxBatch:  maxBatch,
		flushDur:  flushDur,
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go h.worker(ctx, wg)
	}

	return h
}
