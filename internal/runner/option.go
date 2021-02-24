package runner

import (
	"flag"
	"os"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectpandora/depseeker/common/customheader"
	"github.com/projectpandora/depseeker/common/depseeker"
)

// Options of the internal runner
type Options struct {
	Verbose       bool
	NoColor       bool
	Silent        bool
	Version       bool
	HTTPProxy     string
	Timeout       int
	Threads       int
	UserAgent     string
	CustomHeaders customheader.CustomHeaders
	InputTarget   string
	InputFile     string
}

// ParseOptions parses the command line options for application
func ParseOptions() *Options {
	options := &Options{}

	flag.IntVar(&options.Threads, "threads", 16, "Number of threads")
	flag.IntVar(&options.Timeout, "timeout", 60, "Timeout in seconds")
	flag.StringVar(&options.UserAgent, "user-agent", depseeker.DefaultUserAgent, "User agent")
	flag.Var(&options.CustomHeaders, "H", "Custom headers")
	flag.StringVar(&options.HTTPProxy, "http-proxy", "", "HTTP Proxy, eg http://127.0.0.1:8080")
	flag.BoolVar(&options.Silent, "silent", false, "Silent mode")
	flag.BoolVar(&options.Version, "version", false, "Show version of depseeker")
	flag.BoolVar(&options.Verbose, "verbose", false, "Verbose mode")
	flag.BoolVar(&options.NoColor, "no-color", false, "No color")
	flag.StringVar(&options.InputFile, "l", "", "File containing urls")
	flag.StringVar(&options.InputTarget, "target", "", "A single url to run")

	flag.Parse()

	// read the inputs and configure the logging
	options.configureOutput()

	if options.Silent == false {
		showBanner()
	}

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", Version)
		os.Exit(0)
	}

	options.validateOptions()

	return options
}

func (options *Options) validateOptions() {
}

func (options *Options) configureOutput() {
	gologger.DefaultLogger.SetFormatter(formatter.NewCLI(options.NoColor))
	if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}
	// } else if options.Silent {
	// 	gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	// }
}
