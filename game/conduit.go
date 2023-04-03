package game

import "log"

type ConduitRole int

const (
	Producer ConduitRole = iota
	Consumer
)

type conduit struct {
	role ConduitRole
	on   bool
}

var (
	producerRune = '0'
	consumerRune = 'X'
)

func (c conduit) rune() rune {
	if c.role == Producer {
		return producerRune
	} else {
		return consumerRune
	}
}

func (c conduit) display() string {
	if c.on {
		return onStyle.Render(string(c.rune()))
	} else {
		return offStyle.Render(string(c.rune()))
	}
}

func findConduitPositions(cb *cellbuffer) []cellpos {
	var conduitPositions []cellpos

	for i, c := range cb.cells {
		_, ok := c.(conduit)
		if ok {
			conduitPositions = append(conduitPositions, cellpos{
				i % cb.stride,
				i / cb.stride,
			})
		}
	}

	return conduitPositions
}

func neighboringPieces(cb *cellbuffer, cp cellpos) (*piece, *piece, *piece, *piece) {
	var topPiece, rightPiece, downPiece, leftPiece *piece

	top := cb.get(cp.x, cp.y-1)
	if top != nil {
		p, ok := top.(piece)

		if ok {
			topPiece = &p
		}
	}

	right := cb.get(cp.x+1, cp.y)
	if right != nil {
		p, ok := right.(piece)

		if ok {
			rightPiece = &p
		}
	}

	down := cb.get(cp.x, cp.y+1)
	if down != nil {
		p, ok := down.(piece)

		if ok {
			downPiece = &p
		}
	}

	left := cb.get(cp.x-1, cp.y)
	if left != nil {
		p, ok := left.(piece)

		if ok {
			leftPiece = &p
		}
	}

	return topPiece, rightPiece, downPiece, leftPiece
}

func possiblePaths(cb *cellbuffer, cp cellpos) []cellpos {
	var paths []cellpos
	var top, right, down, left = neighboringPieces(cb, cp)
	self, isPiece := cb.get(cp.x, cp.y).(piece)

	if top != nil && (!isPiece || self.state.up) && top.state.down && !top.on {
		paths = append(paths, cellpos{cp.x, cp.y - 1})
	}

	if right != nil && (!isPiece || self.state.right) && right.state.left && !right.on {
		paths = append(paths, cellpos{cp.x + 1, cp.y})
	}

	if down != nil && (!isPiece || self.state.down) && down.state.up && !down.on {
		paths = append(paths, cellpos{cp.x, cp.y + 1})
	}

	if left != nil && (isPiece || self.state.left) && left.state.right && !left.on {
		paths = append(paths, cellpos{cp.x - 1, cp.y})
	}

	return paths
}

const MAX_CYCLE = 1000

func recalculateProducer(cb *cellbuffer, cp cellpos) {
	cycle := 0
	paths := possiblePaths(cb, cp)

	for len(paths) != 0 || cycle >= MAX_CYCLE {
		path := paths[len(paths)-1]
		paths = paths[:len(paths)-1]

		c := cb.get(path.x, path.y)
		p := c.(piece)
		p.on = true
		cb.set(p, path.x, path.y)

		newPaths := possiblePaths(cb, path)
		paths = append(paths, newPaths...)
	}

	if cycle >= MAX_CYCLE {
		log.Fatal("hit max cycle while recalculate pipes")
	}
}

func recalculateConsumer(cb *cellbuffer, cp cellpos) {
	var top, right, down, left = neighboringPieces(cb, cp)

	on := (top != nil && top.on && top.state.down) ||
		(right != nil && right.on && right.state.left) ||
		(down != nil && down.on && down.state.up) ||
		(left != nil && left.on && left.state.right)

	p := cb.get(cp.x, cp.y).(conduit)
	p.on = on
	cb.set(p, cp.x, cp.y)
}

func recalculatePipes(cb *cellbuffer) {
	for i, c := range cb.cells {
		p, ok := c.(piece)

		if ok {
			p.on = false
			cb.set(p, i%cb.stride, i/cb.stride)
		}
	}

	conduitPositions := findConduitPositions(cb)

	for _, pos := range conduitPositions {
		p, ok := cb.get(pos.x, pos.y).(conduit)
		if ok && p.role == Producer {
			recalculateProducer(cb, pos)
		}
	}

	for _, pos := range conduitPositions {
		p, ok := cb.get(pos.x, pos.y).(conduit)
		if ok && p.role == Consumer {
			recalculateConsumer(cb, pos)
		}
	}
}
