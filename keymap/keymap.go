package keymap

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	MoveTaskNext key.Binding
	AddTask      key.Binding
	DeleteTask   key.Binding
	Quit         key.Binding
	Help         key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.MoveTaskNext, k.AddTask, k.DeleteTask, k.Help, k.Quit}, // first column

	}
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),           // actual keybindings
		key.WithHelp("↑", "move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "move right"),
	),
	MoveTaskNext: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "move task to the next stage"),
	),
	AddTask: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "add task"),
	),
	DeleteTask: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}
