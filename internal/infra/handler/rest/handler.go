package rest

import (
	"context"
	"github.com/atlant1da-404/internal/infra/handler/semafore"
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

	semaphore interface {
		Acquire() bool
		Release()
		Usage() int
	}
)

type apiHandler struct {
	uc        Usecase
	fastJSON  fastJSON
	semaphore semaphore

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
		semaphore: semafore.NewSemaphore(200),
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go h.worker(ctx, wg)
	}

	return h
}
