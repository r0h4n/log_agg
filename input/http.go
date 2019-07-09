package input

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
)

// GenerateHttpInput creates and returns an http handler that can be dropped into the api.
func GenerateHttpInput() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(500)
			return
		}

		var msg log_agg.Message
		err = json.Unmarshal(body, &msg)
		if err != nil {
			if !strings.Contains(err.Error(), "invalid character") {
				res.WriteHeader(500)
				res.Write([]byte(err.Error()))
				return
			}

			// keep body as "message" and make up priority
			msg.Content = string(body)
			msg.Priority = 2
			msg.Tag = []string{"http-raw"}
		}

		if msg.Type == "" {
			msg.Type = config.LogType
		}
		msg.Time = time.Now()
		msg.UTime = msg.Time.UnixNano()

		// config.Log.Trace("Message: %q", msg)
		log_agg.WriteMessage(msg)

		res.WriteHeader(200)
		res.Write([]byte("success!\n"))
	}
}
