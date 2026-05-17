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

	if m.reviewing {
		var bodyBuilder strings.Builder
		std := m.tempDNA.Standards[m.reviewIdx]

		if m.editingRule {
			bodyBuilder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true).Render("EDITING RULE DETAILS") + "\n")
			bodyBuilder.WriteString(subtleStyle.Render("Modify rules below. Use Ctrl+S to save, Esc to cancel.") + "\n\n")
			fmt.Fprintf(&bodyBuilder, "Title:    %s\n", std.Title)
			fmt.Fprintf(&bodyBuilder, "Category: %s\n\n", std.Category)
			bodyBuilder.WriteString(m.textarea.View() + "\n\n")
			bodyBuilder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render("[Ctrl+S] Save changes    [Esc] Cancel") + "\n")
		} else {
			bodyBuilder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true).Render("NEW RULE DETECTED FOR REVIEW") + " (" + infoStyle.Render(fmt.Sprintf("%d/%d", m.reviewIdx+1, len(m.tempDNA.Standards))) + ")\n")
			bodyBuilder.WriteString(subtleStyle.Render("─────────────────────────────────────────────────────────────") + "\n\n")
			
			fmt.Fprintf(&bodyBuilder, "Title:    %s\n", lipgloss.NewStyle().Bold(true).Render(std.Title))
			fmt.Fprintf(&bodyBuilder, "Category: %s  |  Severity: %s\n", std.Category, std.Severity)
			fmt.Fprintf(&bodyBuilder, "Rationale: %s\n\n", std.Rationale)
			
			bodyBuilder.WriteString(lipgloss.NewStyle().Bold(true).Render("Guidelines:") + "\n")
			for _, r := range std.Rules {
				fmt.Fprintf(&bodyBuilder, "  • %s\n", r)
			}
			bodyBuilder.WriteString("\n")
			
			if len(std.Samples) > 0 {
				bodyBuilder.WriteString(lipgloss.NewStyle().Bold(true).Render("Sample Location:") + "\n")
				fmt.Fprintf(&bodyBuilder, "  %s\n\n", subtleStyle.Render(std.Samples[0].FilePath))
			}
			
			bodyBuilder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("[A] Accept") + "   " +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00")).Render("[E] Edit") + "   " +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("[R] Reject") + "   " +
				lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")).Render("[S] Skip Remaining") + "\n")
		}

		bodyBox := focusedStyle.Width(m.width - 6).Height(m.height - 10).Render(bodyBuilder.String())

		// Progress bar with fixed width
		m.progress.Width = m.width / 2
		percent := float64(m.reviewIdx) / float64(len(m.tempDNA.Standards))
		progressBar := m.progress.ViewAs(percent)

		statusView := footerStyle.Background(lipgloss.Color("#FF00FF")).Render("REVIEWING")

		footer := lipgloss.JoinHorizontal(lipgloss.Center,
			statusView,
			"  ",
			progressBar,
			"  ",
			subtleStyle.Render("press q to exit"),
		)

		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			bodyBox,
			footer,
		)
	}

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
