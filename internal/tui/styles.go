package tui

import "github.com/charmbracelet/lipgloss"

// ─── Color Palette ──────────────────────────────────────────────────────────

var (
	// Core brand colors
	ColorPrimary = lipgloss.Color("#7c3aed") // violet-600
	ColorAccent  = lipgloss.Color("#a78bfa") // violet-400
	ColorMagenta = lipgloss.Color("#c084fc") // purple-400
	ColorCyan    = lipgloss.Color("#22d3ee") // cyan-400
	ColorGreen   = lipgloss.Color("#34d399") // emerald-400
	ColorYellow  = lipgloss.Color("#fbbf24") // amber-400
	ColorRed     = lipgloss.Color("#f87171") // red-400
	ColorOrange  = lipgloss.Color("#fb923c") // orange-400
	ColorPink    = lipgloss.Color("#f472b6") // pink-400
	ColorIndigo  = lipgloss.Color("#818cf8") // indigo-400
	ColorTeal    = lipgloss.Color("#2dd4bf") // teal-400

	// Neutral tones
	ColorDimmed    = lipgloss.Color("#6b7280") // gray-500
	ColorSubtle    = lipgloss.Color("#374151") // gray-700
	ColorMuted     = lipgloss.Color("#4b5563") // gray-600
	ColorBorder    = lipgloss.Color("#4c1d95") // violet-900
	ColorBorderLit = lipgloss.Color("#6d28d9") // violet-700

	// Backgrounds
	ColorBg        = lipgloss.Color("#09090b") // zinc-950
	ColorBgPanel   = lipgloss.Color("#0f0a1e") // deep violet-black
	ColorBgSurface = lipgloss.Color("#18122b") // slightly lighter
	ColorHighlight = lipgloss.Color("#1e1045") // selected row bg
	ColorHighAlt   = lipgloss.Color("#2a1a5e") // alternate highlight

	// Foregrounds
	ColorFg       = lipgloss.Color("#e5e7eb") // gray-200
	ColorFgBright = lipgloss.Color("#f9fafb") // gray-50
	ColorFgDim    = lipgloss.Color("#9ca3af") // gray-400
)

// ─── Layout ─────────────────────────────────────────────────────────────────

var (
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(1, 2)

	SidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtle).
			Padding(1, 1)
)

// ─── Header / Status ────────────────────────────────────────────────────────

var (
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFgBright).
			Background(ColorPrimary).
			Padding(0, 2).
			MarginBottom(1)

	SubHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginBottom(1)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Background(lipgloss.Color("#110a28")).
			Padding(0, 1)

	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorMagenta)

	VersionStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed).
			Italic(true)
)

// ─── Text Styles ────────────────────────────────────────────────────────────

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFgBright)

	ArtistStyle = lipgloss.NewStyle().
			Foreground(ColorAccent)

	AlbumStyle = lipgloss.NewStyle().
			Foreground(ColorCyan)

	GenreStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	DurationStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ColorFgBright).
			Bold(true)

	CountBadgeStyle = lipgloss.NewStyle().
			Foreground(ColorBg).
			Background(ColorAccent).
			Padding(0, 1).
			Bold(true)

	FavBadgeStyle = lipgloss.NewStyle().
			Foreground(ColorPink).
			Bold(true)
)

// ─── Stat Cards ─────────────────────────────────────────────────────────────

var (
	StatCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtle).
			Padding(1, 3).
			Width(20).
			Align(lipgloss.Center)

	StatValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFgBright).
			Align(lipgloss.Center)

	StatLabelStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed).
			Align(lipgloss.Center)
)

// ─── Tabs ───────────────────────────────────────────────────────────────────

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorFgBright).
			Background(ColorPrimary).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(ColorDimmed).
				Padding(0, 2)

	TabBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#110a28")).
			Padding(0, 0)
)

var (
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(ColorFgBright).
				Background(ColorHighlight).
				Bold(true).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Padding(0, 1)

	NormalItemAltStyle = lipgloss.NewStyle().
				Foreground(ColorFg).
				Background(lipgloss.Color("#0d0919")).
				Padding(0, 1)

	PlayingItemStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Bold(true).
				Padding(0, 1)

	SelectedPlayingStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Background(ColorHighlight).
				Bold(true).
				Padding(0, 1)
)

// ─── Metadata Detail ────────────────────────────────────────────────────────

var (
	MetaLabelStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed).
			Width(14).
			Align(lipgloss.Right).
			PaddingRight(1)

	MetaValueStyle = lipgloss.NewStyle().
			Foreground(ColorFg)
)

// ─── Help ───────────────────────────────────────────────────────────────────

var (
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed)

	HelpSepStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle)
)

// ─── Now Playing Bar & View ─────────────────────────────────────────────────

var (
	NowPlayingBarStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	ProgressFilledStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary)

	ProgressEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle)

	ProgressKnobStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	NowPlayingTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorFgBright)

	NowPlayingArtistStyle = lipgloss.NewStyle().
				Foreground(ColorAccent)

	NowPlayingAlbumStyle = lipgloss.NewStyle().
				Foreground(ColorCyan)

	NowPlayingTimeStyle = lipgloss.NewStyle().
				Foreground(ColorDimmed)

	PlayButtonStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	PauseButtonStyle = lipgloss.NewStyle().
				Foreground(ColorYellow).
				Bold(true)

	ControlStyle = lipgloss.NewStyle().
			Foreground(ColorDimmed)

	ActiveControlStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)
)

// ─── Queue View ─────────────────────────────────────────────────────────────

var (
	QueueHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorIndigo)

	QueueCurrentStyle = lipgloss.NewStyle().
				Foreground(ColorGreen).
				Bold(true)

	QueueItemStyle = lipgloss.NewStyle().
			Foreground(ColorFgDim)

	QueueIndexStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Width(4).
			Align(lipgloss.Right)
)

// ─── Favorites ──────────────────────────────────────────────────────────────

var (
	FavHeartStyle = lipgloss.NewStyle().
			Foreground(ColorPink).
			Bold(true)

	FavHeartEmptyStyle = lipgloss.NewStyle().
				Foreground(ColorSubtle)
)

// ─── Spinner ────────────────────────────────────────────────────────────────

var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
