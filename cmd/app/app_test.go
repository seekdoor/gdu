package app

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dundee/gdu/v5/device"
	"github.com/dundee/gdu/v5/internal/testdev"
	"github.com/dundee/gdu/v5/internal/testdir"
	"github.com/dundee/gdu/v5/stdout"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzePath(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	out, err := runApp(
		&Flags{LogFile: "/dev/null"},
		"test_dir",
		false,
		testdev.DevicesInfoGetterMock{},
		ActionAnalyzePath,
	)

	assert.Contains(t, out, "nested")
	assert.Nil(t, err)
}

func TestNoCross(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	out, err := runApp(
		&Flags{LogFile: "/dev/null", NoCross: true},
		"test_dir",
		false,
		testdev.DevicesInfoGetterMock{},
		ActionAnalyzePath,
	)

	assert.Contains(t, out, "nested")
	assert.Nil(t, err)
}

func TestNoCrossWithErr(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	out, err := runApp(
		&Flags{LogFile: "/dev/null", NoCross: true},
		"test_dir",
		false,
		device.LinuxDevicesInfoGetter{MountsPath: "/xxxyyy"},
		ActionAnalyzePath,
	)

	assert.Equal(t, "Error loading mount points: open /xxxyyy: no such file or directory", err.Error())
	assert.Empty(t, out)
}

func TestListDevices(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	out, err := runApp(
		&Flags{LogFile: "/dev/null"},
		"",
		false,
		testdev.DevicesInfoGetterMock{},
		ActionListDevices,
	)

	assert.Contains(t, out, "Device")
	assert.Nil(t, err)
}

func TestListDevicesWithErr(t *testing.T) {
	fin := testdir.CreateTestDir()
	defer fin()

	_, err := runApp(
		&Flags{LogFile: "/dev/null"},
		"",
		false,
		device.LinuxDevicesInfoGetter{MountsPath: "/xxxyyy"},
		ActionListDevices,
	)

	assert.Equal(t, "Error loading mount points: open /xxxyyy: no such file or directory", err.Error())
}

func runApp(flags *Flags, path string, istty bool, getter device.DevicesInfoGetter, action Action) (string, error) {
	buff := bytes.NewBufferString("")
	ui := stdout.CreateStdoutUI(
		buff,
		!flags.NoColor && istty,
		!flags.NoProgress && istty,
		flags.ShowApparentSize,
	)

	app := CreateApp(path, flags, ui, getter)
	app.SetAction(action)
	err := app.Run()

	return strings.TrimSpace(buff.String()), err
}
