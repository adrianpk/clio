package am

import (
	"github.com/google/uuid"
)

type Action struct {
	Path   string
	Text   string
	Style  string
	IsForm bool
}

func ListAction(i Identifiable, style string) Action {
	return Action{
		Path:  ListPath(i),
		Text:  "Back",
		Style: style,
	}
}

func EditAction(i Identifiable, id uuid.UUID, style string) Action {
	return Action{
		Path:  EditPath(i, id),
		Text:  "Edit",
		Style: style,
	}
}

func DeleteAction(i Identifiable, id uuid.UUID, style string) Action {
	return Action{
		Path:   DeletePath(i, id),
		Text:   "Delete",
		Style:  style,
		IsForm: true,
	}
}

func NewAction(url, text, style string) Action {
	return Action{
		Path:  url,
		Text:  text,
		Style: style,
	}
}
