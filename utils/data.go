package utils

import "path"

var dataKinds = map[string]string{
	".mp4":  "video/mp4",
	".m3u8": "application/x-mpegURL",
	".mkv":  "video/x-matroska",
	"==":    "file",
}

func getDataKind(link string) string {
	if ext, ok := dataKinds[path.Ext(link)]; ok {
		return ext
	}
	return "video/*"
}
