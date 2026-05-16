# Pattern Detected in internal/tui/model.go

Category: General
Rationale: Automatically extracted from codebase structure.

## Guidelines
- Follow the structural pattern demonstrated in the sample.

## Examples
```go
// File: internal/tui/model.go
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
		m.Result = msg.dna
		m.logs = append(m.logs, fmt.Sprintf("Clustered into %d patterns", len(msg.dna.Standards)))
		m.logs = append(m.logs, "Selected golden samples for all patterns")
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
```

