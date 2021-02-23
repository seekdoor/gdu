package web

import (
	"net/http"

	"github.com/dundee/gdu/v5/analyze"
	"github.com/dundee/gdu/v5/device"
)

// UI is struct for web UI
type UI struct {
	showApparentSize bool
	ignoreDirPaths   map[string]struct{}
}

// CreateWebUI creates web UI
func CreateWebUI(showApparentSize bool) *UI {
	return &UI{
		showApparentSize: showApparentSize,
	}
}

// StartUILoop stub
func (ui *UI) StartUILoop() error {
	http.Handle("/", http.FileServer(getFileSystem()))
	return http.ListenAndServe(":8888", nil)
}

// ListDevices lists mounted devices and shows their disk usage
func (ui *UI) ListDevices(getter device.DevicesInfoGetter) error {
	return nil
}

// AnalyzePath analyzes recursively disk usage in given path
func (ui *UI) AnalyzePath(path string, _ *analyze.File) {

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
