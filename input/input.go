// Package input initializes http servers for collecting logs.
package input

import (
	"net/http"

	"github.com/r0h4n/log_agg/config"
)

var (
	// InputHandler handles the posting of logs via http. It is passed to
	// the api on start.
	InputHandler http.HandlerFunc
)

// Init initializes the http server, if configured
func Init() error {

	if config.ListenHttp != "" {
		InputHandler = GenerateHttpInput()
		config.Log.Info("Input listening on http://%s...", config.ListenHttp)
	}

	return nil
}
