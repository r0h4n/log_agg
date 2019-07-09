// log aggregation service.
//
// To start log_agg:
//
//  log_agg
//
// Usage information, refer to the help doc `log_agg -h`:
//
//  Usage:
//    log_agg [flags]
//
//
//  Flags:
//    -c, --config-file string    config file location for log_agg
//    -d, --db-address string     Log storage address (default "boltdb:///var/db/log_agg.bolt")
//    -a, --listen-http string    API listen address (same endpoint for http log collection) (default "0.0.0.0:6360")
//    -k, --log-keep string       Age or number of logs to keep per type '{"app":"2w", "deploy": 10}' (int or X(m)in, (h)our,  (d)ay, (w)eek, (y)ear) (default "{\"app\":\"2w\"}")
//    -l, --log-level string      Level at which to log (default "info")
//    -L, --log-type string       Default type to apply to incoming logs (commonly used: app|deploy) (default "app")
//    -v, --version               Print version info and exit
//
package main

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/r0h4n/log_agg/api"
	"github.com/r0h4n/log_agg/input"
	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
	"github.com/r0h4n/log_agg/output"
)

var (
	configFile string
	portFile   string


	// provides the log_agg server functionality
	Log_agg = &cobra.Command{
		Use:               "log_agg",
		Short:             "log_agg logging server",
		Long:              ``,
		PersistentPreRunE: readConfig,
		PreRunE:           preFlight,
		RunE:              startLog_agg,
		SilenceErrors:     true,
		SilenceUsage:      true,
	}

	// version information (populated by go linker)
	// -ldflags="-X main.tag=${tag} -X main.commit=${commit}"
	tag    string
	commit string
)

func main() {
	Log_agg.Flags().StringVarP(&configFile, "config-file", "c", "", "config file location for server")

	config.AddFlags(Log_agg)

	err := Log_agg.Execute()
	if err != nil && err.Error() != "" {
		fmt.Println(err)
	}
}

func readConfig(ccmd *cobra.Command, args []string) error {
	if err := config.ReadConfigFile(configFile); err != nil {
		return err
	}
	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf("log_agg %s (%s)\n", tag, commit)
		return fmt.Errorf("")
	}

	return nil
}

func startLog_agg(ccmd *cobra.Command, args []string) error {
	// initialize logger
	lumber.Level(lumber.LvlInt(config.LogLevel)) // for clients using lumber too
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt(config.LogLevel))

	// initialize log_agg
	log_agg.Init()


	// initialize outputs
	err := output.Init()
	if err != nil {
		return fmt.Errorf("Output failed to initialize - %s", err)
	}

	// initializes inputs
	err = input.Init()
	if err != nil {
		return fmt.Errorf("Input failed to initialize - %s", err)
	}

	err = api.Start(input.InputHandler)
	if err != nil {
		return fmt.Errorf("Api failed to initialize - %s", err)
	}

	return nil
}


