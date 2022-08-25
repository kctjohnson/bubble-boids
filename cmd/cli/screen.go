package cli

import (
	"strings"

)

type VirtualScreen struct {
	Width  float64 // Working with floats since the conversion between terminal
	Height float64 // to virtual won't always result in round numbers
}

func NewVirtualScreen(width int, height int) *VirtualScreen {
	var newVirtualScreen = new(VirtualScreen)
	newVirtualScreen.Width = VirtScreenWidth
	newVirtualScreen.Height = (newVirtualScreen.Width / float64(width) * float64(height)) * TermRatio
	return newVirtualScreen
}

func (vs *VirtualScreen) UpdateScreenSize(width int, height int) {
	vs.Width = VirtScreenWidth
	vs.Height = (vs.Width / float64(width) * float64(height)) * TermRatio
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
