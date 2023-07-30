package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"os"

	"github.com/themichaellai/blame-tui/git"
)

type model struct {
	filename  string
	lines     []git.BlameLine
	cursorPos int

	err error
}

func initialModel(filename string) model {
	var lines []git.BlameLine
	return model{
		filename:  filename,
		lines:     lines,
		cursorPos: 0,
	}
}

func (m model) Init() tea.Cmd {
	return getBlameLines(m.filename)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursorPos > 0 {
				m.cursorPos--
			}

		case "down", "j":
			if m.cursorPos < len(m.lines)-1 {
				m.cursorPos++
			}

		case "enter", " ":
			//_, ok := m.selected[m.cursorPos]
			//if ok {
			//	delete(m.selected, m.cursorPos)
			//} else {
			//	m.selected[m.cursorPos] = struct{}{}
			//}
		}

	case blameLinesMsg:
		m.lines = msg.lines

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m model) View() string {
	res := ""
	for _, line := range m.lines {
		res += fmt.Sprintf("%s %s\n", line.AuthorName, line.Code)
	}
	return res
}

func getBlameLines(filename string) tea.Cmd {
	return func() tea.Msg {
		lines, err := git.Blame(filename)
		if err != nil {
			return errMsg{err}
		}
		return blameLinesMsg{lines}
	}
}

type blameLinesMsg struct {
	lines []git.BlameLine
}
type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }

func run() error {
	p := tea.NewProgram(initialModel(os.Args[1]))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
