package web

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/net/websocket"

	"github.com/dundee/gdu/v5/analyze"
	"github.com/dundee/gdu/v5/device"
)

// UI is struct for web UI
type UI struct {
	analyzer         analyze.Analyzer
	showApparentSize bool
	ignoreDirPaths   map[string]struct{}
	topDir           *analyze.File
	topDirPath       string
	currentDir       *analyze.File
	currentDirPath   string
}

// CreateWebUI creates web UI
func CreateWebUI(showApparentSize bool) *UI {
	return &UI{
		showApparentSize: showApparentSize,
		analyzer:         analyze.CreateAnalyzer(),
	}
}

// StartUILoop stub
func (ui *UI) StartUILoop() error {
	addr := "127.0.0.1:8888"

	// err := openBrowser("http://" + addr)
	// if err != nil {
	// 	return err
	// }
	http.Handle("/ws", websocket.Handler(ui.handleWs))
	http.Handle("/", http.FileServer(getFileSystem()))
	log.Printf("Starting web server on %s", addr)
	return http.ListenAndServe(addr, nil)
}

// ListDevices lists mounted devices and shows their disk usage
func (ui *UI) ListDevices(getter device.DevicesInfoGetter) error {
	return nil
}

// AnalyzePath analyzes recursively disk usage in given path
func (ui *UI) AnalyzePath(path string, _ *analyze.File) {
	abspath, _ := filepath.Abs(path)

	go func() {
		ui.currentDir = ui.analyzer.AnalyzeDir(abspath, ui.ShouldDirBeIgnored)
		ui.topDirPath = abspath
		ui.topDir = ui.currentDir
		log.Println("Analysis done")
	}()
}

// SetIgnoreDirPaths sets paths to ignore
func (ui *UI) SetIgnoreDirPaths(paths []string) {
	ui.ignoreDirPaths = make(map[string]struct{}, len(paths))
	for _, path := range paths {
		ui.ignoreDirPaths[path] = struct{}{}
	}
}

// ShouldDirBeIgnored returns true if given path should be ignored
func (ui *UI) ShouldDirBeIgnored(path string) bool {
	_, ok := ui.ignoreDirPaths[path]
	return ok
}

func openBrowser(url string) error {
	var err error

	log.Println("Opening browser")

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}
