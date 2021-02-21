package app

import (
	"fmt"

	"github.com/dundee/gdu/v5/common"
	"github.com/dundee/gdu/v5/device"
)

// Flags define flags accepted by Run
type Flags struct {
	LogFile          string
	IgnoreDirs       []string
	ShowApparentSize bool
	NoColor          bool
	NoProgress       bool
	NoCross          bool
}

// App defines the main application
type App struct {
	path   string
	flags  *Flags
	getter device.DevicesInfoGetter
	ui     common.UI
	action Action
}

// Action is the action to be performed
type Action func(a *App) error

// CreateApp creates app struct
func CreateApp(path string, flags *Flags, ui common.UI, getter device.DevicesInfoGetter) *App {
	return &App{
		path:   path,
		ui:     ui,
		flags:  flags,
		getter: getter,
	}
}

// SetAction sets action to be run
func (a *App) SetAction(action Action) {
	a.action = action
}

// Run runs the app
func (a *App) Run() error {
	if err := a.setNoCross(a.path); err != nil {
		return err
	}
	a.ui.SetIgnoreDirPaths(a.flags.IgnoreDirs)

	if err := a.action(a); err != nil {
		return err
	}

	return a.ui.StartUILoop()
}

func (a *App) setNoCross(path string) error {
	if a.flags.NoCross {
		mounts, err := a.getter.GetMounts()
		if err != nil {
			return fmt.Errorf("Error loading mount points: %w", err)
		}
		paths := device.GetNestedMountpointsPaths(path, mounts)
		a.flags.IgnoreDirs = append(a.flags.IgnoreDirs, paths...)
	}
	return nil
}

// ActionAnalyzePath analyzes given path
func ActionAnalyzePath(a *App) error {
	a.ui.AnalyzePath(a.path, nil)
	return nil
}

// ActionListDevices list devices and shows their usage
func ActionListDevices(a *App) error {
	if err := a.ui.ListDevices(a.getter); err != nil {
		return fmt.Errorf("Error loading mount points: %w", err)
	}
	return nil

}
