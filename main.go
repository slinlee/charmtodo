package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	// "github.com/charmbracelet/lipgloss"

	"os"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices: []string{"buy a", "buy b", "buy c"},

		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	// The header

	// s := "@slee 2023\n\n"
	s := ""
	title, _ := glamour.Render("# @slee 2023", "dark")
	fmt.Print(title)

	// Iterate over our choices

	// for i, choice := range m.choices {
	for j := 0; j < 7; j++ {
		switch j {
		case 0:

			s += "S "
		case 1:

			s += "M "
		case 2:

			s += "T "
		case 3:

			s += "W "
		case 4:

			s += "T "
		case 5:

			s += "F "
		case 6:

			s += "S "
		}

		for i := 0; i < 52; i++ {

			s += "▓"
		}
		s += "\n"
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UY for rendering

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}