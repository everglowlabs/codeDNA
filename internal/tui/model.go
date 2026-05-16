package tui

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/everglowlabs/codedna/internal/cluster"
	"github.com/everglowlabs/codedna/internal/optimizer"
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
		m.step = done
		m.logs = append(m.logs, fmt.Sprintf("Clustered into %d patterns", msg.numClusters))
		m.logs = append(m.logs, fmt.Sprintf("Selected %d golden samples", msg.numSamples))
		return m, m.progress.SetPercent(1.0)
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
		var allBlocks []parser.CodeBlock
		for _, f := range files {
			// Notify UI of file being parsed for progress
			// Note: This requires a bit of a trick in standard tea.Cmd
			// For simplicity, we'll just parse and then cluster
			blocks, err := p.ExtractBlocks(f)
			if err == nil {
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

		return clusteringFinishedMsg{
			numClusters: len(clusters),
			numSamples:  len(samples),
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
	numClusters int
	numSamples  int
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
