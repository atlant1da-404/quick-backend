package rest

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/atlant1da-404/internal/model"
)

type (
	Usecase interface {
		CreateNote(ctx context.Context, dto *model.NoteCreate) (*model.Note, error)
		GetNote(ctx context.Context, filter model.NoteGet) (*model.Note, error)
	}
	apiHandlerType = func(ctx *fasthttp.RequestCtx) ([]byte, error)
)

type apiHandler struct {
	uc Usecase
}

func NewAPIHandler(uc Usecase) *apiHandler {
	return &apiHandler{
		uc: uc,
	}
}

func (a apiHandler) wrapperHandler(handler apiHandlerType) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		res, err := handler(ctx)
		if err != nil {
			fmt.Fprintf(ctx, "Something wrong: :(")
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(res)
	}
}
