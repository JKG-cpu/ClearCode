package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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

func InitModelMode() (string, Mode, error) {
	if len(os.Args) == 1 {
		return ".", FileMode, nil
	}

	filepathArg := filepath.Join(".", os.Args[1])

	info, err := os.Stat(filepathArg)

	if err != nil {
		return ".", NormalMode, err
	}

	if info.IsDir() {
		return filepathArg, FileMode, nil
	} else {
		return filepathArg, NormalMode, nil
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

func ShortenVerticalLines(content string, maxLines int) string {
	lines := strings.Split(content, "\n")
	if len(lines) <= maxLines || maxLines <= 0 {
		return content
	}
	return strings.Join(lines[:maxLines], "\n")
}