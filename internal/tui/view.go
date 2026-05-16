package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#5C5C5C"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("#874BFD")).
			MarginBottom(1)

	columnStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3C3C3C"))

	focusedStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD"))

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3C3C3C")).
			Padding(0, 1).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF"))
	dnaStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
)

func (m model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	// 1. Header
	header := titleStyle.Render("CodeDNA 🧬") + " " + subtleStyle.Render("v0.1.0") + " | " + infoStyle.Render(m.path)
	header = headerStyle.Width(m.width - 2).Render(header)

	// 2. Body (Two Columns)
	sidebarWidth := m.width / 4
	if sidebarWidth < 25 {
		sidebarWidth = 25
	}
	mainWidth := m.width - sidebarWidth - 6

	// Sidebar: Engineering DNA Summary
	var sidebarBuilder strings.Builder
	sidebarBuilder.WriteString(dnaStyle.Render("ENGINEERING DNA") + "\n\n")

	if len(m.patternsFound) > 0 {
		for _, p := range m.patternsFound {
			fmt.Fprintf(&sidebarBuilder, " %s %s\n", successStyle.Render("✔"), p)
		}
	} else if m.step == done {
		sidebarBuilder.WriteString(subtleStyle.Render(" No patterns detected.\n"))
	} else {
		fmt.Fprintf(&sidebarBuilder, " %s Analyzing code...\n", m.spinner.View())
	}

	if m.step == done {
		sidebarBuilder.WriteString(
			"\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00")).Bold(true).Render("SUMMARY") + "\n")
		fmt.Fprintf(&sidebarBuilder, " Files: %d\n", m.totalFiles)
		fmt.Fprintf(&sidebarBuilder, " Patterns: %d\n", len(m.patternsFound))
	}

	sidebar := focusedStyle.Width(sidebarWidth).Height(m.height - 10).Render(sidebarBuilder.String())

	// Main: Extraction Logs
	var mainBuilder strings.Builder
	mainBuilder.WriteString(lipgloss.NewStyle().Bold(true).Render("EXTRACTION LOGS") + "\n\n")
	for _, log := range m.logs {
		icon := infoStyle.Render("→")
		if strings.Contains(log, "Clustered") || strings.Contains(log, "Selected") {
			icon = successStyle.Render("✔")
		}
		fmt.Fprintf(&mainBuilder, " %s %s\n", icon, log)
	}
	if m.step != done {
		fmt.Fprintf(&mainBuilder, "\n %s Current: %s", m.spinner.View(), infoStyle.Render(m.currentFile))
	} else {
		mainBuilder.WriteString("\n " + successStyle.Render("⠿") + " All rules synthesized in internal schema.")
	}
	mainArea := columnStyle.Width(mainWidth).Height(m.height - 10).Render(mainBuilder.String())

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainArea)

	// 3. Status Bar & Footer
	status := "READY"
	statusColor := lipgloss.Color("#3C3C3C")
	switch m.step {
	case scanningFiles:
		status = "SCANNING"
		statusColor = lipgloss.Color("#00BFFF")
	case parsingAST:
		status = "PARSING"
		statusColor = lipgloss.Color("#7D56F4")
	case clustering:
		status = "CLUSTERING"
		statusColor = lipgloss.Color("#FF00FF")
	case done:
		status = "COMPLETE"
		statusColor = lipgloss.Color("#00FF00")
	}

	statusView := footerStyle.Background(statusColor).Render(status)

	// Progress bar with fixed width
	m.progress.Width = m.width / 2
	progressBar := m.progress.View()

	footer := lipgloss.JoinHorizontal(lipgloss.Center,
		statusView,
		"  ",
		progressBar,
		"  ",
		subtleStyle.Render("press q to exit"),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		footer,
	)
}
