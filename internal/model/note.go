package model

import gonanoid "github.com/matoous/go-nanoid/v2"

type Note struct {
	Id    string
	Title string
}

type NoteCreate struct {
	Title string
}

type NoteGet struct {
	Id string
}

func NewNote(title string) (*Note, error) {
	id, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	return &Note{
		Id:    id,
		Title: title,
	}, nil
}
