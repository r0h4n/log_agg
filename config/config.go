// location for configuration options, also contains config file parsing logic.
package config

import (
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// inputs
	ListenHttp = "0.0.0.0:6360" // address the api and http log inputs listen on

	// outputs
	DbAddress  = "boltdb:///var/db/log_agg.bolt" // database address

	// other
	CorsAllow = "*"            // sets `Access-Control-Allow-Origin` header
	LogKeep   = `{"app":"2w"}` // LogType and expire (X(m)in, (h)our,  (d)ay, (w)eek, (y)ear) (1, 10, 100 == keep up to that many) // todo: maybe map[string]interface
	LogType   = "app"          // default incoming log type when not set
	LogLevel  = "info"         // level which log_agg will log at
	Log       lumber.Logger    // logger to write logs
	Version   = false          // whether or not to print version info and exit
	CleanFreq = 60             // how often to clean log database
)

// AddFlags adds cli flags to log_agg
func AddFlags(cmd *cobra.Command) {
	// inputs
	cmd.Flags().StringVarP(&ListenHttp, "listen-http", "a", ListenHttp, "API listen address (same endpoint for http log collection)")

	// outputs
	cmd.Flags().StringVarP(&DbAddress, "db-address", "d", DbAddress, "Log storage address")

	// other
	cmd.Flags().StringVarP(&CorsAllow, "cors-allow", "C", CorsAllow, "Sets the 'Access-Control-Allow-Origin' header")
	cmd.Flags().StringVarP(&LogKeep, "log-keep", "k", LogKeep, "Age or number of logs to keep per type '{\"app\":\"2w\", \"deploy\": 10}' (int or X(m)in, (h)our,  (d)ay, (w)eek, (y)ear)")
	cmd.Flags().StringVarP(&LogLevel, "log-level", "l", LogLevel, "Level at which to log")
	cmd.Flags().StringVarP(&LogType, "log-type", "L", LogType, "Default type to apply to incoming logs (commonly used: app|deploy)")
	cmd.Flags().BoolVarP(&Version, "version", "v", Version, "Print version info and exit")
	cmd.Flags().IntVar(&CleanFreq, "clean-frequency", CleanFreq, "How often to clean log database")
	cmd.Flags().MarkHidden("clean-frequency")

	Log = lumber.NewConsoleLogger(lumber.LvlInt("ERROR"))
}

// ReadConfigFile reads in the config file, if any
func ReadConfigFile(configFile string) error {
	if configFile == "" {
		return nil
	}

	// Set defaults to whatever might be there already
	viper.SetDefault("listen-http", ListenHttp)
	viper.SetDefault("db-address", DbAddress)
	viper.SetDefault("cors-allow", CorsAllow)
	viper.SetDefault("log-keep", LogKeep)
	viper.SetDefault("log-level", LogLevel)
	viper.SetDefault("log-type", LogType)

	filename := filepath.Base(configFile)
	viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
	viper.AddConfigPath(filepath.Dir(configFile))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Set values. Config file will override commandline
	ListenHttp = viper.GetString("listen-http")
	DbAddress = viper.GetString("db-address")
	CorsAllow = viper.GetString("cors-allow")
	LogKeep = viper.GetString("log-keep")
	LogLevel = viper.GetString("log-level")
	LogType = viper.GetString("log-type")

	return nil
}
