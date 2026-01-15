package tui

import (
	"archon/internal/adapters/gemini"
	"archon/internal/adapters/vectordb"
	"archon/internal/config"
	"archon/internal/core"
	"archon/internal/utils"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/charmbracelet/lipgloss"
)


func StartApp() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type state int

const (
	stateMenu state = iota
	stateChat
	stateIndex
	stateStatus
	stateContext
	stateInputPath
)

type model struct {
	state          state
	choices        []string
	cursor         int
	textInput      textinput.Model
	viewport       viewport.Model
	viewportContext viewport.Model
	table          table.Model
	spinner        spinner.Model
	progress       progress.Model
	thinking       bool
	indexing       bool
	chatHistory    string
	lastContext    string
	err            error
	indexStats     string
	statusInfo     string
	width          int
	height         int
	selectedAction string
	totalTokens    int
	lastPromptTokens int
	lastAnswerTokens int
	totalCost      float64
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(SecondaryColor)

	prog := progress.New(progress.WithDefaultGradient())

	vp := viewport.New(80, 15)
	vp.SetContent("Welcome to Archon Chat Mode! Ask anything about your codebase.")

	vpc := viewport.New(80, 15)
	vpc.SetContent("RAG Context will appear here when you ask a question.")

	columns := []table.Column{
		{Title: "Component", Width: 20},
		{Title: "Status/Value", Width: 50},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	return model{
		state:           stateMenu,
		choices:         []string{
			"Chat Mode", 
			"Index Codebase", 
			"AI Code Review",
			"Smart Commit",
			"Explain File/Symbol",
			"Refactor Code",
			"Generate Unit Tests",
			"Architectural Analysis",
			"Generate Diagram",
			"System Status", 
			"Clear Index",
			"Exit",
		},
		textInput:       ti,
		spinner:         s,
		progress:        prog,
		viewport:        vp,
		viewportContext: vpc,
		table:           tbl,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = max(0, msg.Width-4)
		m.viewport.Height = max(0, msg.Height-16)
		m.viewportContext.Width = max(0, msg.Width-4)
		m.viewportContext.Height = max(0, msg.Height-16)
		m.progress.Width = max(0, msg.Width-10)
		m.table.SetWidth(max(0, msg.Width-10))
		m.table.SetHeight(max(0, msg.Height-20))

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.state != stateMenu {
				m.state = stateMenu
				m.err = nil
				return m, nil
			}
		case "tab":
			if m.state == stateChat {
				m.state = stateContext
			} else if m.state == stateContext {
				m.state = stateChat
			}
		case "up", "k":
			if m.state == stateMenu && m.cursor > 0 {
				m.cursor--
			} else if m.state == stateStatus {
				m.table.MoveUp(1)
			}
		case "down", "j":
			if m.state == stateMenu && m.cursor < len(m.choices)-1 {
				m.cursor++
			} else if m.state == stateStatus {
				m.table.MoveDown(1)
			}
		case "enter":
			if m.state == stateMenu {
				choice := m.choices[m.cursor]
				switch choice {
				case "Chat Mode":
					m.state = stateChat
				case "Index Codebase":
					m.state = stateIndex
					m.indexing = true
					return m, m.startIndexing()
				case "AI Code Review":
					m.state = stateChat
					m.thinking = true
					m.chatHistory += fmt.Sprintf("\n%s You: Running AI Code Review on staged changes\n", UserMsgStyle.Render("●"))
					m.viewport.SetContent(m.chatHistory)
					m.viewport.GotoBottom()
					
					diff, err := utils.GetStagedDiff()
					if err != nil {
						m.err = err
						m.thinking = false
						return m, nil
					}
					if strings.TrimSpace(diff) == "" {
						m.chatHistory += "\nArchon: No staged changes (git add) to review.\n"
						m.viewport.SetContent(m.chatHistory)
						m.thinking = false
						return m, nil
					}
					
					prompt := fmt.Sprintf("Perform a code review on the following git diff:\n\n```diff\n%s\n```\n\nIdentify potential bugs, best practice violations, and provide improvement suggestions.", diff)
					return m, m.askGemini(prompt)
					
				case "Smart Commit":
					m.state = stateChat
					m.thinking = true
					m.chatHistory += fmt.Sprintf("\n%s You: Generating smart commit message\n", UserMsgStyle.Render("●"))
					m.viewport.SetContent(m.chatHistory)
					m.viewport.GotoBottom()
					
					diff, err := utils.GetStagedDiff()
					if err != nil {
						m.err = err
						m.thinking = false
						return m, nil
					}
					if strings.TrimSpace(diff) == "" {
						m.chatHistory += "\nArchon: No staged changes to generate a commit message for.\n"
						m.viewport.SetContent(m.chatHistory)
						m.thinking = false
						return m, nil
					}
					
					prompt := fmt.Sprintf("Analyze the following changes and generate a commit message following Conventional Commits standards:\n\n```diff\n%s\n```\n\nProvide 3 commit message options.", diff)
					return m, m.askGemini(prompt)

				case "Explain File/Symbol", "Refactor Code", "Generate Unit Tests", "Generate Diagram":
					m.state = stateInputPath
					m.selectedAction = choice
					m.textInput.SetValue("")
					if choice == "Generate Diagram" {
						m.textInput.Placeholder = "Focus area (e.g., AuthSystem)..."
					} else {
						m.textInput.Placeholder = "File path (e.g., ./main.go)..."
					}
					return m, nil
				case "Architectural Analysis":
					m.state = stateChat
					m.thinking = true
					m.chatHistory += fmt.Sprintf("\n%s You: Running Architectural Analysis\n", UserMsgStyle.Render("●"))
					m.viewport.SetContent(m.chatHistory)
					m.viewport.GotoBottom()
					return m, m.askGemini("Perform a deep architectural analysis on this project. Detect anomalies, code smells, or design pattern violations.")
				case "System Status":
					m.state = stateStatus
					return m, m.loadStatus()
				case "Clear Index":
					m.state = stateStatus
					return m, m.clearIndex()
				case "Exit":
					return m, tea.Quit
				}
			} else if m.state == stateInputPath {
				input := m.textInput.Value()
				if input != "" {
					m.state = stateChat
					m.thinking = true
					var prompt string
					switch m.selectedAction {
					case "Explain File/Symbol":
						prompt = fmt.Sprintf("Explain in detail about: %s", input)
						m.chatHistory += fmt.Sprintf("\n%s You: Explain %s\n", UserMsgStyle.Render("●"), input)
					case "Refactor Code":
						prompt = fmt.Sprintf("Provide refactoring suggestions for file: %s", input)
						m.chatHistory += fmt.Sprintf("\n%s You: Refactor %s\n", UserMsgStyle.Render("●"), input)
					case "Generate Unit Tests":
						prompt = fmt.Sprintf("Create comprehensive unit tests for file: %s", input)
						m.chatHistory += fmt.Sprintf("\n%s You: Create test for %s\n", UserMsgStyle.Render("●"), input)
					case "Generate Diagram":
						prompt = fmt.Sprintf("Create Mermaid diagram code for focus area: %s", input)
						m.chatHistory += fmt.Sprintf("\n%s You: Create diagram for %s\n", UserMsgStyle.Render("●"), input)
					}
					m.viewport.SetContent(m.chatHistory)
					m.viewport.GotoBottom()
					m.textInput.SetValue("")
					m.textInput.Placeholder = "Type your question here..."
					return m, m.askGemini(prompt)
				}
			} else if m.state == stateChat {
				query := m.textInput.Value()
				if query != "" {
					m.thinking = true
					m.chatHistory += fmt.Sprintf("\n%s You: %s\n", UserMsgStyle.Render("●"), query)
					m.viewport.SetContent(m.chatHistory)
					m.viewport.GotoBottom()
					m.textInput.SetValue("")
					return m, m.askGemini(query)
				}
			}
		}
	case indexCompleteMsg:
		m.indexing = false
		m.indexStats = string(msg)
		return m, nil
	case geminiResponseMsg:
		m.thinking = false
		m.chatHistory += fmt.Sprintf("\n%s Archon: %s\n", BotMsgStyle.Render("◆"), msg.answer)
		m.viewport.SetContent(m.chatHistory)
		m.viewport.GotoBottom()
		m.lastContext = msg.context
		m.viewportContext.SetContent(m.lastContext)
		m.lastPromptTokens = msg.promptTokens
		m.lastAnswerTokens = msg.answerTokens
		m.totalTokens += msg.totalTokens
		
		// Hitung biaya
		cost := m.calculateCost(msg.promptTokens, msg.answerTokens)
		m.totalCost += cost
		
		return m, nil
	case statusTableMsg:
		m.table.SetRows(msg)
		return m, nil
	case statusMsg:
		m.statusInfo = string(msg)
		return m, nil
	case errMsg:
		m.thinking = false
		m.err = msg
		return m, nil
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.state == stateChat || m.state == stateInputPath {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.state == stateContext {
		m.viewportContext, cmd = m.viewportContext.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.thinking || m.indexing {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

type statusMsg string
type statusTableMsg []table.Row

func (m model) loadStatus() tea.Cmd {
	return func() tea.Msg {
		cfg, _ := config.LoadConfig()
		
		var rows []table.Row
		rows = append(rows, table.Row{"Model", cfg.ModelID})
		
		apiKeyStatus := "NOT Configured"
		if cfg.GeminiKey != "" {
			apiKeyStatus = "Configured"
		}
		rows = append(rows, table.Row{"API Key", apiKeyStatus})

		vectorDBStatus := "Not initialized"
		if _, err := os.Stat("./chromem_db"); err == nil {
			vectorDBStatus = "Ready (chromem_db)"
		}
		rows = append(rows, table.Row{"Vector DB", vectorDBStatus})

		// Caching status
		hash, _ := gemini.CalculateProjectHash(".")
		if hash != "" {
			rows = append(rows, table.Row{"Project Hash", hash[:8] + "..."})
		} else {
			rows = append(rows, table.Row{"Project Hash", "N/A"})
		}
		
		cacheStatus := "Inactive"
		if cfg.CacheName != "" {
			if cfg.ProjectHash == hash {
				cacheStatus = "Active (" + filepath.Base(cfg.CacheName) + ")"
			} else {
				cacheStatus = "Inactive (Hash Mismatch)"
			}
		} else {
			cacheStatus = "Inactive (Not Created)"
		}
		rows = append(rows, table.Row{"Context Cache", cacheStatus})

		// List all caches if client is available
		if cfg.GeminiKey != "" {
			client, err := gemini.NewClient(context.Background(), cfg.GeminiKey, cfg.ModelID)
			if err == nil {
				defer client.Close()
				cm := gemini.NewCacheManager(client.Client())
				caches, err := cm.ListCaches(context.Background())
				if err == nil && len(caches) > 0 {
					rows = append(rows, table.Row{"---", "---"})
					rows = append(rows, table.Row{"Remote Caches", fmt.Sprintf("%d found", len(caches))})
					for _, c := range caches {
						rows = append(rows, table.Row{"  ID", filepath.Base(c.Name)})
						rows = append(rows, table.Row{"  Model", c.Model})
					}
				}
			}
		}

		return statusTableMsg(rows)
	}
}

func (m model) clearIndex() tea.Cmd {
	return func() tea.Msg {
		err := os.RemoveAll("./chromem_db")
		if err != nil {
			return errMsg(err)
		}
		return statusMsg("Index successfully deleted. Vector database has been cleared.")
	}
}

type geminiResponseMsg struct {
	answer       string
	context      string
	promptTokens int
	answerTokens int
	totalTokens  int
}
type indexCompleteMsg string
type errMsg error

type indexProgressMsg struct {
	current int
	total   int
	file    string
}

func (m model) startIndexing() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		cfg, _ := config.LoadConfig()
		if cfg.GeminiKey == "" {
			return errMsg(fmt.Errorf("Gemini API key not found. Please use 'archon auth --key [key]' in CLI first."))
		}
		
		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		if err != nil {
			return errMsg(err)
		}
		defer store.Close()

		orchestrator := core.NewOrchestrator(store)
		
		err = orchestrator.IndexDirectory(ctx, ".", nil)

		if err != nil {
			return errMsg(err)
		}
		return indexCompleteMsg("Indexing complete. The codebase is now ready for semantic search.")
	}
}

func (m model) askGemini(query string) tea.Cmd {
	return func() tea.Msg {
		cfg, _ := config.LoadConfig()
		if cfg.GeminiKey == "" {
			return errMsg(fmt.Errorf("Gemini API key not found. Please use 'archon auth' in CLI."))
		}
		ctx := context.Background()

		// RAG Context
		store, err := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
		var prompt string
		var contextText string
		if err == nil {
			defer store.Close()
			orchestrator := core.NewOrchestrator(store)
			contextText, err = orchestrator.SearchContext(ctx, query)
			if err == nil {
				prompt = fmt.Sprintf("%s\n\nUser Question: %s", contextText, query)
			} else {
				prompt = query
			}
		} else {
			prompt = query
		}

		client, err := gemini.NewClient(ctx, cfg.GeminiKey, cfg.ModelID)
		if err != nil {
			return errMsg(err)
		}
		defer client.Close()

		// Sync and Use Context Cache
		cacheName, err := m.syncCache(ctx, client)
		if err == nil && cacheName != "" {
			client.SetCachedContent(cacheName)
		}

		resp, err := client.Ask(ctx, prompt)
		if err != nil {
			return errMsg(err)
		}
		return geminiResponseMsg{
			answer:       resp.Text,
			context:      contextText,
			promptTokens: resp.PromptTokens,
			answerTokens: resp.AnswerTokens,
			totalTokens:  resp.TotalTokens,
		}
	}
}

func (m model) syncCache(ctx context.Context, client *gemini.Client) (string, error) {
	cfg, _ := config.LoadConfig()
	hash, err := gemini.CalculateProjectHash(".")
	if err != nil {
		return "", err
	}

	if cfg.CacheName != "" && cfg.ProjectHash == hash {
		// Verify if cache still exists on server
		cm := gemini.NewCacheManager(client.Client())
		caches, err := cm.ListCaches(ctx)
		if err == nil {
			for _, c := range caches {
				if c.Name == cfg.CacheName {
					return cfg.CacheName, nil
				}
			}
		}
	}

	// Cache invalid or not found, create new one
	store, _ := vectordb.NewStore(ctx, "./chromem_db", cfg.GeminiKey)
	if store != nil {
		defer store.Close()
		orchestrator := core.NewOrchestrator(store)
		files, err := orchestrator.GetFilesForIndexing(".")
		if err == nil && len(files) > 0 {
			cm := gemini.NewCacheManager(client.Client())
			// For demo, we use existing files
			cacheName, err := cm.CreateContextCache(ctx, cfg.ModelID, files)
			if err == nil {
				viper.Set("project_hash", hash)
				viper.Set("cache_name", cacheName)
				viper.WriteConfig()
				return cacheName, nil
			} else {
				// If failed (e.g., insufficient tokens), still update hash to prevent constant retries
				// but clear cache_name
				viper.Set("project_hash", hash)
				viper.Set("cache_name", "")
				viper.WriteConfig()
			}
		}
	}

	return "", nil
}

func (m model) calculateCost(promptTokens, answerTokens int) float64 {
	cfg, _ := config.LoadConfig()
	modelID := cfg.ModelID
	
	// Gemini 3 Price Estimation (per 1M tokens)
	// Flash: Input $0.10, Output $0.40
	// Pro: Input $1.25, Output $5.00
	
	var inputPrice, outputPrice float64
	
	if strings.Contains(modelID, "flash") {
		inputPrice = 0.10 / 1000000
		outputPrice = 0.40 / 1000000
	} else if strings.Contains(modelID, "pro") {
		// Includes gemini-3-pro-preview
		inputPrice = 1.25 / 1000000
		outputPrice = 5.00 / 1000000
	} else {
		// Default to Pro price if unknown
		inputPrice = 1.25 / 1000000
		outputPrice = 5.00 / 1000000
	}
	
	return float64(promptTokens)*inputPrice + float64(answerTokens)*outputPrice
}

func (m model) drawTokenGraph(width int) string {
	if m.totalTokens == 0 || width <= 10 {
		return ""
	}
	
	// We use a 100k tokens threshold for bar visualization
	maxTokens := 100000.0
	percent := float64(m.totalTokens) / maxTokens
	if percent > 1.0 {
		percent = 1.0
	}
	
	barWidth := width
	if barWidth > 40 {
		barWidth = 40
	}
	
	filledWidth := int(percent * float64(barWidth))
	if filledWidth < 1 && m.totalTokens > 0 {
		filledWidth = 1
	}
	
	// Ensure filledWidth does not exceed barWidth
	if filledWidth > barWidth {
		filledWidth = barWidth
	}
	
	emptyWidth := barWidth - filledWidth
	if emptyWidth < 0 {
		emptyWidth = 0
	}
	
	bar := GraphBarFilledStyle.Render(strings.Repeat("█", filledWidth))
	bar += GraphBarEmptyStyle.Render(strings.Repeat("░", emptyWidth))
	
	return fmt.Sprintf("Tokens: %s %.1f%%", bar, percent*100)
}

func (m model) View() string {
	var s string

	var logo string
	if m.width > 80 {
		logo = GetLogo()
	} else {
		logo = GetCompactLogo()
	}

	logoContent := LogoStyle.Render(logo)
	geminiTag := GeminiStyle.Render("powered by gemini 3.5 pro")

	s += lipgloss.Place(m.width, 0, lipgloss.Center, lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Center, logoContent, geminiTag)) + "\n"

	switch m.state {
	case stateMenu:
		var choices string
		s += lipgloss.NewStyle().Foreground(AccentColor).Italic(true).Render("Welcome, Architect. What would you like to do today?") + "\n\n"
		for i, choice := range m.choices {
			if m.cursor == i {
				choices += SelectedItemStyle.Render("> " + choice) + "\n"
			} else {
				choices += ItemStyle.Render("  " + choice) + "\n"
			}
		}
		s += choices + "\n"
		s += FooterStyle.Render("Use arrow keys ↑↓ and Enter to select • ctrl+c to exit")

	case stateChat:
		s += StatusStyle.Render("CHAT MODE") + " (Tab: View Context)\n\n"
		s += ChatBoxStyle.Render(m.viewport.View()) + "\n\n"
		
		if m.thinking {
			s += m.spinner.View() + " Archon is thinking...\n"
		} else {
			s += m.textInput.View() + "\n"
		}
		
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
		}
		s += FooterStyle.Render("\n(esc: back • tab: context)")

	case stateContext:
		s += StatusStyle.Render("LAST RAG CONTEXT") + " (Tab: Back to Chat)\n\n"
		s += ChatBoxStyle.Render(m.viewportContext.View()) + "\n\n"
		s += FooterStyle.Render("\n(esc: back • tab: chat)")

	case stateInputPath:
		s += StatusStyle.Render("INPUT REQUIRED") + "\n\n"
		s += lipgloss.NewStyle().Foreground(AccentColor).Render(fmt.Sprintf("Action: %s", m.selectedAction)) + "\n\n"
		s += m.textInput.View() + "\n\n"
		s += FooterStyle.Render("\n(enter: run • esc: back)")

	case stateIndex:
		s += StatusStyle.Render("CODEBASE INDEXING") + "\n\n"
		if m.indexing {
			s += m.spinner.View() + " Scanning files and generating embeddings...\n"
			s += "\n" + m.progress.ViewAs(0.5) + "\n"
		} else {
			s += SelectedStyle.Render("Success: ") + m.indexStats + "\n"
		}
		if m.err != nil {
			s += ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
		}
		s += FooterStyle.Render("\n(esc to return to menu)")

	case stateStatus:
		s += StatusStyle.Render("SYSTEM STATUS") + "\n\n"
		s += TableStyle.Render(m.table.View()) + "\n"
		if m.statusInfo != "" {
			s += "\n" + m.statusInfo + "\n"
		}
		s += FooterStyle.Render("\n(esc to return to menu)")
	}

	// Token & Cost Info Bar
	if m.totalTokens > 0 {
		tokenGraph := m.drawTokenGraph(m.width - 20)
		costInfo := CostStyle.Render(fmt.Sprintf(" Estimated Cost: $%.6f", m.totalCost))
		
		sessionInfo := fmt.Sprintf(" [ Session: %d tokens | %s ]", m.totalTokens, costInfo)
		
		s += "\n\n" + tokenGraph + "\n" + TokenStyle.Render(sessionInfo)
		
		if m.lastAnswerTokens > 0 {
			lastInfo := fmt.Sprintf(" Last: %d prompt + %d answer = %d total", 
				m.lastPromptTokens, m.lastAnswerTokens, m.lastPromptTokens+m.lastAnswerTokens)
			s += "\n " + lipgloss.NewStyle().Foreground(NeutralColor).Italic(true).Render(lastInfo)
		}
	}

	// Fill screen
	if m.height > 0 {
		lines := lipgloss.Height(s)
		if lines < m.height {
			s += strings.Repeat("\n", max(0, m.height-lines))
		}
	}

	return s
}
