package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Terminal
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

// Init
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

// Modes
func GetModeString(m model) string {
	if m.mode == 1 {
		return "INSERT"
	}

	if m.mode == 2 {
		return "FILE"
	}

	return "NORMAL"
}

// Files
func ReadFile(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	return string(bytes), err
}

func InsertRune(m model, r rune) model {
	cursorRow := m.cursor[0]
	cursorCol := m.cursor[1]

	runes := []rune(m.lines[cursorRow])

	var builder strings.Builder
	builder.WriteString(string(runes[:cursorCol]))
	builder.WriteRune(r)
	builder.WriteString(string(runes[cursorCol:]))

	m.lines[cursorRow] = builder.String()

	m.cursor[1]++
	m.desiredCol = m.cursor[1]
	
	return m
}

func InsertEmptyLine(m model) model {
    currentLine := []rune(m.lines[m.cursor[0]])

    before := string(currentLine[:m.cursor[1]])
    after := string(currentLine[m.cursor[1]:])

    m.lines[m.cursor[0]] = before
    m.lines = slices.Insert(m.lines, m.cursor[0]+1, after)

    m.cursor[0]++
    m.cursor[1] = 0
    m.desiredCol = 0

    return m
}

func DeleteCharacter(m model) model {
	if m.cursor == [2]int{0, 0} {
		return m
	}

	// Cursor at start of line
	if m.cursor[1] == 0 {
		prevLine := m.lines[m.cursor[0] - 1]
		currentLine := m.lines[m.cursor[0]][m.cursor[1]:]

		var builder strings.Builder
		builder.WriteString(prevLine)
		builder.WriteString(currentLine)

		m.lines[m.cursor[0] - 1] = builder.String()
		m.lines = slices.Delete(m.lines, m.cursor[0], m.cursor[0] + 1)
		m.cursor[0]--
		m.cursor[1] = len([]rune(prevLine))
	} else {
		cursorRow := m.cursor[0]
		cursorCol := m.cursor[1]

		runes := []rune(m.lines[cursorRow])
		m.lines[cursorRow] = string(runes[:cursorCol-1]) + string(runes[cursorCol:])	
		m.cursor[1] = max(0, m.cursor[1] - 1)
	}

	m.desiredCol = m.cursor[1]
	
	return m
}

func SaveFile(m model) model {
	content := strings.Join(m.lines, "\n")
	err := os.WriteFile(m.currentFile, []byte(content), 0644)
	if err != nil {
		m.saveErr = err
	}

	return m
}

// Cursors
func ResetCursors(m model) model {
	// Text Cursor
	m.cursor = [2]int{0, 0}
	m.cursorScrollOffset = 0
	m.horizontalScrollOffset = 0
	m.desiredCol = 0
	m.endLine = false

	// File Cursor
	m.fileCursor = 0
	m.fileScrollOffset = 0

	return m
}

// Lines
func ShortenVerticalLines(content []string, maxLines int) string {
	if len(content) <= maxLines || maxLines <= 0 {
		return strings.Join(content, "\n")
	}
	return strings.Join(content[:maxLines], "\n")
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

// Get
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

func GetVisualCol(m model) int {
	_, runeToVisual := expandTabsWithMap(m.lines[m.cursor[0]], 4)
	return runeToVisual[m.cursor[1]]
}

func WindowLine(expandedLine string, offset int, width int) string {
	runes := []rune(expandedLine)

	if offset >= len(runes) {
		return ""
	}

	end := min(offset+width, len(runes))

	return string(runes[offset:end])
}

// Rendering
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
	lineLength := len([]rune(m.lines[m.cursor[0]]))
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

func RenderWithCursor(content []string, cursor [2]int) string {
	cursorRow := cursor[0]
	cursorCol := cursor[1]

	if cursorRow < 0 || cursorRow >= len(content) {
		return strings.Join(content, "\n")
	}

	line := content[cursorRow]
	expandedLine, runeToVisual := expandTabsWithMap(line, 4)

	if cursorCol < 0 || cursorCol >= len(runeToVisual) {
		return strings.Join(content, "\n")
	}

	visualCol := runeToVisual[cursorCol]
	runes := []rune(expandedLine)

	if visualCol >= len(runes) {
		content[cursorRow] = expandedLine + cursorStyle.Render(" ")
		return strings.Join(content, "\n")
	}

	before := string(runes[:visualCol])
	char := string(runes[visualCol])
	after := string(runes[visualCol+1:])

	content[cursorRow] = before + cursorStyle.Render(char) + after

	return strings.Join(content, "\n")
}
