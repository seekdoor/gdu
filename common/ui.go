package common

import (
	"github.com/dundee/gdu/v5/analyze"
	"github.com/dundee/gdu/v5/device"
)

// UI is common interface for both terminal UI and text output
type UI interface {
	ListDevices(getter device.DevicesInfoGetter) error
	AnalyzePath(path string, parentDir *analyze.File)
	SetIgnoreDirPaths(paths []string)
	StartUILoop() error
}
