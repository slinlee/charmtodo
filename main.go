package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	heatmap "github.com/slinlee/bubbletea-heatmap"

	"os"
	"time"
)

type model struct {
	heatmap heatmap.Model
	calData []heatmap.CalDataPoint
}

func addCalData(date time.Time, val float64) {
	// Create new cal data point and add to cal data
	newPoint := heatmap.CalDataPoint{date, val}
	calData = append(calData, newPoint)
}

func saveToFile(filename string) {
	// return // debug
	// ** To save a file
	file, err := json.MarshalIndent(calData, "", " ")
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	_ = ioutil.WriteFile(filename, file, 0644)
}

func readFromFile(filename string) {

	// Get Data from File
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &calData)
	if err != nil {
		log.Fatal("Error during Unmarshall(): ", err)
	}

}

func readMockData() {
	// Generate mock data for debugging

	today := time.Now()

	for i := 0; i < 350; i++ {
		addCalData(today.AddDate(0, 0, -i), float64(i%2))
	}

}

func initialModel() model {
	hm := heatmap.New()
	return model{
		calData: []heatmap.CalDataPoint,
		heatmap: hm,
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
			saveToFile("./s0br.json")
			return m, tea.Quit
		case "enter", " ":
			// Hard coded to add a new entry with value `1.0`
			// addCalData(
			// 	getIndexDate(m.selectedX, m.selectedY),

			// 	1.0)
			// parseCalToView(calData)

		}

	}
	m.heatmap, cmd = m.heatmap.Update(msg)

	return m, cmd
}

func (m model) View() string {
	// The header

	theTime := getIndexDate(m.selectedX, m.selectedY) //time.Now()

	title, _ := glamour.Render(theTime.Format("# Monday, January 02, 2006"), "dark")
	s := title

	selectedDetail := "    Value: " + fmt.Sprint(viewData[m.selectedX][m.selectedY].actual) + " normalized: " + fmt.Sprint(viewData[m.selectedX][m.selectedY].normalized) + "\n\n"

	s += selectedDetail

	s += m.heatmap.View()

	// The footer
	s += "\nPress q to quit.\n"


	return s
}

func main() {
	readFromFile("./s0br.json")
	// readMockData() // debug


	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}