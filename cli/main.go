package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackharrhy/pipes/puzzles"
)

var (
	fourPiece     = []rune{'╬'}
	threePiece    = []rune{'╩', '╠', '╦', '╣'}
	flatTwoPiece  = []rune{'║', '═'}
	curveTwoPiece = []rune{'╚', '╔', '╗', '╝'}
)

var pieceTransition = map[rune]rune{
	fourPiece[0]: fourPiece[0],

	threePiece[0]: threePiece[1],
	threePiece[1]: threePiece[2],
	threePiece[2]: threePiece[3],
	threePiece[3]: threePiece[0],

	flatTwoPiece[0]: flatTwoPiece[1],
	flatTwoPiece[1]: flatTwoPiece[0],

	curveTwoPiece[0]: curveTwoPiece[1],
	curveTwoPiece[1]: curveTwoPiece[2],
	curveTwoPiece[2]: curveTwoPiece[3],
	curveTwoPiece[3]: curveTwoPiece[0],
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#52bf90")).
	Background(lipgloss.Color("black"))

type model struct {
	init       bool
	mouseEvent tea.MouseEvent
	width      int
	minHeight  int
	height     int
	minWidth   int
	board      cellbuffer
}

func Setup(width int, height int) (tea.Model, []tea.ProgramOption) {
	m := model{
		width:     width,
		minWidth:  32,
		height:    height,
		minHeight: 16,
		board:     cellbuffer{},
	}
	m.board.init(puzzles.TestA)
	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseAllMotion()}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.MouseMsg:
		m.init = true
		m.mouseEvent = tea.MouseEvent(msg)

		if msg.Type == tea.MouseLeft {
			v := m.board.get(msg.X, msg.Y)
			next, ok := pieceTransition[v]
			if ok {
				m.board.set(next, msg.X, msg.Y)
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	if m.width != m.minWidth || m.height != m.minHeight {
		s += fmt.Sprintf("resize the window to %dx%d\n", m.minWidth, m.minHeight)
		s += fmt.Sprintf("you are at %d x %d\n", m.width, m.height)
		return s
	}

	s += style.Render(m.board.String())

	return s
}
