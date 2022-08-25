package cli

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kctjohnson/bubble-boids/internal/boid"
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
	"golang.org/x/term"
)

type TickMsg time.Time

type model struct {
	screen         *Screen
	virtualScreen  *VirtualScreen
	boids          *[]*boid.Boid
	scatterCounter int // Starts at 0, when it hits ScatterCounterCap, all of the boids are scattered
}

func initialModel() model {
	width, height, err := term.GetSize(0)
	if err != nil {
		panic("Yikes")
	}

	newBoidSlice := new([]*boid.Boid)
	*newBoidSlice = make([]*boid.Boid, 0)

	return model{
		screen:         NewScreen(width, height),
		boids:          newBoidSlice,
		scatterCounter: 0,
	}
}

func (m model) tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Frame() (tea.Model, tea.Cmd) {
	m.scatterCounter++
	for _, b := range *m.boids {
		if m.scatterCounter >= boid.ScatterCounterCap {
			// Randomize velocity and acceleration
			b.Velocity = mathutil.RandomVec2(-boid.MaxSpeed, boid.MaxSpeed)
			b.Acceleration = mathutil.RandomVec2(-boid.MaxSpeed, boid.MaxSpeed)
		} else {
			b.Edges(m.screen.Width, m.screen.Height)
			b.Flock(m.boids)
		}
		b.Update()
	}

	if m.scatterCounter >= boid.ScatterCounterCap {
		m.scatterCounter = 0
	}
	return m, m.tick()
}

func (m model) Init() tea.Cmd {
	numberOfBoids := boid.BoidCount
	for i := 0; i < numberOfBoids; i++ {
		*m.boids = append(*m.boids, boid.NewBoid(i, m.screen.Width, m.screen.Height))
	}
	return m.tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		return m.Frame()
	case tea.WindowSizeMsg:
		m.screen.UpdateScreenSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	m.screen.Clear()
	for _, b := range *m.boids {
		posX := int(b.Position.X())
		posY := int(b.Position.Y())
		if posX >= m.screen.Width {
			posX = m.screen.Width - 1
		}
		if posY >= m.screen.Height {
			posY = m.screen.Height - 1
		}
		if posX < 0 {
			posX = 0
		}
		if posY < 0 {
			posY = 0
		}
		m.screen.SetRune(posX, posY/boid.TermRatio, '*')
	}
	return m.screen.GetScreen()
}

func Execute() {
	rand.Seed(time.Now().UTC().UnixNano())
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
