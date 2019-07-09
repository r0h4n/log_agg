// Package api handles the api routes and related funtionality.
//
// ROUTES 
//
// | Action | Route | Description       | Payload                          | Output          |
// |--------|-------|-------------------|----------------------------------|-----------------|
// | POST   | /logs | Publish a log     | Log Message                      | Success message |
// | GET    | /logs | Fetch stored logs |                                  | Success message |
//
package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/jcelliott/lumber"
	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/output"
)

// starts the web server with the log_agg functions
func Start(input http.HandlerFunc) error {
	retriever := GenerateArchiveEndpoint(output.Archiver)

	router := pat.New()

	router.Post("/logs", handleRequest(input))
	router.Get("/logs", handleRequest(retriever))

	httpListener, err := net.Listen("tcp", config.ListenHttp)
	if err != nil {
		return err
	}

	config.Log.Info("Api Listening on http://%s...", config.ListenHttp)
	return http.Serve(httpListener, router)

}


// adds a bit of logging
func handleRequest(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", config.CorsAllow)
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST")

		fn(rw, req)

		// must be after req returns
		getStatus := func(trw http.ResponseWriter) string {
			r, _ := regexp.Compile("status:([0-9]*)")
			return r.FindStringSubmatch(fmt.Sprintf("%+v", trw))[1]
		}

		getWrote := func(trw http.ResponseWriter) string {
			r, _ := regexp.Compile("written:([0-9]*)")
			return r.FindStringSubmatch(fmt.Sprintf("%+v", trw))[1]
		}

		config.Log.Debug(`%s - [%s] %s %s %s(%s) - "User-Agent: %s"`,
			req.RemoteAddr, req.Proto, req.Method, req.RequestURI,
			getStatus(rw), getWrote(rw), // %s(%s)
			req.Header.Get("User-Agent"))
	}
}


// generates the endpoint for fetching filtered logs
func GenerateArchiveEndpoint(archive output.Output) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// /logs?id=&type=app&start=0&end=0&limit=50
		query := req.URL.Query()

		host := query.Get("id")
		tag := query["tag"]

		kind := query.Get("type")
		if kind == "" {
			kind = config.LogType // "app"
		}
		start := query.Get("start")
		if start == "" {
			start = "0"
		}
		end := query.Get("end")
		if end == "" {
			end = "0"
		}
		limit := query.Get("limit")
		if limit == "" {
			limit = "100"
		}
		level := query.Get("level")
		if level == "" {
			level = "TRACE"
		}
		config.Log.Trace("type: %s, start: %s, end: %s, limit: %s, level: %s, id: %s, tag: %s", kind, start, end, limit, level, host, tag)
		logLevel := lumber.LvlInt(level)
		realOffset, err := strconv.ParseInt(start, 0, 64)
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte("bad start offset"))
			return
		}
		realEnd, err := strconv.ParseInt(end, 0, 64)
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte("bad end value"))
			return
		}
		realLimit, err := strconv.Atoi(limit)
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte("bad limit"))
			return
		}
		slices, err := archive.Slice(kind, host, tag, realOffset, realEnd, int64(realLimit), logLevel)
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
			return
		}
		body, err := json.Marshal(slices)
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
			return
		}

		res.WriteHeader(200)
		res.Write(append(body, byte('\n')))
	}
}

// parses the request into v
func parseBody(req *http.Request, v interface{}) error {

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	config.Log.Trace("Parsed body - %s", b)

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}
