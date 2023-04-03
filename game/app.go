package game

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackharrhy/pipes/puzzles"
	"golang.org/x/exp/slices"
)

type Screen int64

const (
	LevelSelectScreen Screen = iota
	GameScreen
	CompleteScreen
)

type model struct {
	display      Screen
	levelsList   list.Model
	init         bool
	mouseEvent   tea.MouseEvent
	width        int
	minHeight    int
	height       int
	minWidth     int
	interactive  bool
	board        cellbuffer
	currentLevel *puzzles.Level
	nextLevel    *puzzles.Level
}

func Setup(width int, height int) (tea.Model, []tea.ProgramOption) {
	m := model{
		display:      LevelSelectScreen,
		levelsList:   list.New(puzzles.Levels, list.NewDefaultDelegate(), 0, 0),
		width:        width,
		minWidth:     32,
		height:       height,
		minHeight:    16,
		board:        cellbuffer{},
		interactive:  true,
		currentLevel: nil,
		nextLevel:    nil,
	}

	m.levelsList.Title = "pipes"

	m.levelsList.SetShowStatusBar(false)
	m.levelsList.SetShowHelp(false)

	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseAllMotion()}
}

func (m model) Init() tea.Cmd {
	return nil
}

type changeScreen struct {
	newScreen Screen
}

func delayedChangeScreen(delay time.Duration, screen Screen) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(delay * time.Second)
		return changeScreen{newScreen: screen}
	}
}

func loadLevel(m *model, level puzzles.Level) {
	m.currentLevel = &level
	idx := slices.IndexFunc(
		puzzles.Levels,
		func(i list.Item) bool {
			return i.(puzzles.Level).Title() == (&level).Title()
		},
	) + 1

	if idx < len(puzzles.Levels) {
		nextLevel := puzzles.Levels[idx].(puzzles.Level)
		m.nextLevel = &nextLevel
	} else {
		m.nextLevel = nil
	}

	m.board.init(*level.Board())
	recalculatePipes(&m.board)
	m.display = GameScreen
	m.interactive = true
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

	case changeScreen:
		m.display = msg.newScreen
		m.interactive = true

	}

	var cmd tea.Cmd

	if !m.interactive {
		return m, cmd
	}

	switch m.display {

	case LevelSelectScreen:
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				level, ok := m.levelsList.SelectedItem().(puzzles.Level)

				if ok {
					m.nextLevel = &level
					cmd = delayedChangeScreen(0, GameScreen)
				}

				return m, cmd
			}

		}
		m.levelsList, cmd = m.levelsList.Update(msg)
		return m, cmd

	case GameScreen:
		switch msg := msg.(type) {

		case changeScreen:
			if m.nextLevel == nil {
				cmd = delayedChangeScreen(0, LevelSelectScreen)
			} else {
				loadLevel(&m, *m.nextLevel)
			}

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

					if consumersAllHappy(&m.board) {
						m.interactive = false
						cmd = delayedChangeScreen(1, CompleteScreen)
					}
				}
			}

		}
		return m, cmd

	case CompleteScreen:
		switch msg.(type) {

		case changeScreen:
			cmd = delayedChangeScreen(3, GameScreen)

		}
		return m, cmd

	}

	return m, cmd
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
		s := ""
		s += "\nlevel complete!\n"
		s += "\n"
		if m.nextLevel != nil {
			s += "moving onto next level:\n"
			s += m.nextLevel.Title()
		} else {
			s += "that was the last level\n"
			s += "thanks for playing!\n"
		}
		return completeStyle.Render(s)
	}
	return ""
}
