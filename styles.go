package main

import "github.com/charmbracelet/lipgloss"

var sidebarStyle = lipgloss.NewStyle().
	Width(24).
	Border(lipgloss.NormalBorder())

var mainStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder())

var statusBarStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder())

var selectedTextStyle = lipgloss.NewStyle().
	Bold(true).
	Italic(true)

var cursorStyle = lipgloss.NewStyle().
    Background(lipgloss.Color("205"))