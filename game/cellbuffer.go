package game

import (
	"strings"
)

type cellpos struct {
	x int
	y int
}

type cell interface {
	rune() rune
	display() string
}

type cellbuffer struct {
	cells  []cell
	stride int
}

type displaycell struct {
	char rune
}

func (d displaycell) rune() rune {
	return d.char
}

func (d displaycell) display() string {
	return defaultStyle.Render(string(d.char))
}

func (c *cellbuffer) init(s string) {
	c.stride = 32 // TODO support dynamic based on string
	var chars []rune = []rune(strings.Replace(s, "\n", "", -1))

	c.cells = make([]cell, len(chars))
	for i, char := range chars {
		p, ok := runeToRuneinfo[char]

		if ok {
			piece := p.piece
			piece.state = p.piecestate
			c.cells[i] = piece
		} else if char == producerRune {
			c.cells[i] = conduit{role: Producer, on: true}
		} else if char == consumerRune {
			c.cells[i] = conduit{role: Consumer, on: false}
		} else {
			c.cells[i] = displaycell{char}
		}
	}
}

func (c cellbuffer) set(v cell, x, y int) {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return
	}
	c.cells[i] = v
}

func (c cellbuffer) get(x, y int) cell {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return nil
	}
	return c.cells[i]
}

func (c cellbuffer) width() int {
	return c.stride
}

func (c cellbuffer) height() int {
	h := len(c.cells) / c.stride
	if len(c.cells)%c.stride != 0 {
		h++
	}
	return h
}

func (c cellbuffer) Display() string {
	var b strings.Builder
	for i := 0; i < len(c.cells); i++ {
		if i > 0 && i%c.stride == 0 && i < len(c.cells)-1 {
			b.WriteRune('\n')
		}
		b.WriteString(c.cells[i].display())
	}
	return b.String()
}
