package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit     key.Binding
	Help     key.Binding
	Search   key.Binding
	Escape   key.Binding
	Enter    key.Binding
	Back     key.Binding
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Home     key.Binding
	End      key.Binding
	Refresh  key.Binding

	// Playback
	PlayPause  key.Binding
	Stop       key.Binding
	NextTrack  key.Binding
	PrevTrack  key.Binding
	VolumeUp   key.Binding
	VolumeDown key.Binding
	Mute       key.Binding
	SeekFwd    key.Binding
	SeekBack   key.Binding
	ShuffleTog key.Binding
	RepeatTog  key.Binding
	PlayAll    key.Binding
	NowPlaying key.Binding

	// Extended
	Favorite key.Binding
	Queue    key.Binding
	Detail   key.Binding

	// Direct tabs
	Tab1 key.Binding
	Tab2 key.Binding
	Tab3 key.Binding
	Tab4 key.Binding
	Tab5 key.Binding
	Tab6 key.Binding
	Tab7 key.Binding
	Tab8 key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/close"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "play/select"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("bksp", "go back"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("s-tab", "prev tab"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdn", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g", "go to top"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G", "go to bottom"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("C-r", "rescan"),
		),

		// Playback
		PlayPause: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "play/pause"),
		),
		Stop: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop"),
		),
		NextTrack: key.NewBinding(
			key.WithKeys("n", ">"),
			key.WithHelp("n/>", "next track"),
		),
		PrevTrack: key.NewBinding(
			key.WithKeys("N", "<"),
			key.WithHelp("N/<", "prev track"),
		),
		VolumeUp: key.NewBinding(
			key.WithKeys("+", "="),
			key.WithHelp("+", "volume up"),
		),
		VolumeDown: key.NewBinding(
			key.WithKeys("-", "_"),
			key.WithHelp("-", "volume down"),
		),
		Mute: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "mute"),
		),
		SeekFwd: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "seek +5s"),
		),
		SeekBack: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "seek −5s"),
		),
		ShuffleTog: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "shuffle"),
		),
		RepeatTog: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "repeat mode"),
		),
		PlayAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "play all"),
		),
		NowPlaying: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "now playing"),
		),

		// Extended
		Favorite: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "♥ favorite"),
		),
		Queue: key.NewBinding(
			key.WithKeys("Q"),
			key.WithHelp("Q", "queue"),
		),
		Detail: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "track details"),
		),

		// Direct tab access
		Tab1: key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "dashboard")),
		Tab2: key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "artists")),
		Tab3: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "albums")),
		Tab4: key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "genres")),
		Tab5: key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "songs")),
		Tab6: key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "now playing")),
		Tab7: key.NewBinding(key.WithKeys("7"), key.WithHelp("7", "queue")),
		Tab8: key.NewBinding(key.WithKeys("8"), key.WithHelp("8", "favorites")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.PlayPause, k.NextTrack, k.Enter, k.Search, k.Tab, k.Favorite, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown, k.Home, k.End},
		{k.Enter, k.Escape, k.Back, k.Tab, k.ShiftTab},
		{k.PlayPause, k.Stop, k.NextTrack, k.PrevTrack},
		{k.VolumeUp, k.VolumeDown, k.Mute},
		{k.SeekFwd, k.SeekBack},
		{k.ShuffleTog, k.RepeatTog, k.PlayAll, k.NowPlaying},
		{k.Favorite, k.Queue, k.Detail},
		{k.Search, k.Refresh, k.Help, k.Quit},
		{k.Tab1, k.Tab2, k.Tab3, k.Tab4, k.Tab5, k.Tab6, k.Tab7, k.Tab8},
	}
}
