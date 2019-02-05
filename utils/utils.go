package utils

import (
	"path"
	"strings"
)

var dataKinds = map[string]string{
	".mp4": "video/mp4",
	".m3u8": "application/x-mpegURL",
	".mkv": "video/x-matroska",
	"==": "file",
	_: "",
}

func getDataKind(link string) {
	_, file := path.Split(link)

}
