package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
)

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("CodeDNA // Standard Extraction Engine"))
	s.WriteString("\n\n")

	// Progress bar
	s.WriteString(fmt.Sprintf("Scanning: %s %d%%\n", m.progress.View(), int(m.progress.Percent()*100)))
	s.WriteString(fmt.Sprintf("Currently processing: %s\n\n", infoStyle.Render(m.currentFile)))

	// Analytics section
	s.WriteString(lipgloss.NewStyle().Bold(true).Render("ANALYTICS") + "\n")
	for _, p := range m.patternsFound {
		s.WriteString(fmt.Sprintf(" ├─ [%s] %s\n", successStyle.Render("✔"), p))
	}
	if len(m.patternsFound) == 0 {
		s.WriteString(" ├─ " + m.spinner.View() + " Searching for patterns...\n")
	}

	s.WriteString("\n")

	// Logs section
	s.WriteString(lipgloss.NewStyle().Bold(true).Render("LOGS") + "\n")
	for _, log := range m.logs {
		s.WriteString(fmt.Sprintf(" │ %s\n", log))
	}

	return borderStyle.Render(s.String())
}
