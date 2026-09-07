package main

import (
	"path/filepath"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Normal Mode
		if m.mode == NormalMode {
			// Quitting
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}

			// Cursor Movement
			if msg.String() == "j" || msg.String() == "down" {
				if m.cursor >= len(m.filteredFiles) - 1 {
					m.cursor = 0
				} else {
					m.cursor++
				}
			}

			if msg.String() == "k" || msg.String() == "up" {
				if m.cursor <= 0 {
					m.cursor = len(m.filteredFiles) - 1
				} else {
					m.cursor--
				}
			}

			// Switching Modes
			if msg.String() == "f" {
				m.mode = FileMode
			}

			if msg.String() == "i" {
				m.mode = InsertMode
			}
		}
		
		// Insert Mode
		if m.mode == InsertMode {
			if msg.String() == "esc" {
				m.mode = NormalMode
			}
		}

		// File Mode
		if m.mode == FileMode {
			if msg.String() == "esc" {
				m.mode = NormalMode
			}

			// Cursor Movement
			if msg.String() == "j" || msg.String() == "down" {
				if m.fileCursor >= len(m.filteredFiles) - 1 {
					m.fileCursor = 0
				} else {
					m.fileCursor++
				}
			}

			if msg.String() == "k" || msg.String() == "up" {
				if m.fileCursor <= 0 {
					m.fileCursor = len(m.filteredFiles) - 1
				} else {
					m.fileCursor--
				}
			}

			// Selecting Folders + Files
			if msg.String() == "enter" || msg.String() == "l" || msg.String() == "right" {
				if len(m.filteredFiles) == 0 {
					return m, nil
				}

				file := m.filteredFiles[m.fileCursor]
				if file.IsDir() {
					m.fileCursor = 0
					m.currentPath = filepath.Join(m.currentPath, file.Name())
					return m, func() tea.Msg {
						return ReadCurrentDir(m.currentPath)
					}
				} else {
					newPath := filepath.Join(m.currentPath, file.Name())
					bytes, err := ReadFile(newPath)

					m.currentFileContent = string(bytes)
					m.fileContentErr = err
					m.mode = NormalMode
				}
			}

			if msg.String() == "backspace" || msg.String() == "h" || msg.String() == "left" {
				m.fileCursor = 0
				m.currentPath = filepath.Join(m.currentPath, "..")
				return m, func() tea.Msg {
					return ReadCurrentDir(m.currentPath)
				}
			}
		}
		
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case dirReadMessage:
		if msg.err != nil {
			m.dirErr = msg.err
			return m, tea.Quit
		}
		m.filteredFiles = msg.filteredFiles
		m.files = msg.files
		
	case fileReadMessage:
		m.currentFileContent = msg.content
		m.fileContentErr = msg.err

		// Read Parent Folder as well
		dir := ReadCurrentDir(filepath.Join(m.currentPath, ".."))

		m.filteredFiles = dir.filteredFiles
		m.files = dir.files
		m.dirErr = dir.err
	}

	return m, nil
}