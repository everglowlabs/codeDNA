package tui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/everglowlabs/codedna/internal/cluster"
	"github.com/everglowlabs/codedna/internal/optimizer"
	"github.com/everglowlabs/codedna/internal/parser"
	"github.com/everglowlabs/codedna/internal/scanner"
	"github.com/everglowlabs/codedna/internal/schema"
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
	path          string
	step          scanStep
	progress      progress.Model
	spinner       spinner.Model
	currentFile    string
	totalFiles     int
	processedFiles int
	patternsFound  []string
	logs          []string
	width         int
	height        int
	Result        schema.DNA_Schema // Final result

	// Phase 4 additions
	noInteractive bool
	provider      string
	llmModel      string
	scrubStrings  bool

	// Interactive editor state
	reviewing         bool
	reviewIdx         int
	editingRule       bool
	textarea          textarea.Model
	tempDNA           schema.DNA_Schema
	acceptedStandards []schema.Standard
}

func InitialModel(path string, noInteractive bool, provider string, llmModel string, scrubStrings bool) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	
	ta := textarea.New()
	ta.Placeholder = "Enter guidelines (one per line)..."
	ta.SetWidth(60)
	ta.SetHeight(6)
	ta.Focus()

	return model{
		path:          path,
		step:          scanningFiles,
		spinner:       s,
		progress:      progress.New(progress.WithDefaultGradient()),
		noInteractive: noInteractive,
		provider:      provider,
		llmModel:      llmModel,
		scrubStrings:  scrubStrings,
		textarea:      ta,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.startScan(),
		tea.EnterAltScreen,
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
	// 1. Textarea keys if in editing mode
	if m.editingRule {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+s":
				lines := strings.Split(m.textarea.Value(), "\n")
				var cleaned []string
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						cleaned = append(cleaned, line)
					}
				}
				m.tempDNA.Standards[m.reviewIdx].Rules = cleaned
				m.editingRule = false
				return m, nil
			case "esc":
				m.editingRule = false
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	// 2. Normal review keys if reviewing rules
	if m.reviewing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "a", "A":
				m.acceptedStandards = append(m.acceptedStandards, m.tempDNA.Standards[m.reviewIdx])
				m.patternsFound = append(m.patternsFound, m.tempDNA.Standards[m.reviewIdx].Title)
				m.reviewIdx++
				if m.reviewIdx >= len(m.tempDNA.Standards) {
					m.reviewing = false
					m.step = done
					m.Result = m.tempDNA
					m.Result.Standards = m.acceptedStandards
					m.logs = append(m.logs, fmt.Sprintf("Completed review. Accepted %d standards.", len(m.Result.Standards)))
					return m, m.progress.SetPercent(1.0)
				}
				m.textarea.SetValue(strings.Join(m.tempDNA.Standards[m.reviewIdx].Rules, "\n"))
				return m, nil
			case "r", "R":
				m.logs = append(m.logs, fmt.Sprintf("Rejected standard: %s", m.tempDNA.Standards[m.reviewIdx].Title))
				m.reviewIdx++
				if m.reviewIdx >= len(m.tempDNA.Standards) {
					m.reviewing = false
					m.step = done
					m.Result = m.tempDNA
					m.Result.Standards = m.acceptedStandards
					m.logs = append(m.logs, fmt.Sprintf("Completed review. Accepted %d standards.", len(m.Result.Standards)))
					return m, m.progress.SetPercent(1.0)
				}
				m.textarea.SetValue(strings.Join(m.tempDNA.Standards[m.reviewIdx].Rules, "\n"))
				return m, nil
			case "e", "E":
				m.editingRule = true
				m.textarea.Focus()
				m.textarea.SetValue(strings.Join(m.tempDNA.Standards[m.reviewIdx].Rules, "\n"))
				return m, nil
			case "s", "S":
				for i := m.reviewIdx; i < len(m.tempDNA.Standards); i++ {
					m.acceptedStandards = append(m.acceptedStandards, m.tempDNA.Standards[i])
					m.patternsFound = append(m.patternsFound, m.tempDNA.Standards[i].Title)
				}
				m.reviewing = false
				m.step = done
				m.Result = m.tempDNA
				m.Result.Standards = m.acceptedStandards
				m.logs = append(m.logs, fmt.Sprintf("Skipped remaining reviews. Accepted %d standards total.", len(m.Result.Standards)))
				return m, m.progress.SetPercent(1.0)
			}
		}
	}

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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case scanFinishedMsg:
		m.step = parsingAST
		m.totalFiles = len(msg.files)
		m.logs = append(m.logs, fmt.Sprintf("Found %d files to analyze", m.totalFiles))
		return m, tea.Batch(
			m.progress.SetPercent(0.2), // 20% done after scanning
			m.parseFiles(msg.files),
		)
	case fileParsedMsg:
		m.processedFiles++
		m.currentFile = msg.file
		percent := 0.2 + (float64(m.processedFiles)/float64(m.totalFiles))*0.6 // 20% to 80% for parsing
		return m, m.progress.SetPercent(percent)
	case patternFoundMsg:
		m.patternsFound = append(m.patternsFound, msg.pattern)
		return m, nil
	case clusteringFinishedMsg:
		m.tempDNA = msg.dna
		if m.noInteractive || len(m.tempDNA.Standards) == 0 {
			m.step = done
			m.Result = m.tempDNA
			m.logs = append(m.logs, fmt.Sprintf("Clustered into %d patterns", len(m.Result.Standards)))
			m.logs = append(m.logs, "Selected golden samples for all patterns")
			return m, m.progress.SetPercent(1.0)
		} else {
			m.reviewing = true
			m.reviewIdx = 0
			m.acceptedStandards = []schema.Standard{}
			m.textarea.SetValue(strings.Join(m.tempDNA.Standards[0].Rules, "\n"))
			m.logs = append(m.logs, fmt.Sprintf("Extracted %d standards. Starting interactive TUI review...", len(m.tempDNA.Standards)))
			return m, nil
		}
	case logMsg:
		m.logs = append(m.logs, msg.text)
		if len(m.logs) > 5 {
			m.logs = m.logs[1:]
		}
		return m, nil
	}
	return m, nil
}

func (m model) GetResult() schema.DNA_Schema {
	return m.Result
}

func (m model) parseFiles(files []string) tea.Cmd {
	return func() tea.Msg {
		var allBlocks []parser.CodeBlock
		for _, f := range files {
			ext := filepath.Ext(f)
			p := parser.NewParser(ext)
			if p == nil {
				continue
			}

			blocks, err := p.ExtractBlocks(f)
			if err == nil {
				if m.scrubStrings {
					for idx := range blocks {
						blocks[idx].Content = parser.ScrubStrings(blocks[idx].Content)
					}
				}
				allBlocks = append(allBlocks, blocks...)
			}
		}

		// Now clustering
		cm, err := cluster.NewClusterManager("dna", dummyEmbedder)
		if err != nil {
			return logMsg{text: fmt.Sprintf("Cluster error: %v", err)}
		}

		ctx := context.Background()
		err = cm.AddBlocks(ctx, allBlocks)
		if err != nil {
			return logMsg{text: fmt.Sprintf("Add blocks error: %v", err)}
		}

		clusters, err := cm.GroupByPattern(ctx, allBlocks)
		if err != nil {
			return logMsg{text: fmt.Sprintf("Grouping error: %v", err)}
		}

		// Select golden samples
		opt := optimizer.NewOptimizer(4096)
		samples := opt.SelectGoldenSamples(clusters)
		dna := opt.GenerateDNA("Project", samples)

		return clusteringFinishedMsg{
			dna: dna,
		}
	}
}

func dummyEmbedder(ctx context.Context, text string) ([]float32, error) {
	res := make([]float32, 384)
	// Simple normalization: first element is 1, others 0
	res[0] = 1.0
	return res, nil
}

type clusteringFinishedMsg struct {
	dna schema.DNA_Schema
}

type scanFinishedMsg struct {
	files []string
}

type patternFoundMsg struct {
	pattern string
}

type fileParsedMsg struct {
	file string
}

type logMsg struct {
	text string
}
