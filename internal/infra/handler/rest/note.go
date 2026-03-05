package rest

import (
	"context"
	"github.com/rs/xid"
	"sync"

	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/model"
)

const (
	noteRequestMin   = 5
	noteRequestLimit = 512
	batchPoolLenght  = 2000
)

var (
	notePool = sync.Pool{
		New: func() any {
			return &model.NoteCreate{}
		},
	}

	requestPool = sync.Pool{
		New: func() any { return &createNoteRequestBody{} },
	}

	sem = make(chan struct{}, 10000)
)

type createNoteRequestBody struct {
	Title string `json:"title"`
}

func (c *createNoteRequestBody) Reset() {
	c.Title = ""
}

func (a *apiHandler) CreateNote(ctx *fasthttp.RequestCtx) error {
	// concurrency rate limiter
	select {
	case sem <- struct{}{}:
		defer func() { <-sem }()
	default:
		ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
		return nil
	}

	body := ctx.Request.Body()
	if len(body) <= noteRequestMin {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	if len(body) > noteRequestLimit {
		ctx.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
		return nil
	}

	req := requestPool.Get().(*createNoteRequestBody)
	req.Reset()

	if err := a.fastJSON.Unmarshal(body, req); err != nil || len(req.Title) == 0 {
		requestPool.Put(req)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	note := notePool.Get().(*model.NoteCreate)
	note.Title = req.Title

	requestPool.Put(req)

	select {
	case a.flushChan <- note:
		ctx.SetStatusCode(fasthttp.StatusOK)
	default: // backpressure
		note.Reset()
		notePool.Put(note)
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
	}

	return nil
}

func (a *apiHandler) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	batch := make([]*model.NoteCreate, 0, batchPoolLenght)

	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				a.flush(batch)
			}
			return

		case note := <-a.flushChan:
			batch = append(batch, note)
			if len(batch) >= a.maxBatch {
				a.flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (a *apiHandler) flush(notes []*model.NoteCreate) {
	if len(notes) == 0 {
		return
	}

	for i := range notes {
		notes[i].Id = xid.New().String()
	}

	a.uc.CreateNotesBatch(notes)

	for _, note := range notes {
		note.Reset()
		notePool.Put(note)
	}
}
