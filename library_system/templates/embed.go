package templates

import (
	"embed"
	"github.com/gobuffalo/buffalo"
	"io/fs"
)

//go:embed *.html */*.html
var files embed.FS

func FS() fs.FS {
	return buffalo.NewFS(files, "templates")
}
