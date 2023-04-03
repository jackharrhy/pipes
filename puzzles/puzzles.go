package puzzles

import (
	_ "embed"

	"github.com/charmbracelet/bubbles/list"
)

//go:embed 0-hello-world.pipes
var ZeroHelloWorld string

//go:embed 1-corner.pipes
var OneCorner string

type Level struct {
	title string
	desc  string
	board *string
}

func (i Level) Title() string       { return i.title }
func (i Level) Board() *string      { return i.board }
func (i Level) Description() string { return i.desc }
func (i Level) FilterValue() string { return i.title }

var Levels = []list.Item{
	Level{
		title: "hello pipes",
		desc:  "learn the ropes",
		board: &ZeroHelloWorld,
	},
	Level{
		title: "corner",
		desc:  "another dimension",
		board: &OneCorner,
	},
}
