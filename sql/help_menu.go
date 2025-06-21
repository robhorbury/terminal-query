package sql

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type helpMenu struct {
	keys     help.KeyMap
	help     help.Model
	fullHelp bool
}

func (h helpMenu) View() string {
	if h.fullHelp == true {
		return h.help.FullHelpView(h.keys.FullHelp())
	} else {
		return h.help.View(h.keys)
	}
}

func (h *helpMenu) ToggleFullHelp() {
	h.fullHelp = !h.fullHelp
}

type navigationKeyMap struct {
	Up              key.Binding
	Down            key.Binding
	Left            key.Binding
	Right           key.Binding
	Help            key.Binding
	Quit            key.Binding
	RegexFilter     key.Binding
	SubstringFilter key.Binding
	VisibleColumns  key.Binding
	FilterColumns   key.Binding
}

func (k navigationKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.VisibleColumns, k.SubstringFilter, k.Help}
}

func (k navigationKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.VisibleColumns, k.FilterColumns, k.SubstringFilter, k.RegexFilter},
		{k.Help, k.Quit},
	}
}

type filteringKeyMap struct {
	Apply key.Binding
	Exit  key.Binding
}

func (k filteringKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Apply, k.Exit}
}

func (k filteringKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Apply, k.Exit},
	}
}

var navigationKeys = navigationKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	RegexFilter: key.NewBinding(
		key.WithKeys("\\"),
		key.WithHelp("\\", "regex filter"),
	),
	SubstringFilter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "substring filter"),
	),
	VisibleColumns: key.NewBinding(
		key.WithKeys("."),
		key.WithHelp(".", "visible columns"),
	),
	FilterColumns: key.NewBinding(
		key.WithKeys(","),
		key.WithHelp(",", "filter columns"),
	),
}

var navigationHelpMenu = helpMenu{
	keys:     navigationKeys,
	help:     help.New(),
	fullHelp: false,
}

func newNavigationMenu() helpMenu {
	return navigationHelpMenu
}

var filteringKeys = filteringKeyMap{
	Apply: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "apply"),
	),
	Exit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit"),
	),
}

var filteringHelpMenu = helpMenu{
	keys:     filteringKeys,
	help:     help.New(),
	fullHelp: false,
}

func newFilteringMenu() helpMenu {
	return filteringHelpMenu
}
