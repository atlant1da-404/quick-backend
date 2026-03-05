package rest

import (
	"context"
	"github.com/atlant1da-404/internal/model"
	"sync"
	"time"
)

type (
	Usecase interface {
		CreateNotesBatch(notes []*model.NoteCreate)
	}

	fastJSON interface {
		Unmarshal(buf []byte, val interface{}) error
	}
)

type apiHandler struct {
	uc       Usecase
	fastJSON fastJSON

	flushChan chan *model.NoteCreate
	maxBatch  int
	flushDur  time.Duration
}

func NewAPIHandler(
	ctx context.Context,
	fastJSON fastJSON,
	uc Usecase,
	maxBatch int,
	flushDur time.Duration,
	wg *sync.WaitGroup,
) *apiHandler {

	h := &apiHandler{
		uc:        uc,
		fastJSON:  fastJSON,
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
