package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	heatmap "github.com/slinlee/bubbletea-heatmap"
)

type model struct {
	heatmap heatmap.Model
	calData []heatmap.CalDataPoint
}

func (m model) addCalData(date time.Time, val float64) {
	// Create new cal data point and add to cal data
	newPoint := heatmap.CalDataPoint{Date: date, Value: val}
	m.calData = append(m.calData, newPoint)
}

func (m model) saveToFile(filename string) {
	// return // debug
	// ** To save a file
	file, err := json.MarshalIndent(m.calData, "", " ")
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	_ = os.WriteFile(filename, file, 0o644)
}

func readFromFile(filename string) []heatmap.CalDataPoint {
	var fileData []heatmap.CalDataPoint

	// Get Data from File
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &fileData)
	if err != nil {
		log.Fatal("Error during Unmarshall(): ", err)
	}

	return fileData
}

func initialModel() model {
	fileData := readFromFile("./s0br.json")
	hm := heatmap.New(fileData)
	return model{
		calData: fileData,
		heatmap: hm,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q": // only save when using `q` to quit
			m.saveToFile("./s0br.json")
			return m, tea.Quit
		case "enter", " ":
			// TODO - Decide where to handle this
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

	// theTime := getIndexDate(m.selectedX, m.selectedY) //time.Now()

	// title, _ := glamour.Render(theTime.Format("# Monday, January 02, 2006"), "dark")
	// s := title

	// selectedDetail := "    Value: " + fmt.Sprint(viewData[m.selectedX][m.selectedY].actual) + " normalized: " + fmt.Sprint(viewData[m.selectedX][m.selectedY].normalized) + "\n\n"

	// s += selectedDetail

	s := m.heatmap.View()
	// The footer
	s += "\nPress q to quit.\n"

	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
