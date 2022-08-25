package screen

import (
	"fmt"
	"github.com/kctjohnson/bubble-boids/cmd/cli/utils"
	"strings"
)

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
	finalScreen += fmt.Sprintf("Max Force: %f | Max Speed: %f | Boid Qty: %d", utils.MaxForce, utils.MaxSpeed, utils.BoidCount)
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
