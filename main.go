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

var scaleColors = []string{
	// Light Theme
	// #ebedf0 - Less
	// #9be9a8
	// #40c463
	// #30a14e
	// #216e39 - More

	// Dark Theme
	"#161b22", // Less
	"#0e4429",
	"#006d32",
	"#26a641",
	"#39d353", // - More

}

type CalDataPoint struct {
	date  time.Time
	value float64
}

var calDataMock = []CalDataPoint{

	{date: time.Now(), value: 1.0},
	{date: time.Now().AddDate(0, 0, -1), value: 0.30},
	{date: time.Now().AddDate(0, 0, -5), value: 0.10},
	{date: time.Now().AddDate(0, 0, -5), value: 0.20},
	{date: time.Now().AddDate(0, 0, -15), value: 0.60},
	{date: time.Now().AddDate(0, 0, -50), value: 0.60},
	// {date: time.Date(2023, 2, 11, 0, 0, 0, 0), value: 1.0},
	// {date: time.create("2023, 2, 12"), value: 1.0},
	// {date: time.create("2023, 2, 13"), value: 1.0},
	// {date: time.create("2023, 2, 14"), value: 1.0},
	// {date: time.create("2023, 2, 15"), value: 1.0},
	// {date: time.create("2023, 2, 16"), value: 1.0},
	// {date: time.create("2023, 2, 17"), value: 1.0},
	// {date: time.create("2023, 2, 18"), value: 1.0},
	// {date: time.create("2023, 2, 19"), value: 1.0},
	// {date: time.create("2023, 2, 20"), value: 1.0},
	// {date: time.create("2023, 2, 21"), value: 1.0},
	// {date: time.create("2023, 2, 22"), value: 1.0},
	// {date: time.create("2023, 2, 23"), value: 1.0},
	// {date: time.create("2023, 2, 24"), value: 1.0},
}

func getDateIndex(date time.Time) (int, int) {

	// calculate index
	today := time.Now()

	// How many weeks ago is this day
	difference := int(today.Sub(date).Hours() / 24 / 7)

	dayOfWeek := int(date.Weekday())

	x := 52 - difference - 1

	return x, dayOfWeek
}

func parseCalToView(calData []CalDataPoint) {
	for i, v := range calData {
		x, y := getDateIndex(v.date)
		fmt.Println(i, v.date.Weekday(), int(v.date.Weekday()), x, y)
		viewData[x][y] += v.value
	}
}

var viewData [52][7]float64

func getScaleColor(value float64) string {
	const numColors = 5
	const max = 1.0
	const min = 0.0

	return scaleColors[int((value/max)*(numColors-1))]
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

	title, _ := glamour.Render(theTime.Format("# Monday, 2006-1-2"), "dark")
	s := title

	var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	var boxStyle = lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color(scaleColors[2]))

	var boxSelectedStyle = boxStyle.Copy().
		Background(lipgloss.Color("#9999ff")).
		Foreground(lipgloss.Color(scaleColors[0]))

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

		for i := 0; i < 52; i++ {
			if m.selectedX == i && m.selectedY == j {
				s += boxSelectedStyle.Copy().Foreground(
					lipgloss.Color(
						getScaleColor(
							viewData[i][j]))).
					Render("■")
			} else {
				s += boxStyle.Copy().
					Foreground(
						lipgloss.Color(
							getScaleColor(
								viewData[i][j]))).
					Render("■")
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
	parseCalToView(calDataMock)
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}