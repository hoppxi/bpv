package tui

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hoppxi/bpv/internal/cache"
	"github.com/hoppxi/bpv/internal/metadata"
)

type listEntry struct {
	name  string
	count int
}

func mapToSortedEntries(data map[string]int) []listEntry {
	entries := make([]listEntry, 0, len(data))
	for k, v := range data {
		entries = append(entries, listEntry{k, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].name) < strings.ToLower(entries[j].name)
	})
	return entries
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "…"
}

func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	if total <= 0 {
		return "--:--"
	}
	m := total / 60
	s := total % 60
	if m >= 60 {
		h := m / 60
		m = m % 60
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

func formatBytes(b int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func metaRow(label, value string) string {
	return MetaLabelStyle.Render(label+":") + " " + MetaValueStyle.Render(value)
}

// ─── Tab Bar ────────────────────────────────────────────────────────────────

func renderTabBar(tabs []string, active int, width int) string {
	var rendered []string
	for i, tab := range tabs {
		if i == active {
			rendered = append(rendered, ActiveTabStyle.Render(" "+tab+" "))
		} else {
			rendered = append(rendered, InactiveTabStyle.Render(" "+tab+" "))
		}
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	return TabBarStyle.Width(width).Render(bar)
}

func renderDashboard(lib *cache.CachedLibrary, width int) string {
	logo := LogoStyle.Render("♫  BPV Music Player  ♫")
	subtitle := DimStyle.Render("Music Library · Browser & Player")
	version := VersionStyle.Render("v0.2.0")
	header := lipgloss.JoinVertical(lipgloss.Center, logo, subtitle, version)
	header = lipgloss.PlaceHorizontal(width-4, lipgloss.Center, header)

	cards := []string{
		renderStatCard("♪ Tracks", fmt.Sprintf("%d", lib.FileCount), ColorCyan),
		renderStatCard("♫ Artists", fmt.Sprintf("%d", len(lib.Artists)), ColorAccent),
		renderStatCard("◉ Albums", fmt.Sprintf("%d", len(lib.Albums)), ColorGreen),
		renderStatCard("★ Genres", fmt.Sprintf("%d", len(lib.Genres)), ColorYellow),
	}
	cardRow := lipgloss.JoinHorizontal(lipgloss.Top, cards...)
	cardRow = lipgloss.PlaceHorizontal(width-4, lipgloss.Center, cardRow)

	var nowPlayingHero string
	topArtists := renderTopList("Top Artists", lib.Artists, ArtistStyle, 6)
	topAlbums := renderTopList("Top Albums", lib.Albums, AlbumStyle, 6)
	topGenres := renderTopList("Top Genres", lib.Genres, GenreStyle, 6)

	colWidth := (width - 10) / 3
	if colWidth < 20 {
		colWidth = 20
	}

	col1 := lipgloss.NewStyle().Width(colWidth).Render(topArtists)
	col2 := lipgloss.NewStyle().Width(colWidth).Render(topAlbums)
	col3 := lipgloss.NewStyle().Width(colWidth).Render(topGenres)
	columns := lipgloss.JoinHorizontal(lipgloss.Top, col1, col2, col3)

	scanInfo := DimStyle.Render(fmt.Sprintf(
		"Music: %s  •  Scanned: %s  •  Press 'a' to play all, '↵' on a track to play",
		lib.Dir,
		lib.ScanTime.Format("15:04"),
	))

	parts := []string{"", header, "", cardRow}
	if nowPlayingHero != "" {
		parts = append(parts, "", nowPlayingHero)
	}
	parts = append(parts, "", columns, "", scanInfo)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func renderStatCard(label, value string, color lipgloss.Color) string {
	card := StatCardStyle.BorderForeground(color)
	v := StatValueStyle.Foreground(color).Render(value)
	l := StatLabelStyle.Render(label)
	return card.Render(lipgloss.JoinVertical(lipgloss.Center, v, l))
}

func renderTopList(title string, data map[string]int, style lipgloss.Style, max int) string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range data {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Value > sorted[j].Value })
	if len(sorted) > max {
		sorted = sorted[:max]
	}

	lines := []string{SubHeaderStyle.Render(title)}
	for i, item := range sorted {
		num := DimStyle.Render(fmt.Sprintf("%2d. ", i+1))
		name := style.Render(truncate(item.Key, 22))
		count := DimStyle.Render(fmt.Sprintf(" (%d)", item.Value))
		lines = append(lines, num+name+count)
	}
	if len(sorted) == 0 {
		lines = append(lines, DimStyle.Render("  No data yet"))
	}
	return strings.Join(lines, "\n")
}

func renderCategoryList(items []listEntry, cursor int, title string, width, height int) string {
	header := SubHeaderStyle.Render(fmt.Sprintf("%s  (%d)", title, len(items)))
	scrollInfo := DimStyle.Render(fmt.Sprintf("  %d/%d  •  ↵ browse  •  'a' play all  •  'd' details", cursor+1, len(items)))

	// Calculate visible space specifically for list items
	// Height - header (2 lines roughly) - spacer (1 line) - footer (1 line) - bottom spacer (1 line) = 5 lines overhead
	overhead := 5
	visibleLines := height - overhead
	if visibleLines < 1 {
		visibleLines = 1
	}

	start, end := scrollWindow(cursor, len(items), visibleLines)

	var lines []string
	for i := start; i < end; i++ {
		entry := items[i]
		prefix := "  "
		style := NormalItemStyle
		if i%2 == 1 {
			style = NormalItemAltStyle
		}
		if i == cursor {
			prefix = "▸ "
			style = SelectedItemStyle
		}

		badge := CountBadgeStyle.Render(fmt.Sprintf("%d", entry.count))
		name := style.Render(truncate(entry.name, width-20))
		lines = append(lines, fmt.Sprintf("%s%s  %s", prefix, name, badge))
	}

	// Fill empty lines to maintain stable layout
	for len(lines) < visibleLines {
		lines = append(lines, "")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "", strings.Join(lines, "\n"), "", scrollInfo,
	)
}

// ─── Song List ──────────────────────────────────────────────────────────────

func renderSongList(songs []metadata.AudioFile, cursor int, title string, width, height int, currentTrackPath string, player *Player) string {
	header := SubHeaderStyle.Render(fmt.Sprintf("%s  (%d tracks)", title, len(songs)))

	numW := 5
	titleW := (width - 30) / 3
	if titleW < 15 {
		titleW = 15
	}
	artistW := titleW
	albumW := titleW
	durW := 8

	tableHeader := DimStyle.Render(fmt.Sprintf(
		"%-*s %-*s %-*s %-*s %*s",
		numW, "#",
		titleW, "TITLE",
		artistW, "ARTIST",
		albumW, "ALBUM",
		durW, "DUR",
	))
	sep := DimStyle.Render(strings.Repeat("─", min(width-6, 120)))

	overhead := 6
	if title == "All Tracks" {
		overhead = 7
	}
	visibleLines := height - overhead
	if visibleLines < 1 {
		visibleLines = 1
	}

	start, end := scrollWindow(cursor, len(songs), visibleLines)

	var lines []string
	for i := start; i < end; i++ {
		s := songs[i]
		selected := i == cursor
		isPlaying := currentTrackPath != "" && s.FilePath == currentTrackPath
		isFav := player != nil && player.IsFavorite(s.FilePath)

		num := fmt.Sprintf("%-*d", numW, i+1)
		tTitle := truncate(s.Title, titleW-1)
		artist := truncate(s.Artist, artistW-1)
		album := truncate(s.Album, albumW-1)
		dur := formatDuration(s.Duration)

		row := fmt.Sprintf("%-*s %-*s %-*s %-*s %*s",
			numW, num,
			titleW, tTitle,
			artistW, artist,
			albumW, album,
			durW, dur,
		)

		heart := " "
		if isFav {
			heart = FavBadgeStyle.Render("♥")
		}

		switch {
		case selected && isPlaying:
			row = SelectedPlayingStyle.Render("♫ " + row + " " + heart)
		case selected:
			row = SelectedItemStyle.Render("▸ " + row + " " + heart)
		case isPlaying:
			row = PlayingItemStyle.Render("♫ " + row + " " + heart)
		default:
			base := NormalItemStyle
			if i%2 == 1 {
				base = NormalItemAltStyle
			}
			row = base.Render("  " + row + " " + heart)
		}
		lines = append(lines, row)
	}

	for len(lines) < visibleLines {
		lines = append(lines, "")
	}

	scrollInfo := DimStyle.Render(fmt.Sprintf("  %d/%d  •  ↵ play  •  'a' play all  •  'f' ♥  •  'd' details", cursor+1, len(songs)))

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "", tableHeader, sep, strings.Join(lines, "\n"), "", scrollInfo,
	)
}

func renderTrackDetail(track metadata.AudioFile, player *Player, width, height int) string {
	isFav := player != nil && player.IsFavorite(track.FilePath)
	favIcon := FavHeartEmptyStyle.Render("♡")
	if isFav {
		favIcon = FavHeartStyle.Render("♥")
	}

	title := TitleStyle.Render("♪ " + track.Title + " " + favIcon)

	rows := []string{
		metaRow("Artist", ArtistStyle.Render(track.Artist)),
		metaRow("Album", AlbumStyle.Render(track.Album)),
		metaRow("Genre", GenreStyle.Render(track.Genre)),
	}

	if track.AlbumArtist != "" && track.AlbumArtist != track.Artist {
		rows = append(rows, metaRow("Album Artist", track.AlbumArtist))
	}
	if track.Composer != "" && track.Composer != "Unknown Composer" {
		rows = append(rows, metaRow("Composer", track.Composer))
	}
	if track.Year > 0 {
		rows = append(rows, metaRow("Year", fmt.Sprintf("%d", track.Year)))
	}
	if track.Track > 0 {
		t := fmt.Sprintf("%d", track.Track)
		if track.TotalTracks > 0 {
			t += fmt.Sprintf(" / %d", track.TotalTracks)
		}
		rows = append(rows, metaRow("Track", t))
	}
	if track.Disc > 0 {
		d := fmt.Sprintf("%d", track.Disc)
		if track.TotalDiscs > 0 {
			d += fmt.Sprintf(" / %d", track.TotalDiscs)
		}
		rows = append(rows, metaRow("Disc", d))
	}

	rows = append(rows, "")
	rows = append(rows, metaRow("Duration", DurationStyle.Render(formatDuration(track.Duration))))
	rows = append(rows, metaRow("Format", HighlightStyle.Render(strings.ToUpper(track.FileType))))
	if track.Bitrate > 0 {
		rows = append(rows, metaRow("Bitrate", fmt.Sprintf("%d kbps", track.Bitrate)))
	}
	if track.SampleRate > 0 {
		rows = append(rows, metaRow("Sample Rate", fmt.Sprintf("%d Hz", track.SampleRate)))
	}
	if track.Channels > 0 {
		ch := "Mono"
		if track.Channels == 2 {
			ch = "Stereo"
		} else if track.Channels > 2 {
			ch = fmt.Sprintf("%d channels", track.Channels)
		}
		rows = append(rows, metaRow("Channels", ch))
	}

	rows = append(rows, "")
	rows = append(rows, metaRow("File", DimStyle.Render(truncate(track.FileName, width-20))))
	rows = append(rows, metaRow("Size", DimStyle.Render(formatBytes(track.FileSize))))

	if track.Comment != "" {
		rows = append(rows, "")
		rows = append(rows, metaRow("Comment", truncate(track.Comment, width-20)))
	}

	hint := DimStyle.Render("  Press 'f' to toggle ♥  •  ↵ to play  •  esc to go back")

	return lipgloss.JoinVertical(lipgloss.Left,
		title, "", strings.Join(rows, "\n"), "", hint,
	)
}

// ─── Now Playing View ───────────────────────────────────────────────────────

func renderNowPlaying(player *Player, width, height int) string {
	track := player.CurrentTrack()
	if track == nil {
		return lipgloss.Place(width-4, height,
			lipgloss.Center, lipgloss.Center,
			DimStyle.Render("Nothing playing. Select a track and press ↵ to play."),
		)
	}

	// Cover art dimensions.
	artWidth := 36
	artHeight := 16
	if width < 80 {
		artWidth = width / 3
		artHeight = artWidth / 2
	}
	if height < 25 {
		artHeight = height / 4
		artWidth = artHeight * 2
	}

	coverArt := renderCoverArt(track.CoverArt, artWidth, artHeight)

	// Fav indicator.
	favIcon := FavHeartEmptyStyle.Render("♡")
	if player.IsFavorite(track.FilePath) {
		favIcon = FavHeartStyle.Render("♥")
	}

	titleLine := lipgloss.NewStyle().Bold(true).Foreground(ColorFgBright).Align(lipgloss.Center).Width(width - 8).
		Render(track.Title + " " + favIcon)
	artistLine := lipgloss.NewStyle().Foreground(ColorAccent).Align(lipgloss.Center).Width(width - 8).
		Render(track.Artist)
	albumLine := lipgloss.NewStyle().Foreground(ColorCyan).Align(lipgloss.Center).Width(width - 8).
		Render(track.Album)

	pos := player.Position()
	dur := player.Duration()
	progress := player.Progress()

	barWidth := width - 30
	if barWidth < 20 {
		barWidth = 20
	}
	progressBar := renderProgressBar(progress, barWidth)

	timeLeft := NowPlayingTimeStyle.Render(formatDuration(pos))
	timeRight := NowPlayingTimeStyle.Render(formatDuration(dur))
	timeLine := fmt.Sprintf("  %s  %s  %s", timeLeft, progressBar, timeRight)

	// Controls.
	playIcon := "▶"
	playStyle := PlayButtonStyle
	if player.IsPlaying() {
		playIcon = "▮▮"
		playStyle = PauseButtonStyle
	}

	shuffleStyle := ControlStyle
	if player.Shuffle() {
		shuffleStyle = ActiveControlStyle
	}

	repeatIcon := player.Repeat().Icon()
	repeatStyle := ControlStyle
	if player.Repeat() != RepeatOff {
		repeatStyle = ActiveControlStyle
	}

	volIcon := "🔊"
	if player.IsMuted() {
		volIcon = "🔇"
	}

	controls := fmt.Sprintf(
		"  %s  %s  %s  %s  %s    %s %d%%    %s",
		shuffleStyle.Render("⇄"),
		ControlStyle.Render("⏮"),
		playStyle.Render(playIcon),
		ControlStyle.Render("⏭"),
		repeatStyle.Render(repeatIcon),
		ControlStyle.Render(volIcon),
		player.VolumePercent(),
		DimStyle.Render(fmt.Sprintf("%d/%d tracks", player.QueueIndex()+1, player.QueueLen())),
	)

	formatInfo := DimStyle.Render(fmt.Sprintf(
		"%s · %d kbps · %d Hz",
		strings.ToUpper(track.FileType),
		track.Bitrate,
		track.SampleRate,
	))

	hint := DimStyle.Render("  'f' ♥ favorite  •  'Q' queue  •  space play/pause  •  ←→ seek")

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		coverArt,
		"",
		titleLine,
		artistLine,
		albumLine,
		"",
		timeLine,
		"",
		controls,
		"",
		formatInfo,
		"",
		hint,
	)

	return lipgloss.Place(width-4, height, lipgloss.Center, lipgloss.Center, content)
}

// ─── Now Playing Bar (bottom bar) ───────────────────────────────────────────

func renderNowPlayingBar(player *Player, width int) string {
	track := player.CurrentTrack()
	if track == nil {
		return ""
	}

	title := NowPlayingTitleStyle.Render(truncate(track.Title, 30))
	artist := NowPlayingArtistStyle.Render(truncate(track.Artist, 20))
	trackInfo := title + "  " + artist

	playIcon := "▶"
	playStyle := PlayButtonStyle
	if player.IsPlaying() {
		playIcon = "▮▮"
		playStyle = PauseButtonStyle
	}

	pos := player.Position()
	dur := player.Duration()
	progress := player.Progress()

	barWidth := width - 80
	if barWidth < 10 {
		barWidth = 10
	}
	progressBar := renderProgressBar(progress, barWidth)

	centerPart := fmt.Sprintf(
		"%s %s %s %s %s %s",
		ControlStyle.Render("⏮"),
		playStyle.Render(playIcon),
		ControlStyle.Render("⏭"),
		NowPlayingTimeStyle.Render(formatDuration(pos)),
		progressBar,
		NowPlayingTimeStyle.Render(formatDuration(dur)),
	)

	shuffleIcon := ControlStyle.Render("⇄")
	if player.Shuffle() {
		shuffleIcon = ActiveControlStyle.Render("⇄")
	}

	repeatIcon := ControlStyle.Render(player.Repeat().Icon())
	if player.Repeat() != RepeatOff {
		repeatIcon = ActiveControlStyle.Render(player.Repeat().Icon())
	}

	rightPart := fmt.Sprintf("%s %s", shuffleIcon, repeatIcon)

	leftW := lipgloss.Width(trackInfo)
	centerW := lipgloss.Width(centerPart)
	rightW := lipgloss.Width(rightPart)
	gap1 := (width - leftW - centerW - rightW) / 2
	if gap1 < 1 {
		gap1 = 1
	}
	gap2 := width - leftW - centerW - rightW - gap1
	if gap2 < 0 {
		gap2 = 0
	}

	bar := trackInfo + strings.Repeat(" ", gap1) + centerPart + strings.Repeat(" ", gap2) + rightPart

	return NowPlayingBarStyle.Width(width - 2).Render(bar)
}

// ─── Progress Bar ───────────────────────────────────────────────────────────

func renderProgressBar(progress float64, width int) string {
	filled := int(float64(width) * progress)
	if filled > width {
		filled = width
	}
	empty := width - filled

	knob := ""
	if filled > 0 && filled < width {
		filled--
		knob = ProgressKnobStyle.Render("●")
	}

	bar := ProgressFilledStyle.Render(strings.Repeat("━", filled)) +
		knob +
		ProgressEmptyStyle.Render(strings.Repeat("─", empty))

	return bar
}

// ─── Queue View ─────────────────────────────────────────────────────────────

func renderQueue(player *Player, cursor int, width, height int) string {
	queue := player.Queue()
	currentIdx := player.QueueIndex()

	if len(queue) == 0 {
		return lipgloss.Place(width-4, height,
			lipgloss.Center, lipgloss.Center,
			DimStyle.Render("Queue is empty. Play a track to populate the queue."),
		)
	}

	header := QueueHeaderStyle.Render(fmt.Sprintf("  Queue  (%d tracks)", len(queue)))

	// Overhead: header(1) + spacer(1) + footerSpacer(1) + footer(1) = 4 lines
	overhead := 4
	visibleLines := height - overhead
	if visibleLines < 1 {
		visibleLines = 1
	}

	start, end := scrollWindow(cursor, len(queue), visibleLines)

	var lines []string
	for i := start; i < end; i++ {
		track := queue[i]
		isCurrent := i == currentIdx
		isSelected := i == cursor

		idx := QueueIndexStyle.Render(fmt.Sprintf("%d.", i+1))
		title := truncate(track.Title, (width-30)/2)
		artist := truncate(track.Artist, (width-30)/2)

		var row string
		switch {
		case isCurrent && isSelected:
			row = SelectedPlayingStyle.Render(fmt.Sprintf("♫ %s %s — %s", idx, title, artist))
		case isCurrent:
			row = QueueCurrentStyle.Render(fmt.Sprintf("♫ %s %s — %s", idx, title, artist))
		case isSelected:
			row = SelectedItemStyle.Render(fmt.Sprintf("  %s %s — %s", idx, title, artist))
		default:
			row = QueueItemStyle.Render(fmt.Sprintf("  %s %s — %s", idx, title, artist))
		}
		lines = append(lines, row)
	}

	for len(lines) < visibleLines {
		lines = append(lines, "")
	}

	scrollInfo := DimStyle.Render(fmt.Sprintf("  %d/%d  •  ↵ play  •  esc back", cursor+1, len(queue)))

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "", strings.Join(lines, "\n"), "", scrollInfo,
	)
}

// ─── Favorites View ─────────────────────────────────────────────────────────

func renderFavorites(favTracks []metadata.AudioFile, cursor int, width, height int, currentTrackPath string) string {
	if len(favTracks) == 0 {
		return lipgloss.Place(width-4, height,
			lipgloss.Center, lipgloss.Center,
			DimStyle.Render("No favorites yet. Press 'f' on a track to add it."),
		)
	}

	header := FavHeartStyle.Render(fmt.Sprintf("  ♥ Favorites  (%d tracks)", len(favTracks)))

	// Overhead: header(1) + spacer(1) + footerSpacer(1) + footer(1) = 4 lines
	overhead := 4
	visibleLines := height - overhead
	if visibleLines < 1 {
		visibleLines = 1
	}

	start, end := scrollWindow(cursor, len(favTracks), visibleLines)

	var lines []string
	for i := start; i < end; i++ {
		track := favTracks[i]
		isSelected := i == cursor
		isPlaying := currentTrackPath != "" && track.FilePath == currentTrackPath

		title := truncate(track.Title, (width-20)/2)
		artist := truncate(track.Artist, (width-20)/2)
		heart := FavBadgeStyle.Render("♥ ")

		row := fmt.Sprintf("%s — %s", title, artist)
		switch {
		case isSelected && isPlaying:
			lines = append(lines, SelectedPlayingStyle.Render("♫ "+heart+row))
		case isSelected:
			lines = append(lines, SelectedItemStyle.Render("▸ "+heart+row))
		case isPlaying:
			lines = append(lines, PlayingItemStyle.Render("♫ "+heart+row))
		default:
			lines = append(lines, NormalItemStyle.Render("  "+heart+row))
		}
	}

	for len(lines) < visibleLines {
		lines = append(lines, "")
	}

	scrollInfo := DimStyle.Render(fmt.Sprintf("  %d/%d  •  ↵ play  •  'f' unfavorite  •  esc back", cursor+1, len(favTracks)))

	return lipgloss.JoinVertical(lipgloss.Left,
		header, "", strings.Join(lines, "\n"), "", scrollInfo,
	)
}

// ─── Search Results ─────────────────────────────────────────────────────────

func renderSearchResults(results []metadata.AudioFile, cursor int, query string, width, height int, currentTrackPath string, player *Player) string {
	title := fmt.Sprintf("Search: %s  (%d found)",
		HighlightStyle.Render("\""+query+"\""),
		len(results),
	)

	if len(results) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			SubHeaderStyle.Render(title), "",
			DimStyle.Render("  No results found."),
		)
	}

	return renderSongList(results, cursor, title, width, height, currentTrackPath, player)
}

// ─── Help Overlay ───────────────────────────────────────────────────────────

func renderHelp(keys KeyMap, width, height int) string {
	title := HeaderStyle.Render("  ⌨ Keyboard Shortcuts  ")

	groups := keys.FullHelp()
	// Labels corresponding to groups
	groupLabels := []string{"Navigation", "Actions", "Playback", "Volume", "Seek", "Modes", "Library", "General", "Direct Tabs"}

	numCols := 2
	if width > 85 {
		numCols = 3
	}

	cols := make([][]string, numCols)

	// Distribute groups round-robin or chunks? Chunks usually look better (reading down).
	// Total groups = 9.
	// 3 cols: 3 groups each.
	// 2 cols: 5 and 4.

	perCol := (len(groups) + numCols - 1) / numCols

	for i, group := range groups {
		label := ""
		if i < len(groupLabels) {
			label = groupLabels[i]
		}

		var bindings []string
		for _, b := range group {
			k := HelpKeyStyle.Render(fmt.Sprintf("%-6s", b.Help().Key))
			d := HelpDescStyle.Render(b.Help().Desc)
			bindings = append(bindings, fmt.Sprintf("%s %s", k, d))
		}

		section := SubHeaderStyle.Render(label) + "\n" + strings.Join(bindings, "\n")

		colIdx := i / perCol
		if colIdx >= numCols {
			colIdx = numCols - 1
		}
		cols[colIdx] = append(cols[colIdx], section)
	}

	// Render columns
	var renderedCols []string
	for _, c := range cols {
		renderedCols = append(renderedCols, strings.Join(c, "\n\n"))
	}

	// Join columns horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top,
		renderedCols[0],
		strings.Repeat("   ", 2), // Gap
		renderedCols[1],
	)
	if numCols > 2 {
		content = lipgloss.JoinHorizontal(lipgloss.Top,
			content,
			strings.Repeat("   ", 2),
			renderedCols[2],
		)
	}

	boxWidth := width - 6
	if boxWidth > 120 {
		boxWidth = 120
	}

	// Calculate box height to center nicely, but clamp to screen height
	boxHeight := height - 4
	if boxHeight < 10 {
		boxHeight = 10
	}

	box := PanelStyle.Width(boxWidth).
		Render(
			lipgloss.JoinVertical(lipgloss.Center,
				title,
				"",
				content,
				"",
				DimStyle.Render("Press any key to close"),
			),
		)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

// ─── Scroll Window ──────────────────────────────────────────────────────────

func scrollWindow(cursor, total, visible int) (int, int) {
	if total <= visible {
		return 0, total
	}

	// Keep cursor roughly centered.
	half := visible / 2
	start := cursor - half
	if start < 0 {
		start = 0
	}
	end := start + visible
	if end > total {
		end = total
		start = end - visible
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Ensure math import is used.
var _ = math.Round
