package rest

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/model"
)

var (
	WrongRequestNoteIDErr = errors.New("wrong request note_id")
)

type createNoteRequestBody struct {
	Title string `json:"title"`
}

// notePool comments:
// Use a buffer pool to reduce memory allocations.
// Before unmarshalling, grab an item from the pool and put it back after you're done.
var notePool = sync.Pool{
	New: func() interface{} {
		return &createNoteRequestBody{}
	},
}

func (a apiHandler) createNote(ctx *fasthttp.RequestCtx) ([]byte, error) {
	note := notePool.Get().(*createNoteRequestBody)
	defer notePool.Put(note)

	err := sonic.Unmarshal(ctx.Request.Body(), note)
	if err != nil {
		return nil, fmt.Errorf("sonic.Unmarshal: %w", err)
	}

	resp, err := a.uc.CreateNote(ctx, &model.NoteCreate{Title: note.Title})
	if err != nil {
		return nil, fmt.Errorf("a.uc.CreateNote: %w", err)
	}

	data, err := sonic.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("sonic.Marshal: %w", err)
	}

	return data, nil
}

func (a apiHandler) getNote(ctx *fasthttp.RequestCtx) ([]byte, error) {
	noteID, ok := ctx.UserValue("note_id").(string)
	if !ok {
		return nil, WrongRequestNoteIDErr
	}

	note, err := a.uc.GetNote(ctx, &model.NoteGet{Id: noteID})
	if err != nil {
		return nil, fmt.Errorf("a.uc.GetNote: %w", err)
	}

	data, err := sonic.Marshal(note)
	if err != nil {
		return nil, fmt.Errorf("sonic.Marshal: %w", err)
	}

	return data, nil
}
