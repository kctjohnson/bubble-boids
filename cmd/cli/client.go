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
	"golang.org/x/term"
)

type TickMsg time.Time

type Model struct {
	keys          keyMap
	help          help.Model
	screen        *Screen
	virtualScreen *VirtualScreen
	viewStyle     lipgloss.Style
	flock         *boid.Flock
}

func InitialModel(width int, height int) Model {
	newVirtualScreen := NewVirtualScreen(width-BorderPadding, height-HelpHeight-BorderPadding)
	return Model{ // Subtract the help height to make space for the help UI
		virtualScreen: newVirtualScreen,
		screen:        NewScreen(width-BorderPadding, height-HelpHeight-BorderPadding),
		flock:         boid.NewFlock(newVirtualScreen.Width, newVirtualScreen.Height),
		keys:          keys,
		help:          help.New(),
		viewStyle:     lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#ffffff")),
	}
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second/FPS, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Frame() (tea.Model, tea.Cmd) {
	m.flock.Update(m.virtualScreen.Width, m.virtualScreen.Height)
	return m, m.tick()
}

func (m Model) Init() tea.Cmd {
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
		const adjustAmount = 0.01
		switch {
		case key.Matches(msg, m.keys.ModifyAlignment):
			if key.Matches(msg, m.keys.IncreaseAlignment) {
				m.flock.BoidSettings.MaxAlignmentForce += adjustAmount
			} else {
				if m.flock.BoidSettings.MaxAlignmentForce > 0.0 {
					m.flock.BoidSettings.MaxAlignmentForce -= adjustAmount
				}
			}

		case key.Matches(msg, m.keys.ModifyCohesion):
			if key.Matches(msg, m.keys.IncreaseCohesion) {
				m.flock.BoidSettings.MaxCohesionForce += adjustAmount
			} else {
				if m.flock.BoidSettings.MaxCohesionForce > 0.0 {
					m.flock.BoidSettings.MaxCohesionForce -= adjustAmount
				}
			}

		case key.Matches(msg, m.keys.ModifySeparation):
			if key.Matches(msg, m.keys.IncreaseSeparation) {
				m.flock.BoidSettings.MaxSeparationForce += adjustAmount
			} else {
				if m.flock.BoidSettings.MaxSeparationForce > 0.0 {
					m.flock.BoidSettings.MaxSeparationForce -= adjustAmount
				}
			}

		case key.Matches(msg, m.keys.ModifyPerception):
			if key.Matches(msg, m.keys.IncreasePerception) {
				m.flock.BoidSettings.Perception += 1
			} else {
				if m.flock.BoidSettings.Perception > 0.0 {
					m.flock.BoidSettings.Perception -= 1
				}
			}

		case key.Matches(msg, m.keys.ModifyMaxSpeed):
			if key.Matches(msg, m.keys.IncreaseMaxSpeed) {
				m.flock.BoidSettings.MaxSpeed += adjustAmount
			} else {
				if m.flock.BoidSettings.MaxSpeed > 0.0 {
					m.flock.BoidSettings.MaxSpeed -= adjustAmount
				}
			}

		case key.Matches(msg, m.keys.Scatter):
			m.flock.Scatter()

		case key.Matches(msg, m.keys.ToggleEdgeMode):
			if m.flock.BoidSettings.EdgeMode == boid.EDGE_AVOID {
				m.flock.BoidSettings.EdgeMode = boid.EDGE_WARP
			} else {
				m.flock.BoidSettings.EdgeMode = boid.EDGE_AVOID
			}

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	m.screen.Clear()

	for _, b := range m.flock.Boids {
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
	screenView += fmt.Sprintf("\nAlignment: %.2f | Cohesion: %.2f | Separation: %.2f | Perception: %d | Speed: %.2f", m.flock.BoidSettings.MaxAlignmentForce, m.flock.BoidSettings.MaxCohesionForce, m.flock.BoidSettings.MaxSeparationForce, m.flock.BoidSettings.Perception, m.flock.BoidSettings.MaxSpeed)
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
