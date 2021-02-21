package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/dundee/gdu/v5/build"
	"github.com/dundee/gdu/v5/cmd/app"
	"github.com/dundee/gdu/v5/common"
	"github.com/dundee/gdu/v5/device"
	"github.com/dundee/gdu/v5/stdout"
	"github.com/dundee/gdu/v5/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-isatty"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var af *app.Flags

var rootCmd = &cobra.Command{
	Use:   "gdu [directory_to_scan]",
	Short: "Pretty fast disk usage analyzer written in Go",
	Long: `Pretty fast disk usage analyzer written in Go.

Gdu is intended primarily for SSD disks where it can fully utilize parallel processing.
However HDDs work as well, but the performance gain is not so huge.
`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui, err := createStdoutUI()
		if err != nil {
			return err
		}
		return runApp(args, ui, app.ActionAnalyzePath)
	},
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of gdu",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t", build.Version)
		fmt.Println("Built time:\t", build.Time)
		fmt.Println("Built user:\t", build.User)
	},
}

var interactiveCmd = &cobra.Command{
	Use:     "interactive [directory_to_scan]",
	Aliases: []string{"i"},
	Short:   "Run in interactive mode",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ui, fini, err := createTui()
		defer fini()
		if err != nil {
			return err
		}
		return runApp(args, ui, app.ActionAnalyzePath)
	},
}

var disksCmd = &cobra.Command{
	Use:     "disks",
	Aliases: []string{"d"},
	Short:   "Show all mounted disks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui, err := createStdoutUI()
		if err != nil {
			return err
		}
		return runApp(args, ui, app.ActionListDevices)
	},
}

var interactiveDisksCmd = &cobra.Command{
	Use:     "disks",
	Aliases: []string{"d"},
	Short:   "Show all mounted disks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui, fini, err := createTui()
		defer fini()
		if err != nil {
			return err
		}
		return runApp(args, ui, app.ActionListDevices)
	},
}

func init() {
	af = &app.Flags{}
	rootFlags := rootCmd.PersistentFlags()
	rootFlags.StringVarP(&af.LogFile, "log-file", "l", "/dev/null", "Path to a logfile")
	rootFlags.StringSliceVarP(&af.IgnoreDirs, "ignore-dirs", "i", []string{"/proc", "/dev", "/sys", "/run"}, "Absolute paths to ignore (separated by comma)")
	rootFlags.BoolVarP(&af.ShowApparentSize, "show-apparent-size", "a", false, "Show apparent size")
	rootFlags.BoolVarP(&af.NoColor, "no-color", "c", false, "Do not use colorized output")
	rootFlags.BoolVarP(&af.NoCross, "no-cross", "x", false, "Do not cross filesystem boundaries")

	rootCmd.Flags().BoolVarP(&af.NoProgress, "no-progress", "p", false, "Do not show progress in non-interactive mode")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(disksCmd)

	interactiveCmd.AddCommand(interactiveDisksCmd)

	// we are not able to analyze disk usage on Windows and Plan9
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		af.ShowApparentSize = true
	}
	if runtime.GOOS == "windows" && af.LogFile == "/dev/null" {
		af.LogFile = "nul"
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getPath(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return "."
}

func createTermApp() (common.TermApplication, func(), error) {
	istty := isatty.IsTerminal(os.Stdout.Fd())

	if !istty {
		return nil, nil, errors.New("Interactive mode cannot be started, not running in valid TTY")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating screen: %w", err)
	}
	screen.Init()
	fini := func() {
		defer screen.Clear()
		defer screen.Fini()
	}

	termApp := tview.NewApplication()
	termApp.SetScreen(screen)

	if !af.NoColor {
		tview.Styles.TitleColor = tcell.NewRGBColor(27, 161, 227)
	}
	return termApp, fini, nil
}

func createStdoutUI() (common.UI, error) {
	if err := setLogging(); err != nil {
		return nil, err
	}
	istty := isatty.IsTerminal(os.Stdout.Fd())
	ui := stdout.CreateStdoutUI(
		os.Stdout,
		!af.NoColor && istty,
		!af.NoProgress && istty,
		af.ShowApparentSize,
	)
	return ui, nil
}

func createTui() (common.UI, func(), error) {
	if err := setLogging(); err != nil {
		return nil, nil, err
	}
	termApp, fini, err := createTermApp()
	if err != nil {
		return nil, nil, err
	}
	ui := tui.CreateUI(termApp, !af.NoColor, af.ShowApparentSize)
	return ui, fini, nil
}

func runApp(args []string, ui common.UI, action func(*app.App) error) error {
	a := app.CreateApp(
		getPath(args),
		af,
		ui,
		device.Getter,
	)
	a.SetAction(action)
	return a.Run()
}

func setLogging() error {
	f, err := os.OpenFile(af.LogFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("Error opening log file: %w", err)
	}
	defer f.Close()
	log.SetOutput(f)
	return nil
}
