package rest

import (
	"context"
	"github.com/rs/xid"
	"sync"

	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/model"
)

const (
	noteRequestLimit = 512
	batchPoolLenght  = 2000
)

var (
	notePool = sync.Pool{
		New: func() any {
			return &model.NoteCreate{}
		},
	}

	batchPool = sync.Pool{
		New: func() any {
			b := make([]*model.NoteCreate, 0, batchPoolLenght)
			return &b
		},
	}

	requestPool = sync.Pool{
		New: func() any { return &createNoteRequestBody{} },
	}
)

type createNoteRequestBody struct {
	Title string `json:"title"`
}

func (c *createNoteRequestBody) Reset() {
	c.Title = ""
}

func (a *apiHandler) CreateNote(ctx *fasthttp.RequestCtx) error {
	body := ctx.Request.Body()
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
	default:
		note.Reset()
		notePool.Put(note)
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
	}

	return nil
}

func (a *apiHandler) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	batch := batchPool.Get().(*[]*model.NoteCreate)
	*batch = (*batch)[:0]

	for {
		select {
		case <-ctx.Done():
			if len(*batch) > 0 {
				a.flush(*batch)
				*batch = (*batch)[:0]
			}
			batchPool.Put(batch)
			return

		case note := <-a.flushChan:
			*batch = append(*batch, note)

			if len(*batch) >= a.maxBatch {
				a.flush(*batch)
				*batch = (*batch)[:0]
			}
		}
	}
}

func (a *apiHandler) flush(notes []*model.NoteCreate) {
	for i := range notes {
		notes[i].Id = xid.New().String()
	}

	a.uc.CreateNotesBatch(notes)

	for _, note := range notes {
		note.Reset()
		notePool.Put(note)
	}
}
