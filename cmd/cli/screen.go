package cli

import (
	"fmt"
	"strings"
)

// Virtual screen is made based on the size of the terminal screen in rows and columns.
// The virtual screen will be TERM_RATIO times taller than the main screen, since the
// terminal characters are taller than they are wide (generally)

// First step will be to get the ratio of the terminal
// Next step will be to create a sceen at least VIRTUAL_SCREEN_WIDTH wide,
// 	then get the height based on the ratio, multiplied by TERM_RATIO
// The boids will operate inside of the virtual screen dimensions,
// 	then be rendered down into the terminal screen on draw

type VirtualScreen struct {
	Width  int
	Height int
	Zoom   float64
}

type Screen struct {
	Runes  []rune
	Width  int
	Height int
}

func NewScreen(width int, height int) *Screen {
	var newScreen = new(Screen)
	newScreen.Runes = make([]rune, width*height)
	newScreen.Width = width
	newScreen.Height = height
	newScreen.Clear()
	return newScreen
}

func (s Screen) GetScreen() string {
	finalScreen := ""
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			finalScreen += string(s.Runes[x+y*s.Width])
		}
		finalScreen += "\n"
	}
	finalScreen += fmt.Sprintf("Max Force: %f | Max Speed: %f | Boid Qty: %d", MaxForce, MaxSpeed, BoidCount)
	return strings.Trim(finalScreen, "\n")
}

func (s *Screen) UpdateScreenSize(width int, height int) {
	s.Runes = make([]rune, width*height)
	s.Width = width
	s.Height = height
}

func (s *Screen) Clear() {
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.SetRune(x, y, ' ')
		}
	}
}

func (s *Screen) SetRune(x int, y int, c rune) {
	s.Runes[x+y*s.Width] = c
}
