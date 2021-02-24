package runner

import (
	"bufio"
	"context"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cornelk/hashmap"
	aurora "github.com/logrusorgru/aurora/v3"
	"github.com/olekukonko/tablewriter"
	"github.com/projectdiscovery/gologger"
	"github.com/projectpandora/depseeker/common/depseeker"
	"github.com/projectpandora/depseeker/common/fileutil"
)

// Runner is a client for running the enumeration process.
type Runner struct {
	options    *Options
	depseeker  *depseeker.Depseeker
	hm         hashmap.HashMap
	numTargets uint64
}

// New creates a new client for running enumeration process.
func New(options *Options) (*Runner, error) {
	runner := &Runner{
		options: options,
		hm:      hashmap.HashMap{},
	}

	depseekerOptions := depseeker.DefaultOptions
	depseekerOptions.Timeout = time.Duration(options.Timeout) * time.Second
	depseekerOptions.HTTPProxy = options.HTTPProxy
	depseekerOptions.UserAgent = options.UserAgent

	var key, value string
	depseekerOptions.CustomHeaders = make(map[string]string)
	for _, customHeader := range options.CustomHeaders {
		tokens := strings.SplitN(customHeader, ":", 2)

		// continue normally
		if len(tokens) < 2 {
			continue
		}
		key = strings.TrimSpace(tokens[0])
		value = strings.TrimSpace(tokens[1])
		depseekerOptions.CustomHeaders[key] = value
	}

	var err error
	runner.depseeker, err = depseeker.New(&depseekerOptions)
	if err != nil {
		gologger.Fatal().Msgf("Could not create depseeker instance: %s\n", err)
	}

	return runner, nil
}

// Close the instance
func (runner *Runner) Close() {
}

func (runner *Runner) prepareInput() {
	var (
		finput  *os.File
		scanner *bufio.Scanner
		err     error
	)
	// handle single target
	if runner.options.InputTarget != "" {
		runner.numTargets = 1
		// nolint:errcheck // ignoring error
		runner.hm.Set(runner.options.InputTarget, nil)
	} else {
		// check if file has been provided
		if fileutil.FileExists(runner.options.InputFile) {
			finput, err = os.Open(runner.options.InputFile)
			if err != nil {
				gologger.Fatal().Msgf("Could read input file '%s': %s\n", runner.options.InputFile, err)
			}
			scanner = bufio.NewScanner(finput)
		} else if fileutil.HasStdin() {
			scanner = bufio.NewScanner(os.Stdin)
		} else {
			// get from target
			gologger.Fatal().Msgf("No input provided")
		}
		var numTargets uint64
		for scanner.Scan() {
			target := strings.TrimSpace(scanner.Text())
			// get the exact number of targets
			if _, ok := runner.hm.Get(target); ok {
				continue
			}
			numTargets++
			// nolint:errcheck // ignore
			runner.hm.Set(target, nil)
		}
		runner.numTargets = numTargets
	}

	if runner.options.InputFile != "" {
		err := finput.Close()
		if err != nil {
			gologger.Fatal().Msgf("Could close input file '%s': %s\n", runner.options.InputFile, err)
		}
	}
}

// RunEnumeration on targets
func (runner *Runner) RunEnumeration(ctx context.Context) {
	runner.prepareInput()

	chanTargets := runner.hm.Iter()
	var wg sync.WaitGroup
	for i := 0; i < runner.options.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			target, more := <-chanTargets
			if more {
				ctxRun, cancel := context.WithTimeout(ctx, time.Second*time.Duration(runner.options.Timeout))
				defer cancel()
				sTarget, ok := target.Key.(string)
				if ok {
					dependencies, err := runner.depseeker.Run(ctxRun, sTarget)
					if err == nil {
						if len(dependencies) == 0 {
							gologger.Print().Msgf("[+] %v process completed, found 0 dependencies.\n", target.Key)
						} else {
							gologger.Print().Msgf("[+] %v process completed, found %d dependencies.\n", target.Key, aurora.Yellow(len(dependencies)))
							table := tablewriter.NewWriter(os.Stdout)
							table.SetHeader([]string{"Name", "Version", "Private"})
							for _, dependency := range dependencies {
								rowName := dependency.Name
								rowVersion := dependency.Version
								rowPrivate := "No"
								if dependency.IsPrivate == true {
									rowName = aurora.Red(dependency.Name).String()
									rowVersion = aurora.Red(dependency.Version).String()
									rowPrivate = aurora.Red("Yes").String()
								}
								table.Append([]string{rowName, rowVersion, rowPrivate})
							}
							table.Render()
						}
					} else {
						gologger.Error().Msgf("Error = %v\n", err)
					}
				}
			}
		}()
	}
	wg.Wait()
}
