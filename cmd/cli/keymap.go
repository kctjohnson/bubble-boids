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

	ToggleEdgeMode key.Binding

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
		k.ToggleEdgeMode,
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
		{k.Scatter, k.ToggleEdgeMode},
		{k.Quit},
	}
}

var keys = keyMap{
	ModifyAlignment:    key.NewBinding(key.WithKeys("a", "z"), key.WithHelp("a/z", "Alignment")),
	ModifyCohesion:     key.NewBinding(key.WithKeys("s", "x"), key.WithHelp("s/x", "Cohesion")),
	ModifySeparation:   key.NewBinding(key.WithKeys("d", "c"), key.WithHelp("d/c", "Separation")),
	ModifyPerception:   key.NewBinding(key.WithKeys("f", "v"), key.WithHelp("f/v", "Perception")),
	ModifyMaxSpeed:     key.NewBinding(key.WithKeys("g", "b"), key.WithHelp("g/b", "Max Speed")),
	IncreaseAlignment:  key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "Inc Alignment")),
	DecreaseAlignment:  key.NewBinding(key.WithKeys("z"), key.WithHelp("z", "Dec Alignment")),
	IncreaseCohesion:   key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "Inc Cohesion")),
	DecreaseCohesion:   key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "Dec Cohesion")),
	IncreaseSeparation: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "Inc Separation")),
	DecreaseSeparation: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "Dec Separation")),
	IncreasePerception: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "Inc Perception")),
	DecreasePerception: key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "Dec Perception")),
	IncreaseMaxSpeed:   key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "Inc Max Speed")),
	DecreaseMaxSpeed:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "Dec Max Speed")),
	Scatter:            key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "Scatter")),
	ToggleEdgeMode:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Edge Mode")),
	Quit:               key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "Quit")),
}
