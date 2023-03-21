package game

type piecestate struct {
	off   rune
	on    rune
	up    bool
	right bool
	down  bool
	left  bool
}

type piece struct {
	states []piecestate
	state  piecestate
	on     bool
}

func (p piece) rune() rune {
	if p.on {
		return p.state.on
	} else {
		return p.state.off
	}
}

func (p piece) display() string {
	if p.on {
		return onStyle.Render(string(p.state.on))
	} else {
		return offStyle.Render(string(p.state.off))
	}
}

var (
	fourPieceStates = []piecestate{
		piecestate{off: '╬', on: '╋'},
	}
	threePieceStates = []piecestate{
		piecestate{off: '╣', on: '┫'},
		piecestate{off: '╩', on: '┻'},
		piecestate{off: '╠', on: '┣'},
		piecestate{off: '╦', on: '┳'},
	}
	flatTwoPieceStates = []piecestate{
		piecestate{off: '═', on: '━'},
		piecestate{off: '║', on: '┃'},
	}
	curveTwoPieceStates = []piecestate{
		piecestate{off: '╝', on: '┛'},
		piecestate{off: '╚', on: '┗'},
		piecestate{off: '╔', on: '┏'},
		piecestate{off: '╗', on: '┓'},
	}
)

var (
	fourPiece = piece{
		states: fourPieceStates,
		state:  fourPieceStates[0],
		on:     false,
	}
	threePiece = piece{
		states: threePieceStates,
		state:  threePieceStates[0],
		on:     false,
	}
	flatTwoPiece = piece{
		states: flatTwoPieceStates,
		state:  flatTwoPieceStates[0],
		on:     false,
	}
	curveTwoPiece = piece{
		states: curveTwoPieceStates,
		state:  curveTwoPieceStates[0],
		on:     false,
	}
)

var pieces = []piece{fourPiece, threePiece, flatTwoPiece, curveTwoPiece}

type runeinfo struct {
	piece          piece
	piecestate     piecestate
	nextPiecestate piecestate
	on             bool
}

func generateRuneToRuneinfo() map[rune]runeinfo {
	runeToRuneinfo := make(map[rune]runeinfo)

	for _, p := range pieces {
		for i, s := range p.states {
			var nextPiecestate piecestate

			if i == len(p.states)-1 {
				nextPiecestate = p.states[0]
			} else {
				nextPiecestate = p.states[i+1]
			}

			runeToRuneinfo[s.off] = runeinfo{
				piece:          p,
				piecestate:     s,
				nextPiecestate: nextPiecestate,
				on:             false,
			}

			runeToRuneinfo[s.on] = runeinfo{
				piece:          p,
				piecestate:     s,
				nextPiecestate: nextPiecestate,
				on:             true,
			}
		}
	}

	return runeToRuneinfo
}

var runeToRuneinfo = generateRuneToRuneinfo()
