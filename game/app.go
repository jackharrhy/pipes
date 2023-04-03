package game

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackharrhy/pipes/puzzles"
)

type Screen int64

const (
	LevelSelectScreen Screen = iota
	GameScreen
	CompleteScreen
)

type model struct {
	display    Screen
	levelsList list.Model
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
		display:    LevelSelectScreen,
		levelsList: list.New(puzzles.Levels, list.NewDefaultDelegate(), 0, 0),
		width:      width,
		minWidth:   32,
		height:     height,
		minHeight:  16,
		board:      cellbuffer{},
	}

	m.levelsList.Title = "pipes"

	m.levelsList.SetShowStatusBar(false)
	m.levelsList.SetShowHelp(false)

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
		m.levelsList.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd

	switch m.display {

	case LevelSelectScreen:
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				level, ok := m.levelsList.SelectedItem().(puzzles.Level)

				if ok {
					m.board.init(*level.Board())
					recalculatePipes(&m.board)
					m.display = GameScreen
				}

				return m, nil
			}

		}

		m.levelsList, cmd = m.levelsList.Update(msg)
		return m, cmd

	case GameScreen:
		switch msg := msg.(type) {

		case tea.MouseMsg:
			m.init = true
			m.mouseEvent = tea.MouseEvent(msg)

			if msg.Type == tea.MouseLeft {
				v := m.board.get(msg.X, msg.Y)
				ri, ok := runeToRuneinfo[v.rune()]
				if ok {
					p := v.(piece)
					p.state = ri.nextPiecestate
					m.board.set(p, msg.X, msg.Y)
					recalculatePipes(&m.board)
				}
			}

		}

		return m, nil

	case CompleteScreen:
		return m, nil

	}

	return m, nil
}

func (m model) View() string {
	if m.width != m.minWidth || m.height != m.minHeight {
		s := ""
		s += fmt.Sprintf("resize the window to %dx%d\n", m.minWidth, m.minHeight)
		s += fmt.Sprintf("you are at %d x %d\n", m.width, m.height)
		return s
	}

	switch m.display {
	case LevelSelectScreen:
		return m.levelsList.View()
	case GameScreen:
		return m.board.Display()
	case CompleteScreen:
		return "complete"
	}
	return ""
}
