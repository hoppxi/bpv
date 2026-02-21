package tui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

var (
	coverMu    sync.Mutex
	coverCache = make(map[string]string)
)

func coverCacheKey(b64 string, w, h int) string {
	l := len(b64)
	if l == 0 {
		return ""
	}

	// Create a robust key using length + start + mid + end segments.
	// Just using the prefix is insufficient because many images share the same
	// file headers (magic bytes).
	start := 50
	if start > l {
		start = l
	}

	end := 50
	if l < 50 {
		end = 0
	}

	midStart := l / 2
	midEnd := midStart + 50
	if midEnd > l {
		midEnd = l
	}

	prefix := b64[:start]
	suffix := ""
	if end > 0 {
		suffix = b64[l-end:]
	}
	mid := b64[midStart:midEnd]

	return fmt.Sprintf("%d:%s:%s:%s:%d:%d", l, prefix, mid, suffix, w, h)
}

func renderCoverArt(coverB64 string, targetWidth, targetHeight int) string {
	if coverB64 == "" || targetWidth <= 0 || targetHeight <= 0 {
		return renderPlaceholderArt(targetWidth, targetHeight)
	}

	key := coverCacheKey(coverB64, targetWidth, targetHeight)
	coverMu.Lock()
	if cached, ok := coverCache[key]; ok {
		coverMu.Unlock()
		return cached
	}
	coverMu.Unlock()

	data, err := base64.StdEncoding.DecodeString(coverB64)
	if err != nil {
		return renderPlaceholderArt(targetWidth, targetHeight)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return renderPlaceholderArt(targetWidth, targetHeight)
	}

	result := renderImageToBlocks(img, targetWidth, targetHeight)

	coverMu.Lock()
	// Prevent unbounded growth
	if len(coverCache) > 50 {
		// Simple eviction: clear map.
		// For a TUI player, hitting 50 unique displayed covers is rare in one session
		// unless browsing aggressively.
		coverCache = make(map[string]string)
	}
	coverCache[key] = result
	coverMu.Unlock()

	return result
}

func renderImageToBlocks(img image.Image, targetWidth, targetHeight int) string {
	pixelRows := targetHeight * 2
	pixelCols := targetWidth

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	var lines []string
	for row := 0; row < pixelRows; row += 2 {
		var line strings.Builder
		line.Grow(pixelCols * 30)
		for col := 0; col < pixelCols; col++ {
			srcX := bounds.Min.X + col*srcW/pixelCols
			srcY1 := bounds.Min.Y + row*srcH/pixelRows
			srcY2 := bounds.Min.Y + (row+1)*srcH/pixelRows

			r1, g1, b1 := sampleColor(img, srcX, srcY1)
			r2, g2, b2 := sampleColor(img, srcX, srcY2)

			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorToHex(r1, g1, b1))).
				Background(lipgloss.Color(colorToHex(r2, g2, b2)))

			line.WriteString(style.Render("▀"))
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}

func sampleColor(img image.Image, x, y int) (uint8, uint8, uint8) {
	bounds := img.Bounds()
	if x < bounds.Min.X {
		x = bounds.Min.X
	}
	if x >= bounds.Max.X {
		x = bounds.Max.X - 1
	}
	if y < bounds.Min.Y {
		y = bounds.Min.Y
	}
	if y >= bounds.Max.Y {
		y = bounds.Max.Y - 1
	}

	c := img.At(x, y)
	r, g, b, _ := c.RGBA()
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

func colorToHex(r, g, b uint8) string {
	const hex = "0123456789abcdef"
	return "#" + string([]byte{hex[r>>4], hex[r&0x0f], hex[g>>4], hex[g&0x0f], hex[b>>4], hex[b&0x0f]})
}

func renderPlaceholderArt(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	noteWidth := 3
	centerRow := height / 2
	centerCol := width / 2
	startNoteCol := centerCol - 2
	if startNoteCol < 0 {
		startNoteCol = 0
	}
	showNote := height > 2 && width > 4

	var lines []string
	for row := 0; row < height; row++ {
		var line strings.Builder
		for col := 0; col < width; col++ {
			if showNote && row == centerRow && col >= startNoteCol && col < startNoteCol+noteWidth {
				if col == startNoteCol {
					noteStyle := lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color("#c084fc")).
						Background(lipgloss.Color("#1a0d38"))
					line.WriteString(noteStyle.Render(" ♫ "))
				}
				continue
			}

			t := float64(row) / float64(height)
			ct := float64(col) / float64(width)

			r := uint8(20 + t*40 + ct*15)
			g := uint8(5 + t*10)
			b := uint8(60 + t*80 + ct*20)

			style := lipgloss.NewStyle().
				Background(lipgloss.Color(colorToHex(r, g, b)))
			line.WriteString(style.Render(" "))
		}
		lines = append(lines, line.String())
	}

	return strings.Join(lines, "\n")
}
