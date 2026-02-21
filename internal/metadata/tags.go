package metadata

import (
	"path/filepath"
	"strings"
)

var SupportedFormats = map[string]string{
	"mp3":  "MPEG Audio Layer III",
	"flac": "Free Lossless Audio Codec",
	"wav":  "Waveform Audio File Format",
	"ogg":  "Ogg Vorbis",
	"m4a":  "MPEG-4 Audio",
	"aac":  "Advanced Audio Coding",
	"wma":  "Windows Media Audio",
	"opus": "Opus Audio",
	"alac": "Apple Lossless Audio Codec",
	"aiff": "Audio Interchange File Format",
	"dsf":  "DSD Stream File",
}

func IsSupportedFormat(format string) bool {
	format = strings.ToLower(format)
	_, exists := SupportedFormats[format]
	return exists
}

func GetFormatFromExtension(filename string) string {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if IsSupportedFormat(ext) {
		return strings.ToLower(ext)
	}
	return ""
}

func GetFormatDescription(format string) string {
	if desc, exists := SupportedFormats[strings.ToLower(format)]; exists {
		return desc
	}
	return "Unknown audio format"
}

type FormatInfo struct {
	Extension   string `json:"extension"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Lossless    bool   `json:"lossless"`
}

func GetSupportedFormats() []FormatInfo {
	var formats []FormatInfo
	
	for ext, desc := range SupportedFormats {
		formats = append(formats, FormatInfo{
			Extension:   ext,
			Name:        strings.ToUpper(ext),
			Description: desc,
			Lossless:    isLosslessFormat(ext),
		})
	}
	
	return formats
}

func isLosslessFormat(format string) bool {
	losslessFormats := map[string]bool{
		"flac": true,
		"wav":  true,
		"alac": true,
		"aiff": true,
		"dsf":  true,
	}
	return losslessFormats[strings.ToLower(format)]
}