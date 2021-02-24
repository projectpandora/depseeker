package main

import (
	"context"

	"github.com/projectdiscovery/gologger"
	"github.com/projectpandora/depseeker/internal/runner"
)

func main() {
	// Parse the command line flags
	options := runner.ParseOptions()

	run, err := runner.New(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	run.RunEnumeration(context.Background())
	run.Close()
}
