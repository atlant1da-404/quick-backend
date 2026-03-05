package model

type Note struct {
	Id    string
	Title string
}

type NoteCreate struct {
	Id    string
	Title string
}

func (n *NoteCreate) Reset() {
	n.Id = ""
	n.Title = ""
}
