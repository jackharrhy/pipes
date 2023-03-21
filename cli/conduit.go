package cli

type ConduitRole int

const (
	Producer ConduitRole = iota
	Consumer
)

type conduit struct {
	role ConduitRole
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
	if c.role == Producer {
		return onStyle.Render(string(producerRune))
	} else {
		return offStyle.Render(string(consumerRune))
	}
}
