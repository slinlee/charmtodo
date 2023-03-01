package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"

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
	Date  time.Time
	Value float64
}

var calData []CalDataPoint

func addCalData(date time.Time, val float64) {
	// Create new cal data point and add to cal data
	newPoint := CalDataPoint{date, val}
	calData = append(calData, newPoint)
}

func getIndexDate(x int, y int) time.Time {
	// compare the x,y to today and subtract
	today := time.Now()
	todayX, todayY := getDateIndex(today)

	diffX := todayX - x
	diffY := todayY - y

	diffDays := diffX*7 + diffY

	targetDate := today.AddDate(0, 0, -diffDays)
	return targetDate
}

func saveToFile() {

	// ** To save a file
	file, err := json.MarshalIndent(calData, "", " ")
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	_ = ioutil.WriteFile("s0br.json", file, 0644)
}

func readFromFile() {

	// Get Data from File
	content, err := ioutil.ReadFile("./s0br.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &calData)
	if err != nil {
		log.Fatal("Error during Unmarshall(): ", err)
	}

}

func getDateIndex(date time.Time) (int, int) {

	// calculate index
	today := time.Now()

	// How many weeks ago is this day
	// - compared to the 'Sunday' of this week
	difference := int(math.Ceil(
		(today.
			AddDate(0, 0, -int(today.Weekday())). // get the 'Sunday' of this week
			Sub(date.Local()).Hours() / 24 / 7))) // calculate number of weeks ago

	x := 52 - difference - 1

	if date.Local().Weekday() == time.Sunday {
		x++
	}

	dayOfWeek := int(date.Local().Weekday())

	return x, dayOfWeek
}

func parseCalToView(calData []CalDataPoint) {
	for _, v := range calData {
		x, y := getDateIndex(v.Date)
		viewData[x][y].actual += v.Value
	}
	normalizeViewData()
}

func normalizeViewData() {
	var min float64
	var max float64

	// Find min/max
	min = viewData[0][0].actual
	max = viewData[0][0].actual

	for _, row := range viewData {
		for _, val := range row {

			if val.actual < min {
				min = val.actual
			}
			if val.actual > max {
				max = val.actual
			}
		}

	}

	// Normalize the data
	for i, row := range viewData {
		for j, val := range row {
			viewData[i][j].normalized = (val.actual - min) / (max - min)
		}
	}
}

var viewData [52][7]viewDataPoint

type viewDataPoint struct {
	actual     float64
	normalized float64
}

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
		case "ctrl+c":
			return m, tea.Quit
		case "q": // only save when using `q` to quit
			saveToFile()
			return m, tea.Quit
		case "up", "k":
			if m.selectedY > 0 {
				m.selectedY--
			} else if m.selectedY == 0 && m.selectedX > 0 {
				m.selectedY = 6
				m.selectedX--
			}

		case "down", "j":
			// Don't allow user to scroll beyond today
			if m.selectedY < 6 &&
				(m.selectedX != 51 ||
					m.selectedY < int(time.Now().Weekday())) {
				m.selectedY++
			} else if m.selectedY == 6 && m.selectedX != 51 {
				m.selectedY = 0
				m.selectedX++
			}
		case "right", "l":
			// Don't allow users to scroll beyond today from the previous column
			if m.selectedX < 50 ||
				(m.selectedX == 50 &&
					m.selectedY <= int(time.Now().Weekday())) {
				m.selectedX++
			}
		case "left", "h":
			if m.selectedX > 0 {
				m.selectedX--
			}
		case "enter", " ":
			// Hard coded to add a new entry with value `1.0`
			addCalData(
				getIndexDate(m.selectedX, m.selectedY),

				1.0)
			parseCalToView(calData)

		}
	}
	return m, nil
}

func (m model) View() string {
	// The header

	theTime := getIndexDate(m.selectedX, m.selectedY) //time.Now()

	title, _ := glamour.Render(theTime.Format("# Monday, January 02, 2006"), "dark")
	s := title

	selectedDetail := "    Value: " + fmt.Sprint(viewData[m.selectedX][m.selectedY].actual) + " normalized: " + fmt.Sprint(viewData[m.selectedX][m.selectedY].normalized) + "\n\n"

	s += selectedDetail

	var labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	var boxStyle = lipgloss.NewStyle().
		PaddingRight(1).
		Foreground(lipgloss.Color(scaleColors[2]))

	var boxSelectedStyle = boxStyle.Copy().
		Background(lipgloss.Color("#9999ff")).
		Foreground(lipgloss.Color(scaleColors[0]))

	// Month Labels
	var currMonth time.Month
	s += "  "
	for j := 0; j < 52; j++ {
		// Check the last day of the week for that column
		jMonth := getIndexDate(j, 6).Month()

		if currMonth != jMonth {
			currMonth = jMonth
			s += labelStyle.Render(getIndexDate(j, 6).Format("Jan") + " ")

			// Skip the length of the label we just added
			j += 1
		} else {
			s += "  "
		}
	}
	s += "\n"

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
			// Selected Item
			if m.selectedX == i && m.selectedY == j {
				s += boxSelectedStyle.Copy().Foreground(
					lipgloss.Color(
						getScaleColor(
							viewData[i][j].normalized))).
					Render("■")
			} else if i == 51 &&
				j > int(time.Now().Weekday()) {
				// In the future
				s += boxStyle.Render(" ")
			} else {
				// Not Selected Item and not in the future
				s += boxStyle.Copy().
					Foreground(
						lipgloss.Color(
							getScaleColor(
								viewData[i][j].normalized))).
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
	readFromFile()
	// Parse Data
	parseCalToView(calData)

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}