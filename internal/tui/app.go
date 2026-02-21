package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hoppxi/bpv/internal/cache"
	"github.com/hoppxi/bpv/internal/daemon"
	"github.com/hoppxi/bpv/internal/metadata"
	"github.com/hoppxi/bpv/internal/store"
)

type viewKind int

const (
	viewDashboard viewKind = iota
	viewArtists
	viewAlbums
	viewGenres
	viewSongs
	viewTrackDetail
	viewSearch
	viewNowPlaying
	viewQueue
	viewFavorites
)

var tabNames = []string{"Dashboard", "Artists", "Albums", "Genres", "Songs", "Now Playing", "Queue", "Favorites"}

type libraryScanDone struct {
	lib *cache.CachedLibrary
	err error
}

type favoritesLoaded struct {
	favs []string
	err  error
}

type queueLoaded struct {
	queue *store.QueueState
	err   error
}

type tickMsg time.Time
type spinnerTick time.Time

type Model struct {
	client   *daemon.Client
	lib      *cache.CachedLibrary
	musicDir string
	keys     KeyMap
	width    int
	height   int
	player   *Player

	activeView   viewKind
	activeTab    int
	showHelp     bool
	scanning     bool
	searchActive bool
	searchInput  textinput.Model
	err          error

	artistCursor int
	albumCursor  int
	genreCursor  int
	songCursor   int
	searchCursor int
	queueCursor  int
	favCursor    int

	artistList []listEntry
	albumList  []listEntry
	genreList  []listEntry
	songList   []metadata.AudioFile
	searchRes  []metadata.AudioFile
	favTracks  []metadata.AudioFile

	allFiles []metadata.AudioFile

	detailTrack metadata.AudioFile
	filterLabel string

	viewStack []viewKind

	spinnerIdx int
}

func New(musicDir string) Model {
	ti := textinput.New()
	ti.Placeholder = "Search songs, artists, albums…"
	ti.CharLimit = 100
	ti.Width = 40

	p := NewPlayer()

	return Model{
		musicDir:    musicDir,
		keys:        DefaultKeyMap(),
		activeView:  viewDashboard,
		activeTab:   0,
		scanning:    true,
		searchInput: ti,
		viewStack:   []viewKind{},
		player:      p,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.connectAndLoad(),
		tickCmd(),
		spinnerCmd(),
	)
}

func (m Model) connectAndLoad() tea.Cmd {
	return func() tea.Msg {
		client, err := daemon.Connect()
		if err != nil {
			return libraryScanDone{err: fmt.Errorf("daemon: %w", err)}
		}

		lib, err := client.GetLibrary(m.musicDir)
		if err != nil {
			client.Close()
			return libraryScanDone{err: fmt.Errorf("library: %w", err)}
		}

		return daemonConnected{client: client, lib: lib}
	}
}

type daemonConnected struct {
	client *daemon.Client
	lib    *cache.CachedLibrary
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func spinnerCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTick(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case daemonConnected:
		m.scanning = false
		m.client = msg.client
		m.lib = msg.lib
		m.rebuildCaches()

		return m, tea.Batch(m.loadFavorites(), m.loadQueue())

	case libraryScanDone:
		m.scanning = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.lib = msg.lib
		m.rebuildCaches()
		return m, nil

	case favoritesLoaded:
		if msg.err == nil && msg.favs != nil {
			m.player.SetFavoritePaths(msg.favs)
			m.rebuildFavTracks()
		}
		return m, nil

	case queueLoaded:
		if msg.err == nil && msg.queue != nil && len(msg.queue.FilePaths) > 0 && len(m.allFiles) > 0 {
			byPath := make(map[string]metadata.AudioFile, len(m.allFiles))
			for _, f := range m.allFiles {
				byPath[f.FilePath] = f
			}
			tracks := make([]metadata.AudioFile, 0, len(msg.queue.FilePaths))
			for _, p := range msg.queue.FilePaths {
				if t, ok := byPath[p]; ok {
					tracks = append(tracks, t)
				}
			}
			if len(tracks) > 0 {
				start := msg.queue.CurrentIndex
				if start < 0 {
					start = 0
				}
				if start >= len(tracks) {
					start = len(tracks) - 1
				}
				m.player.SetQueue(tracks, start)
				m.player.SetRepeat(RepeatMode(msg.queue.Repeat))
				m.player.SetShuffle(msg.queue.Shuffle)
			}
		}
		return m, nil

	case tickMsg:

		if m.player.CheckTrackEnd() {
			m.persistQueue()
		}
		return m, tickCmd()

	case spinnerTick:
		if m.scanning {
			m.spinnerIdx = (m.spinnerIdx + 1) % len(SpinnerFrames)
			return m, spinnerCmd()
		}
		return m, nil

	case tea.KeyMsg:
		if m.searchActive {
			return m.updateSearch(msg)
		}
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		return m.updateNormal(msg)
	}

	return m, nil
}

func (m Model) loadFavorites() tea.Cmd {
	return func() tea.Msg {
		if m.client == nil {
			return favoritesLoaded{}
		}
		favs, err := m.client.GetFavorites()
		return favoritesLoaded{favs: favs, err: err}
	}
}

func (m Model) loadQueue() tea.Cmd {
	return func() tea.Msg {
		if m.client == nil {
			return queueLoaded{}
		}
		q, err := m.client.GetQueue()
		return queueLoaded{queue: q, err: err}
	}
}

func (m *Model) persistQueue() {
	if m.client == nil {
		return
	}
	q := store.QueueState{
		FilePaths:    []string{},
		CurrentIndex: 0,
		Shuffle:      m.player.Shuffle(),
		Repeat:       int(m.player.Repeat()),
	}
	tracks := m.player.Queue()
	if len(tracks) > 0 {
		q.FilePaths = make([]string, 0, len(tracks))
		for _, t := range tracks {
			q.FilePaths = append(q.FilePaths, t.FilePath)
		}
		if idx := m.player.CurrentQueueItemIndex(); idx >= 0 {
			q.CurrentIndex = idx
		}
	}

	go m.client.SaveQueue(&q)
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Quit):
		m.player.Close()
		if m.client != nil {
			m.client.Close()
		}
		return m, tea.Quit

	case matchKey(msg, m.keys.Help):
		m.showHelp = true

	case matchKey(msg, m.keys.Search):
		m.searchActive = true
		m.searchInput.Focus()
		m.searchRes = nil
		m.searchCursor = 0
		m.pushView(viewSearch)
		return m, textinput.Blink

	case matchKey(msg, m.keys.Tab1):
		m.activeTab = 0
		m.switchToTab(0)
	case matchKey(msg, m.keys.Tab2):
		m.activeTab = 1
		m.switchToTab(1)
	case matchKey(msg, m.keys.Tab3):
		m.activeTab = 2
		m.switchToTab(2)
	case matchKey(msg, m.keys.Tab4):
		m.activeTab = 3
		m.switchToTab(3)
	case matchKey(msg, m.keys.Tab5):
		m.activeTab = 4
		m.switchToTab(4)
	case matchKey(msg, m.keys.Tab6):
		m.activeTab = 5
		m.switchToTab(5)
	case matchKey(msg, m.keys.Tab7):
		m.activeTab = 6
		m.switchToTab(6)
	case matchKey(msg, m.keys.Tab8):
		m.activeTab = 7
		m.switchToTab(7)

	case matchKey(msg, m.keys.Tab):
		m.activeTab = (m.activeTab + 1) % len(tabNames)
		m.switchToTab(m.activeTab)

	case matchKey(msg, m.keys.ShiftTab):
		m.activeTab = (m.activeTab - 1 + len(tabNames)) % len(tabNames)
		m.switchToTab(m.activeTab)

	case matchKey(msg, m.keys.Escape), matchKey(msg, m.keys.Back):
		m.goBack()

	case matchKey(msg, m.keys.PlayPause):
		if m.player.HasTrack() {
			m.player.TogglePause()
		} else {
			m.playCurrentSong()
		}

	case matchKey(msg, m.keys.Stop):
		m.player.Stop()

	case matchKey(msg, m.keys.NextTrack):
		_ = m.player.Next()
		m.persistQueue()

	case matchKey(msg, m.keys.PrevTrack):
		_ = m.player.Previous()
		m.persistQueue()

	case matchKey(msg, m.keys.VolumeUp):
		m.player.VolumeUp()

	case matchKey(msg, m.keys.VolumeDown):
		m.player.VolumeDown()

	case matchKey(msg, m.keys.Mute):
		m.player.ToggleMute()

	case matchKey(msg, m.keys.ShuffleTog):
		m.player.ToggleShuffle()
		m.persistQueue()

	case matchKey(msg, m.keys.RepeatTog):
		m.player.CycleRepeat()
		m.persistQueue()

	case matchKey(msg, m.keys.NowPlaying):
		if m.player.HasTrack() {
			m.activeTab = 5
			m.pushView(viewNowPlaying)
		}

	case matchKey(msg, m.keys.PlayAll):
		m.playAllFromCursor()

	case matchKey(msg, m.keys.SeekFwd):
		if m.player.HasTrack() {
			m.player.SeekForward(5 * time.Second)
		}

	case matchKey(msg, m.keys.SeekBack):
		if m.player.HasTrack() {
			m.player.SeekBackward(5 * time.Second)
		}

	case matchKey(msg, m.keys.Favorite):
		m.handleFavorite()

	case matchKey(msg, m.keys.Queue):
		m.activeTab = 6
		m.pushView(viewQueue)

	case matchKey(msg, m.keys.Detail):
		m.showDetail()

	case matchKey(msg, m.keys.Home):
		m.setCursor(0)

	case matchKey(msg, m.keys.End):
		l := m.currentListLen()
		if l > 0 {
			m.setCursor(l - 1)
		}

	case matchKey(msg, m.keys.Up):
		m.moveCursor(-1)

	case matchKey(msg, m.keys.Down):
		m.moveCursor(1)

	case matchKey(msg, m.keys.PageUp):
		m.moveCursor(-10)

	case matchKey(msg, m.keys.PageDown):
		m.moveCursor(10)

	case matchKey(msg, m.keys.Enter):
		m.handleEnter()

	case matchKey(msg, m.keys.Refresh):
		if m.lib != nil {
			m.scanning = true
			return m, tea.Batch(
				func() tea.Msg {
					if m.client != nil {
						lib, err := m.client.Scan(m.musicDir)
						return libraryScanDone{lib: lib, err: err}
					}
					return libraryScanDone{err: fmt.Errorf("no daemon connection")}
				},
				spinnerCmd(),
			)
		}
	}

	return m, nil
}

func (m *Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Escape):
		m.searchActive = false
		m.searchInput.Blur()
		m.goBack()
		return m, nil

	case matchKey(msg, m.keys.Enter):
		if len(m.searchRes) > 0 && m.searchCursor < len(m.searchRes) {
			m.player.SetQueue(m.searchRes, m.searchCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
			m.searchActive = false
			m.searchInput.Blur()
		}
		return m, nil

	case msg.Type == tea.KeyUp:
		if m.searchCursor > 0 {
			m.searchCursor--
		}
		return m, nil

	case msg.Type == tea.KeyDown:
		if m.searchCursor < len(m.searchRes)-1 {
			m.searchCursor++
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	if m.searchInput.Value() != "" && len(m.allFiles) > 0 {
		m.searchRes = searchFiles(m.allFiles, m.searchInput.Value())
		if m.searchCursor >= len(m.searchRes) {
			m.searchCursor = max(0, len(m.searchRes)-1)
		}
	} else {
		m.searchRes = nil
		m.searchCursor = 0
	}

	return m, cmd
}

func searchFiles(files []metadata.AudioFile, query string) []metadata.AudioFile {
	q := strings.ToLower(query)
	var results []metadata.AudioFile
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Title), q) ||
			strings.Contains(strings.ToLower(f.Artist), q) ||
			strings.Contains(strings.ToLower(f.Album), q) ||
			strings.Contains(strings.ToLower(f.Genre), q) {
			results = append(results, f)
		}
	}
	return results
}

func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	if m.scanning {
		return m.renderLoading()
	}
	if m.err != nil {
		return m.renderError()
	}
	if m.lib == nil {
		return m.renderLoading()
	}
	if m.showHelp {
		return renderHelp(m.keys, m.width, m.height)
	}

	tabBar := renderTabBar(tabNames, m.activeTab, m.width)
	tabBarHeight := lipgloss.Height(tabBar)

	var npBar string
	npBarHeight := 0
	hasTrack := m.player.HasTrack()

	if hasTrack && m.activeView != viewNowPlaying {
		npBar = renderNowPlayingBar(m.player, m.width)
		npBarHeight = lipgloss.Height(npBar)
	}

	statusBarHeight := 1
	panelOverhead := 4
	totalChromeHeight := tabBarHeight + npBarHeight + statusBarHeight + panelOverhead
	innerContentHeight := m.height - totalChromeHeight
	if innerContentHeight < 5 {
		innerContentHeight = 5
	}

	currentPath := ""
	if ct := m.player.CurrentTrack(); ct != nil {
		currentPath = ct.FilePath
	}

	var content string
	switch m.activeView {
	case viewDashboard:
		content = renderDashboard(m.lib, m.width)
	case viewArtists:
		content = renderCategoryList(m.artistList, m.artistCursor, "Artists", m.width, innerContentHeight)
	case viewAlbums:
		content = renderCategoryList(m.albumList, m.albumCursor, "Albums", m.width, innerContentHeight)
	case viewGenres:
		content = renderCategoryList(m.genreList, m.genreCursor, "Genres", m.width, innerContentHeight)
	case viewSongs:
		content = renderSongList(m.songList, m.songCursor, m.filterLabel, m.width, innerContentHeight, currentPath, m.player)
	case viewTrackDetail:
		content = renderTrackDetail(m.detailTrack, m.player, m.width, innerContentHeight)
	case viewNowPlaying:

		content = renderNowPlaying(m.player, m.width, innerContentHeight)
	case viewQueue:
		content = renderQueue(m.player, m.queueCursor, m.width, innerContentHeight)
	case viewFavorites:
		content = renderFavorites(m.favTracks, m.favCursor, m.width, innerContentHeight, currentPath)
	case viewSearch:
		searchBar := m.searchInput.View()
		if len(m.searchRes) > 0 || m.searchInput.Value() != "" {
			content = searchBar + "\n\n" + renderSearchResults(m.searchRes, m.searchCursor, m.searchInput.Value(), m.width, innerContentHeight-3, currentPath, m.player)
		} else {
			content = searchBar + "\n\n" + DimStyle.Render("  Type to search your library…")
		}
	}

	panel := PanelStyle.
		Width(m.width - 2).
		Height(innerContentHeight).
		Render(content)

	statusLeft := m.statusLine()
	statusRight := DimStyle.Render("? help  space play  / search  q quit")

	wLeft := lipgloss.Width(statusLeft)
	wRight := lipgloss.Width(statusRight)
	gap := m.width - wLeft - wRight - 2
	if gap < 0 {
		gap = 0
	}

	statusBar := StatusBarStyle.Width(m.width).Render(
		statusLeft + strings.Repeat(" ", gap) + statusRight,
	)

	parts := []string{tabBar, panel}
	if npBar != "" {
		parts = append(parts, npBar)
	}
	parts = append(parts, statusBar)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *Model) playCurrentSong() {
	if m.activeView == viewSongs && m.songCursor < len(m.songList) {
		m.player.SetQueue(m.songList, m.songCursor)
		_ = m.player.PlayCurrent()
		m.persistQueue()
	}
}

func (m *Model) playAllFromCursor() {
	switch m.activeView {
	case viewSongs:
		if len(m.songList) > 0 {
			m.player.SetQueue(m.songList, m.songCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
		}
	case viewArtists:
		if m.artistCursor < len(m.artistList) {
			entry := m.artistList[m.artistCursor]
			tracks := m.artistTracks(entry.name)
			if len(tracks) > 0 {
				m.player.SetQueue(tracks, 0)
				_ = m.player.PlayCurrent()
				m.persistQueue()
			}
		}
	case viewAlbums:
		if m.albumCursor < len(m.albumList) {
			entry := m.albumList[m.albumCursor]
			tracks := m.albumTracks(entry.name)
			if len(tracks) > 0 {
				m.player.SetQueue(tracks, 0)
				_ = m.player.PlayCurrent()
				m.persistQueue()
			}
		}
	case viewGenres:
		if m.genreCursor < len(m.genreList) {
			entry := m.genreList[m.genreCursor]
			tracks := m.genreTracks(entry.name)
			if len(tracks) > 0 {
				m.player.SetQueue(tracks, 0)
				_ = m.player.PlayCurrent()
				m.persistQueue()
			}
		}
	case viewDashboard:
		if len(m.allFiles) > 0 {
			m.player.SetQueue(m.allFiles, 0)
			_ = m.player.PlayCurrent()
			m.persistQueue()
		}
	case viewFavorites:
		if len(m.favTracks) > 0 {
			m.player.SetQueue(m.favTracks, m.favCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
		}
	}
}

func (m *Model) handleFavorite() {
	var filePath string

	switch m.activeView {
	case viewSongs:
		if m.songCursor < len(m.songList) {
			filePath = m.songList[m.songCursor].FilePath
		}
	case viewSearch:
		if m.searchCursor < len(m.searchRes) {
			filePath = m.searchRes[m.searchCursor].FilePath
		}
	case viewNowPlaying:
		if track := m.player.CurrentTrack(); track != nil {
			filePath = track.FilePath
		}
	case viewTrackDetail:
		filePath = m.detailTrack.FilePath
	case viewFavorites:
		if m.favCursor < len(m.favTracks) {
			filePath = m.favTracks[m.favCursor].FilePath
		}
	case viewQueue:
		queue := m.player.Queue()
		if m.queueCursor < len(queue) {
			filePath = queue[m.queueCursor].FilePath
		}
	}

	if filePath == "" {
		return
	}

	m.player.ToggleFavorite(filePath)

	if m.client != nil {
		if m.player.IsFavorite(filePath) {
			go m.client.AddFavorite(filePath)
		} else {
			go m.client.RemoveFavorite(filePath)
		}
	}

	if m.activeView == viewFavorites {
		m.rebuildFavTracks()
	}
}

func (m *Model) showDetail() {
	switch m.activeView {
	case viewSongs:
		if m.songCursor < len(m.songList) {
			m.detailTrack = m.songList[m.songCursor]
			m.pushView(viewTrackDetail)
		}
	case viewSearch:
		if m.searchCursor < len(m.searchRes) {
			m.detailTrack = m.searchRes[m.searchCursor]
			m.pushView(viewTrackDetail)
		}
	case viewNowPlaying:
		if track := m.player.CurrentTrack(); track != nil {
			m.detailTrack = *track
			m.pushView(viewTrackDetail)
		}
	case viewQueue:
		queue := m.player.Queue()
		if m.queueCursor < len(queue) {
			m.detailTrack = queue[m.queueCursor]
			m.pushView(viewTrackDetail)
		}
	}
}

func (m *Model) rebuildCaches() {
	if m.lib == nil {
		return
	}

	m.allFiles = m.lib.Files
	m.artistList = mapToSortedEntries(m.lib.Artists)
	m.albumList = mapToSortedEntries(m.lib.Albums)
	m.genreList = mapToSortedEntries(m.lib.Genres)
	m.songList = m.allFiles
	m.filterLabel = "All Tracks"
	m.rebuildFavTracks()
}

func (m *Model) rebuildFavTracks() {
	favPaths := m.player.FavoritePaths()
	if len(favPaths) == 0 {
		m.favTracks = nil
		return
	}
	pathSet := make(map[string]bool, len(favPaths))
	for _, p := range favPaths {
		pathSet[p] = true
	}
	m.favTracks = nil
	for _, f := range m.allFiles {
		if pathSet[f.FilePath] {
			m.favTracks = append(m.favTracks, f)
		}
	}
	if m.favCursor >= len(m.favTracks) {
		m.favCursor = max(0, len(m.favTracks)-1)
	}
}

func (m *Model) artistTracks(artist string) []metadata.AudioFile {
	var tracks []metadata.AudioFile
	for _, f := range m.allFiles {
		if f.Artist == artist {
			tracks = append(tracks, f)
		}
	}
	return tracks
}

func (m *Model) albumTracks(album string) []metadata.AudioFile {
	var tracks []metadata.AudioFile
	for _, f := range m.allFiles {
		if f.Album == album {
			tracks = append(tracks, f)
		}
	}
	return tracks
}

func (m *Model) genreTracks(genre string) []metadata.AudioFile {
	var tracks []metadata.AudioFile
	for _, f := range m.allFiles {
		if f.Genre == genre {
			tracks = append(tracks, f)
		}
	}
	return tracks
}

func (m *Model) switchToTab(tab int) {
	m.viewStack = nil
	switch tab {
	case 0:
		m.activeView = viewDashboard
	case 1:
		m.activeView = viewArtists
	case 2:
		m.activeView = viewAlbums
	case 3:
		m.activeView = viewGenres
	case 4:
		m.activeView = viewSongs
		m.songList = m.allFiles
		m.filterLabel = "All Tracks"
		m.songCursor = 0
	case 5:
		m.activeView = viewNowPlaying
	case 6:
		m.activeView = viewQueue
		m.queueCursor = 0
	case 7:
		m.activeView = viewFavorites
		m.favCursor = 0
	}
}

func (m *Model) pushView(v viewKind) {
	m.viewStack = append(m.viewStack, m.activeView)
	m.activeView = v
}

func (m *Model) goBack() {
	if len(m.viewStack) > 0 {
		m.activeView = m.viewStack[len(m.viewStack)-1]
		m.viewStack = m.viewStack[:len(m.viewStack)-1]
	}
}

func (m *Model) handleEnter() {
	switch m.activeView {
	case viewArtists:
		if m.artistCursor < len(m.artistList) {
			entry := m.artistList[m.artistCursor]
			m.songList = m.artistTracks(entry.name)
			m.filterLabel = "Artist: " + entry.name
			m.songCursor = 0
			m.pushView(viewSongs)
		}
	case viewAlbums:
		if m.albumCursor < len(m.albumList) {
			entry := m.albumList[m.albumCursor]
			m.songList = m.albumTracks(entry.name)
			m.filterLabel = "Album: " + entry.name
			m.songCursor = 0
			m.pushView(viewSongs)
		}
	case viewGenres:
		if m.genreCursor < len(m.genreList) {
			entry := m.genreList[m.genreCursor]
			m.songList = m.genreTracks(entry.name)
			m.filterLabel = "Genre: " + entry.name
			m.songCursor = 0
			m.pushView(viewSongs)
		}
	case viewSongs:
		if m.songCursor < len(m.songList) {
			m.player.SetQueue(m.songList, m.songCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
		}
	case viewSearch:
		if m.searchCursor < len(m.searchRes) {
			m.player.SetQueue(m.searchRes, m.searchCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
			m.searchActive = false
			m.searchInput.Blur()
		}
	case viewQueue:
		if m.queueCursor < len(m.player.Queue()) {
			m.player.SetQueue(m.player.Queue(), m.queueCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
		}
	case viewFavorites:
		if m.favCursor < len(m.favTracks) {
			m.player.SetQueue(m.favTracks, m.favCursor)
			_ = m.player.PlayCurrent()
			m.persistQueue()
		}
	}
}

func (m *Model) moveCursor(delta int) {
	cur := m.currentCursor()
	length := m.currentListLen()
	if length == 0 {
		return
	}
	newCur := cur + delta
	if newCur < 0 {
		newCur = 0
	}
	if newCur >= length {
		newCur = length - 1
	}
	m.setCursor(newCur)
}

func (m *Model) currentCursor() int {
	switch m.activeView {
	case viewArtists:
		return m.artistCursor
	case viewAlbums:
		return m.albumCursor
	case viewGenres:
		return m.genreCursor
	case viewSongs:
		return m.songCursor
	case viewSearch:
		return m.searchCursor
	case viewQueue:
		return m.queueCursor
	case viewFavorites:
		return m.favCursor
	}
	return 0
}

func (m *Model) currentListLen() int {
	switch m.activeView {
	case viewArtists:
		return len(m.artistList)
	case viewAlbums:
		return len(m.albumList)
	case viewGenres:
		return len(m.genreList)
	case viewSongs:
		return len(m.songList)
	case viewSearch:
		return len(m.searchRes)
	case viewQueue:
		return m.player.QueueLen()
	case viewFavorites:
		return len(m.favTracks)
	}
	return 0
}

func (m *Model) setCursor(v int) {
	if v < 0 {
		v = 0
	}
	switch m.activeView {
	case viewArtists:
		m.artistCursor = v
	case viewAlbums:
		m.albumCursor = v
	case viewGenres:
		m.genreCursor = v
	case viewSongs:
		m.songCursor = v
	case viewSearch:
		m.searchCursor = v
	case viewQueue:
		m.queueCursor = v
	case viewFavorites:
		m.favCursor = v
	}
}

func (m Model) statusLine() string {
	if m.lib == nil {
		return ""
	}

	status := fmt.Sprintf(" ♫ %d tracks · %d artists · %d albums",
		m.lib.FileCount,
		len(m.lib.Artists),
		len(m.lib.Albums),
	)

	if m.player.HasTrack() {
		track := m.player.CurrentTrack()
		state := "▮▮"
		if m.player.IsPaused() {
			state = "▶"
		}
		status += fmt.Sprintf("  │  %s %s", state, truncate(track.Title, 25))
	}

	return status
}

func (m Model) renderLoading() string {
	frame := SpinnerFrames[m.spinnerIdx%len(SpinnerFrames)]
	spinner := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render(frame + " Connecting to daemon & loading library…")

	info := DimStyle.Render(m.musicDir)

	box := PanelStyle.
		Width(60).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			"",
			LogoStyle.Render("♫  BPV  ♫"),
			"",
			spinner,
			info,
			"",
		))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m Model) renderError() string {
	errMsg := ErrorStyle.Render("Error: " + m.err.Error())

	box := PanelStyle.
		Width(60).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			"",
			LogoStyle.Render("♫  BPV  ♫"),
			"",
			errMsg,
			"",
			DimStyle.Render("Press q to quit"),
			"",
		))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func matchKey(msg tea.KeyMsg, binding key.Binding) bool {
	for _, k := range binding.Keys() {
		if msg.String() == k {
			return true
		}
	}
	return false
}
