package static

import (
	"embed"
)

//go:embed *.html js/*.js
var FS embed.FS
