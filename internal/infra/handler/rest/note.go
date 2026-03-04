package rest

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/model"
)

type createNoteRequestBody struct {
	Title string `json:"title"`
}

const noteRequestLimit = 512

func (a *apiHandler) CreateNote(ctx *fasthttp.RequestCtx) error {
	body := ctx.Request.Body()
	if len(body) > noteRequestLimit {
		ctx.SetStatusCode(fasthttp.StatusRequestEntityTooLarge)
		return nil
	}

	var req createNoteRequestBody
	if err := sonic.Unmarshal(body, &req); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	if strings.TrimSpace(req.Title) == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}

	id, err := gonanoid.New()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return nil
	}

	select {
	case a.flushChan <- &model.NoteCreate{Id: id, Title: req.Title}:
	default:
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

func (a *apiHandler) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	batch := make([]*model.NoteCreate, 0, a.maxBatch)
	ticker := time.NewTicker(a.flushDur)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			for {
				select {
				case note := <-a.flushChan:
					batch = append(batch, note)
					if len(batch) >= a.maxBatch {
						a.flush(batch)
						batch = batch[:0]
					}
				default:
					if len(batch) > 0 {
						a.flush(batch)
					}
					return
				}
			}

		case note := <-a.flushChan:
			batch = append(batch, note)

			if len(batch) >= a.maxBatch {
				a.flush(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				a.flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (a *apiHandler) flush(notes []*model.NoteCreate) {
	a.uc.CreateNotesBatch(notes)
}
