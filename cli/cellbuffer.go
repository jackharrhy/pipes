package cli

import "strings"

type cellbuffer struct {
	cells  []rune
	stride int
}

func (c *cellbuffer) init(s string) {
	c.stride = 32 // TODO support dynamic based on string
	print(c.stride)
	c.cells = []rune(strings.Replace(s, "\n", "", -1))
}

func (c cellbuffer) set(v rune, x, y int) {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return
	}
	c.cells[i] = v
}

func (c cellbuffer) get(x, y int) rune {
	i := y*c.stride + x
	if i > len(c.cells)-1 || x < 0 || y < 0 || x >= c.width() || y >= c.height() {
		return ' '
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

func (c cellbuffer) String() string {
	var b strings.Builder
	for i := 0; i < len(c.cells); i++ {
		if i > 0 && i%c.stride == 0 && i < len(c.cells)-1 {
			b.WriteRune('\n')
		}
		b.WriteString(string(c.cells[i]))
	}
	return b.String()
}
