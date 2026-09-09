package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func ClearTerminal() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func InitModelMode() (string, string, Mode, error) {
    if len(os.Args) == 1 {
        return ".", "", FileMode, nil
    }

    filepathArg := filepath.Join(".", os.Args[1])
    info, err := os.Stat(filepathArg)

    if err != nil {
        return ".", "", NormalMode, err
    }

    if info.IsDir() {
        return filepathArg, "", FileMode, nil
    } else {
        return filepath.Dir(filepathArg), filepathArg, NormalMode, nil
    }
}

func GetModeString(m model) string {
	if m.mode == 1 {
		return "INSERT"
	}

	if m.mode == 2 {
		return "FILE"
	}

	return "NORMAL"
}

func ReadFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	return string(bytes), err
}

func ResetCursors(m model) model {
	// Text Cursor
	m.cursor = [2]int{0, 0}
	m.cursorScrollOffset = 0
	m.horizontalScrollOffset = 0
	m.desiredCol = 0

	// File Cursor
	m.fileCursor = 0
	m.fileScrollOffset = 0

	return m
}

func ShortenVerticalLines(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines || maxLines <= 0 {
		return content
	}
	return strings.Join(lines[:maxLines], "\n")
}

func TruncateLine(s string, width int) string {
	runes := []rune(s)
	if len(runes) <= width || width <= 0 {
		return s
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

func GetContentHeight(m model) int {
	sidebarWidth := sidebarStyle.GetWidth() + sidebarStyle.GetHorizontalFrameSize()
	mainWidth := m.width - sidebarWidth - mainStyle.GetHorizontalFrameSize()

	status := statusBarStyle.Width(sidebarWidth + mainWidth).Render(GetModeString(m))
	statusHeight := lipgloss.Height(status)

	return m.height - sidebarStyle.GetVerticalFrameSize() - statusHeight
}

func GetMainPanelWidth(m model) int {
	sidebarWidth := sidebarStyle.GetWidth() + sidebarStyle.GetHorizontalFrameSize()
	return m.width - sidebarWidth - mainStyle.GetHorizontalFrameSize()
}

func WindowLine(expandedLine string, offset int, width int) string {
    runes := []rune(expandedLine)

    if offset >= len(runes) {
        return ""
    }

    end := min(offset + width, len(runes))

    return string(runes[offset:end])
}

func expandTabsWithMap(line string, tabWidth int) (string, []int) {
	var expanded strings.Builder
	visualCol := 0
	runeToVisual := make([]int, 0, len(line)+1)

	for _, r := range line {
		runeToVisual = append(runeToVisual, visualCol)
		if r == '\t' {
			spaces := tabWidth - (visualCol % tabWidth)
			expanded.WriteString(strings.Repeat(" ", spaces))
			visualCol += spaces
		} else {
			expanded.WriteRune(r)
			visualCol++
		}
	}
	runeToVisual = append(runeToVisual, visualCol)

	return expanded.String(), runeToVisual
}

func clampCursorCol(m model) int {
	lines := strings.Split(m.currentFileContent, "\n")
	lineLength := len([]rune(lines[m.cursor[0]]))
	if m.cursor[1] > lineLength {
		return lineLength
	}
	return m.cursor[1]
}

func RenderCursorAtCol(line string, col int) string {
    runes := []rune(line)

    if col < 0 {
        return line
    }
    if col >= len(runes) {
        return line + cursorStyle.Render(" ")
    }

    before := string(runes[:col])
    char := string(runes[col])
    after := string(runes[col+1:])

    return before + cursorStyle.Render(char) + after
}

func RenderWithCursor(content string, cursor [2]int) string {
	lines := strings.Split(content, "\n")

	cursorRow := cursor[0]
	cursorCol := cursor[1]

	if cursorRow < 0 || cursorRow >= len(lines) {
		return content
	}

	line := lines[cursorRow]
	expandedLine, runeToVisual := expandTabsWithMap(line, 4)

	if cursorCol < 0 || cursorCol >= len(runeToVisual) {
		return content
	}

	visualCol := runeToVisual[cursorCol]
	runes := []rune(expandedLine)

	if visualCol >= len(runes) {
		lines[cursorRow] = expandedLine + cursorStyle.Render(" ")
		return strings.Join(lines, "\n")
	}

	before := string(runes[:visualCol])
	char := string(runes[visualCol])
	after := string(runes[visualCol+1:])

	lines[cursorRow] = before + cursorStyle.Render(char) + after

	return strings.Join(lines, "\n")
}
