package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Structs
type Mode int

const (
	NormalMode Mode = iota
	InsertMode
	FileMode
)

type model struct {
	width int
	height int
	cursor [2]int
	cursorScrollOffset int
	fileCursor int
	fileScrollOffset int
	mode Mode
	filteredFiles []os.DirEntry
	files []os.DirEntry
	currentPath string
	currentFileContent string
	dirErr error
	fileContentErr error
}

type dirReadMessage struct {
	filteredFiles []os.DirEntry
	files []os.DirEntry
	err error
}

type fileReadMessage struct {
	content string
	err error
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		if m.mode == NormalMode {
			content, err := ReadFile(m.currentPath)
			return fileReadMessage{content: content, err: err}
		}
		return ReadCurrentDir(m.currentPath)
	}
}

func main() {
	ClearTerminal()

	path, startingMode, err := InitModelMode()

	if err != nil {
		fmt.Printf("Invalid Filepath:\n%s", err)
		os.Exit(1)
	}

	program := tea.NewProgram(model{mode: startingMode, currentPath: path}, tea.WithAltScreen())
	finalModel, err := program.Run()

	if err != nil {
		fmt.Printf("Error running program:\n%s", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok && m.dirErr != nil{
		fmt.Println("Error:", m.dirErr)
		os.Exit(1)
	}

}
