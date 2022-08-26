package cli

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kctjohnson/bubble-boids/internal/boid"
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
	"golang.org/x/term"
)

type TickMsg time.Time

type Model struct {
	keys           keyMap
	help           help.Model
	screen         *Screen
	virtualScreen  *VirtualScreen
	boids          *[]*boid.Boid
	scatterCounter int // Starts at 0, when it hits ScatterCounterCap, all of the boids are scattered
	viewStyle      lipgloss.Style
}

func InitialModel(width int, height int) Model {
	newBoidSlice := new([]*boid.Boid)
	*newBoidSlice = make([]*boid.Boid, 0)

	return Model{ // Subtract the help height to make space for the help UI
		virtualScreen:  NewVirtualScreen(width-BorderPadding, height-HelpHeight-BorderPadding),
		screen:         NewScreen(width-BorderPadding, height-HelpHeight-BorderPadding),
		boids:          newBoidSlice,
		scatterCounter: 0,
		keys:           keys,
		help:           help.New(),
		viewStyle:      lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#ffffff")),
	}
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second/FPS, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Frame() (tea.Model, tea.Cmd) {
	m.scatterCounter++
	for _, b := range *m.boids {
		if m.scatterCounter >= boid.ScatterCounterCap {
			// Randomize velocity and acceleration
			b.Velocity = mathutil.RandomVec2(-boid.MaxSpeed, boid.MaxSpeed)
			b.Acceleration = mathutil.RandomVec2(-boid.MaxSpeed, boid.MaxSpeed)
		} else {
			b.Edges(m.virtualScreen.Width, m.virtualScreen.Height)
			b.Flock(m.boids)
		}
		b.Update()
	}

	if m.scatterCounter >= boid.ScatterCounterCap {
		m.scatterCounter = 0
	}
	return m, m.tick()
}

func (m Model) Init() tea.Cmd {
	numberOfBoids := boid.BoidCount
	for i := 0; i < numberOfBoids; i++ {
		*m.boids = append(*m.boids, boid.NewBoid(i, m.virtualScreen.Width, m.virtualScreen.Height))
	}
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		return m.Frame()
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

		// Update both the main screen, and the virtual screen
		m.virtualScreen.UpdateScreenSize(msg.Width-BorderPadding, msg.Height-HelpHeight-BorderPadding)
		m.screen.UpdateScreenSize(msg.Width-BorderPadding, msg.Height-HelpHeight-BorderPadding)
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.ModifyAlignment):
			if key.Matches(msg, m.keys.IncreaseAlignment) {
				boid.MaxAlignmentForce += 0.1
			} else {
				if boid.MaxAlignmentForce > 0.0 {
					boid.MaxAlignmentForce -= 0.1
				}
			}

		case key.Matches(msg, m.keys.ModifyCohesion):
			if key.Matches(msg, m.keys.IncreaseCohesion) {
				boid.MaxCohesionForce += 0.1
			} else {
				if boid.MaxCohesionForce > 0.0 {
					boid.MaxCohesionForce -= 0.1
				}
			}

		case key.Matches(msg, m.keys.ModifySeparation):
			if key.Matches(msg, m.keys.IncreaseSeparation) {
				boid.MaxSeparationForce += 0.1
			} else {
				if boid.MaxSeparationForce > 0.0 {
					boid.MaxSeparationForce -= 0.1
				}
			}

		case key.Matches(msg, m.keys.ModifyPerception):
			if key.Matches(msg, m.keys.IncreasePerception) {
				boid.Perception += 1
			} else {
				if boid.Perception > 0.0 {
					boid.Perception -= 1
				}
			}

		case key.Matches(msg, m.keys.ModifyMaxSpeed):
			if key.Matches(msg, m.keys.IncreaseMaxSpeed) {
				boid.MaxSpeed += 0.1
			} else {
				if boid.MaxSpeed > 0.0 {
					boid.MaxSpeed -= 0.1
				}
			}

		case key.Matches(msg, m.keys.Scatter):
			m.scatterCounter = boid.ScatterCounterCap

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	m.screen.Clear()

	for _, b := range *m.boids {
		// Convert the current virtual boid position to a terminal coordinate
		termX := int(math.Floor(float64(m.screen.Width) / m.virtualScreen.Width * b.Position.X()))
		termY := int(math.Floor(float64(m.screen.Height) / m.virtualScreen.Height * b.Position.Y()))
		if termX >= m.screen.Width {
			termX = m.screen.Width - 1
		}
		if termY >= m.screen.Height {
			termY = m.screen.Height - 1
		}
		if termX < 0 {
			termX = 0
		}
		if termY < 0 {
			termY = 0
		}
		m.screen.SetRune(termX, termY, '*')
	}

	screen := m.screen.GetScreen()
	screenView := m.viewStyle.Render(screen)
	screenView += fmt.Sprintf("\nAlignment: %f | Cohesion: %f | Separation: %f | Perception: %d | Speed: %f", boid.MaxAlignmentForce, boid.MaxCohesionForce, boid.MaxSeparationForce, boid.Perception, boid.MaxSpeed)
	screenView += "\n" + m.help.View(m.keys)

	return screenView
}

func Execute() {
	rand.Seed(time.Now().UTC().UnixNano())

	width, height, err := term.GetSize(0)
	if err != nil {
		panic("Yikes")
	}

	p := tea.NewProgram(InitialModel(width, height), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
