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
	contentHeight := GetContentHeight(m)

	// Status Bar
	status := statusBarStyle.Width(sidebarWidth + mainWidth).Render(GetModeString(m))

	// File Panel
	start := m.fileScrollOffset
	end := min(start+contentHeight, len(m.filteredFiles))
	visibleFiles := m.filteredFiles[start:end]

	fileNames := RenderSidebar(m, visibleFiles)
	formattedFileNameString := strings.Join(fileNames, "\n")

	sidebar := sidebarStyle.Height(contentHeight).Render(formattedFileNameString)

	// Main Panel
	mainContent := m.currentFileContent
	if m.fileContentErr != nil {
		mainContent = fmt.Sprintf("Error opening file:\n%s", m.fileContentErr)
	}

	lines := strings.Split(mainContent, "\n")
	start = m.cursorScrollOffset
	end = min(start + contentHeight, len(lines))
	
	if start > end {
		start = end
	}

	mainContent = strings.Join(lines[start:end], "\n")
	displayContent := RenderWithCursor(mainContent, m.cursor)
	main := mainStyle.Width(mainWidth).Height(contentHeight).Render(displayContent)

	top_panels := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)
	full_panel := lipgloss.JoinVertical(lipgloss.Left, top_panels, status)
	return full_panel
}
