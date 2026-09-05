package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sidebarStyle = lipgloss.NewStyle().
	Width(24).
	Border(lipgloss.NormalBorder())

var mainStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder())

type model struct {
	width int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q"{
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	sidebarTotal := sidebarStyle.GetWidth() + sidebarStyle.GetHorizontalFrameSize()
	mainWidth := m.width - sidebarTotal - mainStyle.GetHorizontalFrameSize()
	contentHeight := m.height - sidebarStyle.GetVerticalFrameSize()

	sidebar := sidebarStyle.Height(contentHeight).Render("Sidebar")
	main := mainStyle.Width(mainWidth).Height(contentHeight).Render("Main")

	output := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	return output
}

func main() {
	program := tea.NewProgram(model{}, tea.WithAltScreen())
	program.Run()
}
