package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/everglowlabs/codedna/internal/parser"
	"github.com/everglowlabs/codedna/internal/scanner"
)

type scanStep int

const (
	scanningFiles scanStep = iota
	parsingAST
	clustering
	generatingRules
	done
)

type model struct {
	path            string
	step            scanStep
	progress        progress.Model
	spinner         spinner.Model
	currentFile     string
	filesScanned    int
	totalFiles      int
	patternsFound   []string
	logs            []string
}

func InitialModel(path string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{
		path:     path,
		step:     scanningFiles,
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.startScan(),
	)
}

func (m model) startScan() tea.Cmd {
	return func() tea.Msg {
		s := scanner.NewScanner(m.path)
		count := 0
		files, err := s.Scan(func(path string, total int) {
			count++
			// Note: In a production app, we'd use a channel to stream updates
			// But for this skeleton, we'll just process after scanning
		})
		if err != nil {
			return logMsg{text: fmt.Sprintf("Error scanning: %v", err)}
		}
		return scanFinishedMsg{files: files}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if p, ok := newModel.(progress.Model); ok {
			m.progress = p
		}
		return m, cmd
	case scanFinishedMsg:
		m.step = parsingAST
		m.totalFiles = len(msg.files)
		m.logs = append(m.logs, fmt.Sprintf("Found %d files to analyze", m.totalFiles))
		return m, m.parseFiles(msg.files)
	case patternFoundMsg:
		m.patternsFound = append(m.patternsFound, msg.pattern)
		return m, nil
	case logMsg:
		m.logs = append(m.logs, msg.text)
		if len(m.logs) > 5 {
			m.logs = m.logs[1:]
		}
		return m, nil
	}
	return m, nil
}

func (m model) parseFiles(files []string) tea.Cmd {
	return func() tea.Msg {
		p := parser.NewGoParser()
		for _, f := range files {
			funcs, err := p.ExtractFunctions(f)
			if err == nil && len(funcs) > 0 {
				return patternFoundMsg{pattern: fmt.Sprintf("Functions in %s: %d", f, len(funcs))}
			}
		}
		return logMsg{text: "AST Analysis complete"}
	}
}

type scanFinishedMsg struct {
	files []string
}

type patternFoundMsg struct {
	pattern string
}

type logMsg struct {
	text string
}
