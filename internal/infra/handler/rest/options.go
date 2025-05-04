package rest

import "github.com/valyala/fasthttp"

const (
	createNote = "/note/create"
	getNote    = "/note/get/{note_id}"
)

func (a apiHandler) Router(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case createNote:
		a.wrapperHandler(a.createNote)
		return
	case getNote:
		a.wrapperHandler(a.getNote)
		return
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.SetBody([]byte("Method Not Allowed"))
	}
}
