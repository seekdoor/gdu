// +build live

package web

import (
	"log"
	"net/http"
	"os"
)

func getFileSystem() http.FileSystem {
	log.Print("using live mode")
	return http.FS(os.DirFS("web/static"))
}
