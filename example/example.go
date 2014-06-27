package main

import (
	"github.com/stuartnelson3/guac"
	"github.com/stuartnelson3/guac/concat"
)

func main() {
	guac.WatchPath("./public/js", concat.Concat("./public/application.js", "./public/js", ".js"))
	guac.WatchPath("./public/css", concat.Concat("./public/application.css", "./public/css", ".css"))

	guac.Run()
}
