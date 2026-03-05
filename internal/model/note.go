package model

import "github.com/rs/xid"

type Note struct {
	Id    xid.ID
	Title string
}

type NoteCreate struct {
	Id    xid.ID
	Title string
}

func (n *NoteCreate) Reset() {
	n.Id = xid.NilID()
	n.Title = ""
}
