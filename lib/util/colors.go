package util

import "github.com/fatih/color"

// Bold returns a color that prints in bold
func Bold() *color.Color {
	return color.New(color.Bold)
}

// BoldCyan returns a color that prints in bold cyan
func BoldCyan() *color.Color {
	return color.New(color.Bold).Add(color.FgCyan)
}

// BoldRed returns a color that prints in bold red
func BoldRed() *color.Color {
	return color.New(color.Bold, color.FgRed)
}

// Red returns a color that prints in red
func Red() *color.Color {
	return color.New(color.FgRed)
}
