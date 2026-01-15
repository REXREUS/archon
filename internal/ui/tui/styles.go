package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor   = lipgloss.Color("#00E5FF") // Cyan Neon
	SecondaryColor = lipgloss.Color("#FF00FF") // Magenta
	AccentColor    = lipgloss.Color("#7C3AED") // Purple
	TertiaryColor  = lipgloss.Color("#F472B6") // Pink
	BotColor       = lipgloss.Color("#00FF9F") // Neon Green
	UserColor      = lipgloss.Color("#FFD93D") // Golden Yellow
	ErrorColor     = lipgloss.Color("#FF4757") // Soft Red
	NeutralColor   = lipgloss.Color("#6B7280") // Gray
	DimColor       = lipgloss.Color("#374151") // Dark Gray

	LogoStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	GeminiStyle = lipgloss.NewStyle().
		Foreground(TertiaryColor).
		Italic(true)

	// Subtitle yang lebih halus
	SubtitleStyle = lipgloss.NewStyle().
		Foreground(NeutralColor).
		Italic(true)

	// Header container
	HeaderContainer = lipgloss.NewStyle().
		Padding(0, 0).
		MarginBottom(1).
		Align(lipgloss.Center)

	// Welcome message style
	WelcomeStyle = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true).
		MarginLeft(2).
		MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(PrimaryColor).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)

	FooterStyle = lipgloss.NewStyle().
		Foreground(DimColor).
		Italic(true).
		Padding(1, 0).
		MarginTop(1)

	CursorStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(BotColor).
		Bold(true)

	ItemStyle = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.Color("#9CA3AF"))

	SelectedItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(PrimaryColor).
		Background(lipgloss.Color("#1F2937")).
		Bold(true).
		Width(45)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true).
		Padding(1, 0)

	ChatBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor).
		Padding(1, 2)

	BotMsgStyle = lipgloss.NewStyle().
		Foreground(BotColor).
		Bold(true)

	UserMsgStyle = lipgloss.NewStyle().
		Foreground(UserColor).
		Bold(true)

	StatusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(AccentColor).
		Bold(true).
		Padding(0, 2)

	TableStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(AccentColor)

	TokenStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	GraphBarFilledStyle = lipgloss.NewStyle().Foreground(BotColor)
	GraphBarEmptyStyle  = lipgloss.NewStyle().Foreground(DimColor)
	CostStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#FBBF24")).Bold(true)
)

// GetMenuIcon returns icon for each menu item
func GetMenuIcon(choice string) string {
	icons := map[string]string{
		"Chat Mode":              "💬",
		"Index Codebase":         "📂",
		"AI Code Review":         "🔍",
		"Smart Commit":           "✍️",
		"Explain File/Symbol":    "📖",
		"Refactor Code":          "🔧",
		"Generate Unit Tests":    "🧪",
		"Architectural Analysis": "🏗️",
		"Generate Diagram":       "📊",
		"System Status":          "⚙️",
		"Clear Index":            "🗑️",
		"Exit":                   "🚪",
	}
	if icon, ok := icons[choice]; ok {
		return icon
	}
	return "•"
}

// GetLogo returns ASCII art logo
func GetLogo() string {
	return `
  ╔═══════════════════════════════════════════════════════════════╗
  ║                                                               ║
  ║      █████╗ ██████╗  ██████╗██╗  ██╗ ██████╗ ███╗   ██╗       ║
  ║     ██╔══██╗██╔══██╗██╔════╝██║  ██║██╔═══██╗████╗  ██║       ║
  ║     ███████║██████╔╝██║     ███████║██║   ██║██╔██╗ ██║       ║
  ║     ██╔══██║██╔══██╗██║     ██╔══██║██║   ██║██║╚██╗██║       ║
  ║     ██║  ██║██║  ██║╚██████╗██║  ██║╚██████╔╝██║ ╚████║       ║
  ║     ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝       ║
  ║                                                               ║
  ╚═══════════════════════════════════════════════════════════════╝`
}

// GetCompactLogo returns smaller logo for narrow terminals
func GetCompactLogo() string {
	return `
  ┌─────────────────────────────────────┐
  │  ▄▀█ █▀█ █▀▀ █░█ █▀█ █▄░█          │
  │  █▀█ █▀▄ █▄▄ █▀█ █▄█ █░▀█  CLI     │
  └─────────────────────────────────────┘`
}