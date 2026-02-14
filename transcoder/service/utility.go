package service

import (
	"path/filepath"
	"strings"
)

func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/MP2T"
	case ".mpd":
		return "application/dash+xml"
	case ".m4s":
		return "video/mp4"
	default:
		return "application/octet-stream"
	}
}
