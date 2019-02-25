package utils

import "path"

var dataKinds = map[string]string{
	".mp4":  "video/mp4",
	".m3u8": "application/x-mpegURL",
	".mkv":  "video/x-matroska",
	"==":    "file",
}

func getDataKind(link string) string {
	_, file := path.Split(link)
	if ext, ok := dataKinds[file]; ok {
		return ext
	}
	return "video/*"
}
