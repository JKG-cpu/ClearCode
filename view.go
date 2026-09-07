package main

import (
	"fmt"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	sidebarWidth := sidebarStyle.GetWidth() + sidebarStyle.GetHorizontalFrameSize()
	mainWidth := m.width - sidebarWidth - mainStyle.GetHorizontalFrameSize()
	
	status := statusBarStyle.Width(sidebarWidth + mainWidth).Render(GetModeString(m))
	statusHeight := lipgloss.Height(status)

	contentHeight := m.height - sidebarStyle.GetVerticalFrameSize() - statusHeight

	fileNames := RenderSidebar(m, m.filteredFiles)
	formattedFileNameString := strings.Join(fileNames, "\n")

	sidebar := sidebarStyle.Height(contentHeight).Render(formattedFileNameString)

	mainContent := m.currentFileContent
	if m.fileContentErr != nil {
		mainContent = fmt.Sprintf("Error opening file:\n%s", m.fileContentErr)
	}
	
	mainContent = ShortenVerticalLines(mainContent, contentHeight)
	main := mainStyle.Width(mainWidth).Height(contentHeight).Render(mainContent)

	top_panels := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	full_panel := lipgloss.JoinVertical(lipgloss.Left, top_panels, status)
	return full_panel
}