// api_test tests the api, from posting logs, to getting them
package api_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/r0h4n/log_agg/api"
	"github.com/r0h4n/log_agg/input"
	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
	"github.com/r0h4n/log_agg/output"
)

var insecureHttp string


func TestMain(m *testing.M) {
	// clean test dir
	os.RemoveAll("/tmp/apiTest")

	// manually configure
	initialize()

	// start insecure api
	go api.Start(input.InputHandler)
	time.Sleep(time.Second)
	<-time.After(time.Second)
	rtn := m.Run()

	// clean test dir
	os.RemoveAll("/tmp/apiTest")

	os.Exit(rtn)
}


// test post logs
func TestPostLogs(t *testing.T) {
	// secure
	body, err := rest("POST", "/logs", "{\"id\":\"log-test\",\"type\":\"app\",\"message\":\"test log\"}")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(body) != "success!\n" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
	// insecure
	body, err = rest("POST", "/logs", "{\"id\":\"log-test\",\"type\":\"app\",\"message\":\"test log\"}")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(body) != "success!\n" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
	// boltdb seems to take some time committing the record (probably the speed/immediate commit tradeoff)
	time.Sleep(500 * time.Millisecond)
}

// test get logs
func TestGetLogs(t *testing.T) {
	body, err := rest("GET", "/logs?type=app&id=log-test&start=0&limit=1", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	msg := []log_agg.Message{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		t.Error(fmt.Errorf("Failed to unmarshal - %s", err))
		t.FailNow()
	}
	if len(msg) != 1 || msg[0].Content != "test log" {
		t.Errorf("%q doesn't match expected out", body)
	}
	_, err = rest("GET", "/logs", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = rest("GET", "/logs?start=word", "")
	if err == nil || strings.Contains(err.Error(), "bad start offset") {
		t.Error("bad start is too forgiving")
		t.FailNow()
	}
	_, err = rest("GET", "/logs?limit=word", "")
	if err == nil || strings.Contains(err.Error(), "bad limit") {
		t.Error("bad limit is too forgiving")
		t.FailNow()
	}
	_, err = rest("GET", "/logs?level=word", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}


// hit api and return response body
func rest(method, route, data string) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))

	req, _ := http.NewRequest(method, fmt.Sprintf("http://%s%s", insecureHttp, route), body)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to %s %s - %s", method, route, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status '200' expected, got '%d'", res.StatusCode)
	}

	b, _ := ioutil.ReadAll(res.Body)

	return b, nil
}



// manually configure and start internals
func initialize() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	insecureHttp = "0.0.0.0:2234"
	config.ListenHttp = "0.0.0.0:2234"
	config.DbAddress = "boltdb:///tmp/apiTest/log_agg.bolt"
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("ERROR"))

	// initialize log_agg
	log_agg.Init()

	// initialize outputs
	err := output.Init()
	if err != nil {
		config.Log.Fatal("Output failed to initialize - %s", err)
		os.Exit(1)
	}

	// initializes inputs
	err = input.Init()
	if err != nil {
		config.Log.Fatal("Input failed to initialize - %s", err)
		os.Exit(1)
	}
}
