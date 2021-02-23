package runner

import (
	"flag"
	"os"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectpandora/depseeker/common/customheader"
)

// Options of the internal runner
type Options struct {
	Verbose             bool
	NoColor             bool
	Silent              bool
	Version             bool
	FollowRedirects     bool
	FollowHostRedirects bool
	HTTPProxy           string
	SocksProxy          string
	Timeout             int
	Threads             int
	Retries             int
	CustomHeaders       customheader.CustomHeaders
	InputTarget         string
	InputFile           string
}

// ParseOptions parses the command line options for application
func ParseOptions() *Options {
	options := &Options{}

	flag.IntVar(&options.Threads, "threads", 50, "Number of threads")
	flag.IntVar(&options.Retries, "retries", 0, "Number of retries")
	flag.IntVar(&options.Timeout, "timeout", 5, "Timeout in seconds")
	flag.Var(&options.CustomHeaders, "H", "Custom Header")
	flag.BoolVar(&options.FollowRedirects, "follow-redirects", false, "Follow Redirects")
	flag.BoolVar(&options.FollowHostRedirects, "follow-host-redirects", false, "Only follow redirects on the same host")
	flag.StringVar(&options.HTTPProxy, "http-proxy", "", "HTTP Proxy, eg http://127.0.0.1:8080")
	flag.BoolVar(&options.Silent, "silent", false, "Silent mode")
	flag.BoolVar(&options.Version, "version", false, "Show version of depseeker")
	flag.BoolVar(&options.Verbose, "verbose", false, "Verbose Mode")
	flag.BoolVar(&options.NoColor, "no-color", false, "No Color")
	flag.StringVar(&options.InputFile, "l", "", "File containing urls")
	flag.StringVar(&options.InputTarget, "target", "", "Target is a single target to scan using template")

	flag.Parse()

	// read the inputs and configure the logging
	options.configureOutput()

	showBanner()

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
	} else if options.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}
}
