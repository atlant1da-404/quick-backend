package rest

import (
	"context"
	"github.com/rs/xid"
	"github.com/valyala/fasthttp"
	"sync"

	"github.com/atlant1da-404/internal/model"
)

const (
	noteRequestMin   = 5
	noteRequestLimit = 512
	batchPoolLength  = 2000
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

	batchPool = sync.Pool{
		New: func() any {
			b := make([]*model.NoteCreate, 0, batchPoolLength)
			return &b
		},
	}
)

type createNoteRequestBody struct {
	Title string `json:"title"`
}

func (a *apiHandler) CreateNote(ctx *fasthttp.RequestCtx) error {
	if !a.semaphore.Acquire() {
		ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
		return nil
	}
	defer a.semaphore.Release()

	body := ctx.PostBody()
	if len(body) <= noteRequestMin {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	if len(body) > noteRequestLimit {
		ctx.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
		return nil
	}

	req := requestPool.Get().(*createNoteRequestBody)
	*req = createNoteRequestBody{}

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

	batchPtr := batchPool.Get().(*[]*model.NoteCreate)
	batch := *batchPtr

	defer func() {
		batch = batch[:0]
		*batchPtr = batch
		batchPool.Put(batchPtr)
	}()

	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				a.flush(batch)
			}
			return

		case note, ok := <-a.flushChan:
			if !ok {
				return
			}
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

	for i := range notes { // zero allocation in heap
		notes[i].Id = xid.New()
	}

	a.uc.CreateNotesBatch(notes)

	for i := range notes {
		n := notes[i]
		*n = model.NoteCreate{} // memclr
		notePool.Put(n)
		notes[i] = nil // flush memory
	}
}
