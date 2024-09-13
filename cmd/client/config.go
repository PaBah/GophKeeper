package main

import "github.com/charmbracelet/lipgloss"

var (
	normalFg = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
	indigo   = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	cream    = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
	fuchsia  = lipgloss.Color("#F780E2")
	help     = lipgloss.Color("240")
	green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
	red      = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	blurBg   = lipgloss.Color("235")
	blurText = lipgloss.Color("240")

	titleStyle = lipgloss.NewStyle().
			Foreground(normalFg).
			Bold(true).
			MarginBottom(1)

	buttonStyle = lipgloss.NewStyle().
			Foreground(cream).
			Background(fuchsia).
			Padding(0, 3).
			MarginTop(1)

	buttonBlurredStyle = lipgloss.NewStyle().
				Foreground(blurText).
				Background(blurBg).
				Padding(0, 3).
				MarginTop(1)
)
