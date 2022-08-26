package cli

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	ModifyAlignment  key.Binding
	ModifyCohesion   key.Binding
	ModifySeparation key.Binding
	ModifyPerception key.Binding
	ModifyMaxSpeed   key.Binding

	IncreaseAlignment  key.Binding
	DecreaseAlignment  key.Binding
	IncreaseCohesion   key.Binding
	DecreaseCohesion   key.Binding
	IncreaseSeparation key.Binding
	DecreaseSeparation key.Binding
	IncreasePerception key.Binding
	DecreasePerception key.Binding
	IncreaseMaxSpeed   key.Binding
	DecreaseMaxSpeed   key.Binding
	Scatter            key.Binding

	Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.ModifyAlignment,
		k.ModifyCohesion,
		k.ModifySeparation,
		k.ModifyPerception,
		k.ModifyMaxSpeed,
		k.Scatter,
		k.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.IncreaseAlignment, k.DecreaseAlignment},
		{k.IncreaseCohesion, k.DecreaseCohesion},
		{k.IncreaseSeparation, k.DecreaseSeparation},
		{k.IncreasePerception, k.DecreasePerception},
		{k.IncreaseMaxSpeed, k.DecreaseMaxSpeed},
		{k.Scatter},
		{k.Quit},
	}
}

var keys = keyMap{
	ModifyAlignment:    key.NewBinding(key.WithKeys("a", "z"), key.WithHelp("a/z", "Modify Alignment")),
	ModifyCohesion:     key.NewBinding(key.WithKeys("s", "x"), key.WithHelp("s/x", "Modify Cohesion")),
	ModifySeparation:   key.NewBinding(key.WithKeys("d", "c"), key.WithHelp("d/c", "Modify Separation")),
	ModifyPerception:   key.NewBinding(key.WithKeys("f", "v"), key.WithHelp("f/v", "Modify Perception")),
	ModifyMaxSpeed:     key.NewBinding(key.WithKeys("g", "b"), key.WithHelp("g/b", "Modify MaxSpeed")),
	IncreaseAlignment:  key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "Increase Alignment")),
	DecreaseAlignment:  key.NewBinding(key.WithKeys("z"), key.WithHelp("z", "Decrease Alignment")),
	IncreaseCohesion:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Increase Cohesion")),
	DecreaseCohesion:   key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "Decrease Cohesion")),
	IncreaseSeparation: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "Increase Separation")),
	DecreaseSeparation: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "Decrease Separation")),
	IncreasePerception: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "Increase Perception")),
	DecreasePerception: key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "Decrease Perception")),
	IncreaseMaxSpeed:   key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "Increase Max Speed")),
	DecreaseMaxSpeed:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "Decrease Max Speed")),
	Scatter:            key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "Scatter Boids")),
	Quit:               key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "Quit")),
}
