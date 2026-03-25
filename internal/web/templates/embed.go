package templates

import (
	"embed"
	"io/fs"
)

//go:embed *.html
var templatesFS embed.FS

func FS() fs.FS {
	return templatesFS
}
