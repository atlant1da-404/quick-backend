package rest

import "github.com/valyala/fasthttp"

const (
	createNote = "/note/create"
)

func (a *apiHandler) Router(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case createNote:
		a.CreateNote(ctx)
		return
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBody([]byte("Method Not Allowed"))
	}
}
