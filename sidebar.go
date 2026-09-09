package main

import (
	"os"
	"strings"
)

func ReadCurrentDir(path string) dirReadMessage {
	filelist, filelistErr := os.ReadDir(path)
	newFilelist := cutVisibleFiles(filelist)
	return dirReadMessage{filteredFiles: newFilelist, files: filelist, err: filelistErr}
}

func cutVisibleFiles(files []os.DirEntry) []os.DirEntry {
	newFiles := []os.DirEntry{}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		} else {
			newFiles = append(newFiles, file)
		}
	}

	return newFiles
}

func RenderSidebar(m model, files []os.DirEntry) []string {
	fileNames := []string{}

	for index, filename := range files {
		fileName := TruncateLine(filename.Name(), sidebarStyle.GetWidth())
		if index+m.fileScrollOffset == m.fileCursor && m.mode == FileMode {
			text := selectedTextStyle.Render(fileName)
			fileNames = append(fileNames, text)
		} else {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames
}
