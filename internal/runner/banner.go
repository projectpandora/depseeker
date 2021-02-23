package runner

import "github.com/projectdiscovery/gologger"

const banner = `
___  ___ ___  ___ ___ ___ _  _____ ___ 
|   \| __| _ \/ __| __| __| |/ / __| _ \
| |) | _||  _/\__ \ _|| _|| ' <| _||   /
|___/|___|_|  |___/___|___|_|\_\___|_|_\ v1.0.0						 
`

// Version is the current version
const Version = `1.0.0`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)

	gologger.Print().Msg("Use with caution. You are responsible for your actions\n")
	gologger.Print().Msg("Developers assume no liability and are not responsible for any misuse or damage.\n")
}
