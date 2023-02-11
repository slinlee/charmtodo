package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"os"
	"time"
)

type model struct {
	selectedX int
	selectedY int
}

func initialModel() model {
	return model{
		selectedX: 0,
		selectedY: 0,
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
			if m.selectedY > 0 {
				m.selectedY--
			}

		case "down", "j":
			if m.selectedY < 6 {
				m.selectedY++
			}
		case "right", "l":
			if m.selectedX < 51 {
				m.selectedX++
			}
		case "left", "h":
			if m.selectedX > 0 {
				m.selectedX--
			}
		case "enter", " ":
		}
	}
	return m, nil
}

func (m model) View() string {
	// The header

	theTime := time.Now()

	title, _ := glamour.Render(theTime.Format("# 2006-1-2"), "dark")
	s := title

	var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	for j := 0; j < 7; j++ {
		switch j {
		case 0:
			s += labelStyle.Render("S ")
		case 1:
			s += labelStyle.Render("M ")
		case 2:
			s += labelStyle.Render("T ")
		case 3:
			s += labelStyle.Render("W ")
		case 4:
			s += labelStyle.Render("T ")
		case 5:
			s += labelStyle.Render("F ")
		case 6:
			s += labelStyle.Render("S ")
		}

		var boxStyle = lipgloss.NewStyle().
			PaddingRight(1).
			Foreground(lipgloss.Color("#04B575"))

		var boxSelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#9999ff")).
			PaddingRight(1).
			Foreground(lipgloss.Color("#04B575"))

		for i := 0; i < 52; i++ {
			if m.selectedX == i && m.selectedY == j {
				s += boxSelectedStyle.Render("■")
			} else {
				s += boxStyle.Render("■")
			}
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