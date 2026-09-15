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

// findOpenSpecProject walks up from startDir looking for a directory containing openspec/.
func findOpenSpecProject(startDir string) string {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}
	for {
		candidate := filepath.Join(dir, "openspec")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "" // reached root
		}
		dir = parent
	}
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

func main() {
	app := cli.NewApp()
	app.Name = "specgetty"
	app.Version = version
	app.Usage = "Finds OpenSpec projects on your local machine"
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, cmd string) {
		fmt.Printf("ERROR: Unknown command '%s'\n", cmd)
	}
	app.Flags = appFlags()

	app.Action = func(c *cli.Context) error {

		config, err := scanner.ParseConfigFile(c.String("config"), defaultConfig)
		if c.Args().Len() > 0 {

			fmt.Println("Arguments given, skipping config")
			config.ScanDirs.Include = c.Args().Slice()

		} else {
			if err != nil {
				return err
			}
			for i := range config.ScanDirs.Include {
				config.ScanDirs.Include[i] = os.ExpandEnv(config.ScanDirs.Include[i])
			}
			for i := range config.ScanDirs.Exclude {
				config.ScanDirs.Exclude[i] = os.ExpandEnv(config.ScanDirs.Exclude[i])
			}
		}

		if c.Bool("debug") {
			var projects scanner.ProjectMap
			projects, err = scanner.Scan(config, c.Bool("ignore_dir_errors"))
			if err != nil {
				panic(err)
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

		// Resolve the startup project. This never walks the configured scan
		// directories: discovery only happens when the picker asks for it.
		var startupPath string
		if startView == "single" {
			if c.String("path") != "" {
				startupPath, _ = filepath.Abs(c.String("path"))
			} else {
				cwd, _ := os.Getwd()
				startupPath = findOpenSpecProject(cwd)
			}
		}

		fields, err := ui.ResolveFields(c.String("change-fields"), config.ChangeFields)
		if err != nil {
			return err
		}

		err = ui.Run(config, c.Bool("ignore_dir_errors"), version, startupPath, startView, fields)
		if err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}
