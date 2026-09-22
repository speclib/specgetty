package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/mipmip/specgetty/src/scanner"
	"github.com/mipmip/specgetty/src/ui"
)

func getDefaultConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configDir, "specgetty", "config.yml")
}

// findOpenSpecProject walks up from startDir looking for a directory whose
// openspec/ is one worth stopping at.
//
// A directory named openspec with nothing in it does not qualify. OpenSpec
// added that rule when stores arrived: the recommended store layout puts one
// at ~/openspec, and without the qualification that empty shell would make the
// home directory capture every project beneath it.
func findOpenSpecProject(startDir string) string {
	return scanner.FindRoot(startDir)
}

// expandScanDirs expands environment variables in the configured scan
// directories, so a config can be written with $HOME rather than a path that
// only works for one person.
func expandScanDirs(config *scanner.Config) {
	for i := range config.ScanDirs.Include {
		config.ScanDirs.Include[i] = os.ExpandEnv(config.ScanDirs.Include[i])
	}
	for i := range config.ScanDirs.Exclude {
		config.ScanDirs.Exclude[i] = os.ExpandEnv(config.ScanDirs.Exclude[i])
	}
}

// resolveStartupPath finds the project to open at startup.
//
// It never walks the configured scan directories: an explicit path is taken as
// given, and otherwise the search walks up from the working directory, which
// costs a few stats. Discovery only happens when the picker asks for it.
func resolveStartupPath(startView, pathFlag string) string {
	if startView != "single" {
		return ""
	}
	if pathFlag != "" {
		abs, err := filepath.Abs(pathFlag)
		if err != nil {
			return ""
		}
		return abs
	}
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return findOpenSpecProject(cwd)
}

// resolveStartView validates the --view value and applies the rule that an
// explicit --path means a single project, whatever --view said.
func resolveStartView(view, path string) (string, error) {
	if view != "single" && view != "all" {
		return "", fmt.Errorf("unknown --view value %q; valid values are: single, all", view)
	}
	if path != "" {
		return "single", nil
	}
	return view, nil
}

//go:embed config.yml
var defaultConfig string

//go:embed VERSION
var version string

func init() {
	version = strings.TrimSpace(version)
}

// appFlags is the command-line surface. It is a function so that tests can
// assert what it does and does not contain: --zoom is gone, and staying gone is
// part of the contract.
func appFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Location of config file",
			Value:   getDefaultConfigPath(),
		},

		&cli.BoolFlag{
			Name:    "ignore_dir_errors",
			Aliases: []string{"i"},
			Value:   true,
			Usage:   "Don't halt on errors while finding dirs",
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug output instead of UI",
		},
		&cli.StringFlag{
			Name:  "view",
			Value: "single",
			Usage: "Which view opens first: single (the project at the working directory) or all (the project picker)",
		},
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Usage:   "OpenSpec project to open; implies --view=single",
		},
		&cli.StringFlag{
			Name:  "change-fields",
			Usage: "Comma-separated columns for the change list (overrides change_fields in the config)",
		},
	}
}

// runApp is the whole command-line surface's behaviour, separated from main so
// that it can be driven by a test. Everything up to the point the terminal is
// taken over is reachable that way, which `--debug` exercises end to end.
func runApp(c *cli.Context) error {
	// Refused rather than ignored. urfave/cli hands a positional argument
	// through without complaint, so saying nothing would open the working
	// directory and look like the argument had been honoured.
	if c.Args().Len() > 0 {
		return fmt.Errorf("%s: directories are not accepted as arguments. "+
			"Use --path to open one project, or scandirs.include in the "+
			"configuration file to choose what the picker searches",
			c.Args().First())
	}

	config, err := scanner.ParseConfigFile(c.String("config"), defaultConfig)
	if err != nil {
		return err
	}
	expandScanDirs(config)

	// A setting that no longer has an effect is reported rather than swallowed.
	// YAML ignores keys it does not know, so without this a released config
	// would keep a line that quietly does nothing.
	for _, key := range scanner.RetiredKeysIn(c.String("config")) {
		fmt.Printf("Note: %s in %s no longer has any effect: %s\n",
			key, c.String("config"), scanner.RetiredConfigKeys[key])
	}

	if c.Bool("debug") {
		projects, err := scanner.Scan(config, c.Bool("ignore_dir_errors"))
		if err != nil {
			return err
		}

		for r, st := range projects {
			fmt.Printf("%-40s %v\n", r, st.ScanTime)
		}
		return nil
	}

	startView, err := resolveStartView(c.String("view"), c.String("path"))
	if err != nil {
		return err
	}

	startupPath := resolveStartupPath(startView, c.String("path"))

	fields, err := ui.ResolveFields(c.String("change-fields"), config.ChangeFields)
	if err != nil {
		return err
	}

	return ui.Run(config, c.Bool("ignore_dir_errors"), version, startupPath, startView, fields)
}

// newApp builds the command-line application.
func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "specgetty"
	app.Version = version
	app.Usage = "Finds OpenSpec projects on your local machine"
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		fmt.Printf("ERROR: Unknown command '%s'\n", cmd)
	}
	app.Flags = appFlags()
	app.Action = runApp
	return app
}

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
