package cmd

import (
	// core
	"flag"
	"fmt"
	"os"

	// gv package
	"github.com/myENA/gv"

	// CLI library
	"github.com/mitchellh/cli"
)

// Command is a Command implementation that returns version information
type Command struct {
	buildInfo *gv.BuildInfo
	ui        cli.Ui
}

// New builds and returns a Command struct
func New(bi *gv.BuildInfo, ui cli.Ui) (*Command, error) {
	return &Command{
		buildInfo: bi,
		ui:        ui,
	}, nil
}

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error    // error holder
	var details bool // detail toggle

	// init flagset
	cmdFlags := flag.NewFlagSet("gv", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.ui.Output(c.Help()); os.Exit(0) }

	// set flags
	cmdFlags.BoolVar(&details, "d", false, "Display detailed dependency information")

	// check and parse
	if err = cmdFlags.Parse(args); err != nil {
		return 1
	}

	// show output
	return c.showVersion(details)
}

// show information output
func (c *Command) showVersion(detail bool) int {
	// print app name and version - should always be present
	c.ui.Output(fmt.Sprintf("==>\t%s v%s",
		c.buildInfo.Name,
		c.buildInfo.Version,
	))

	// print build if present
	if c.buildInfo.Build != "" {
		c.ui.Output(fmt.Sprintf("Build:\t%s",
			c.buildInfo.Build,
		))
	}

	// print branch if present
	if c.buildInfo.Branch != "" {
		c.ui.Output(fmt.Sprintf("Branch:\t%s",
			c.buildInfo.Branch,
		))
	}

	// print commit if present
	if c.buildInfo.Commit != "" {
		c.ui.Output(fmt.Sprintf("Commit:\t%s",
			c.buildInfo.Commit,
		))
	}

	// print date if present
	if c.buildInfo.Date != "" {
		c.ui.Output(fmt.Sprintf("Date:\t%s",
			c.buildInfo.Date,
		))
	}

	// show details if asked
	if detail {
		// ensure we have glide data
		if c.buildInfo.Init() != nil {
			// bail out
			c.ui.Error("Failed to init build information")
			return 1
		}
		// check struct
		if c.buildInfo.GlideLockfile != nil {
			// print header
			c.ui.Output(fmt.Sprintf("\nImports %d Packages",
				len(c.buildInfo.GlideLockfile.Imports)))
			// loop through imports
			for idx, pkg := range c.buildInfo.GlideLockfile.Imports {
				c.ui.Output(fmt.Sprintf("[%03d] % 8s %s", idx+1, pkg.Version[0:8], pkg.Name))
			}
		}
	}

	// all good
	return 0
}

// Synopsis shows the command summary
func (c *Command) Synopsis() string {
	return "Display application version information"
}

// Help shows the detailed command options
func (c *Command) Help() string {
	str := `Usage: %s version [options]

	Display application version and dependency information.

Options:

	-d    Display detailed dependency information
`
	// return help
	return fmt.Sprintf(str, c.buildInfo.Name)
}
