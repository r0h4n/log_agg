// input_test tests the syslog inputs
// (http input is tested in api_test)
package input_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/r0h4n/log_agg/api"
	"github.com/r0h4n/log_agg/input"
	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
	"github.com/r0h4n/log_agg/output"
)

func TestMain(m *testing.M) {
	// clean test dir
	os.RemoveAll("/tmp/syslogTest")

	// manually configure
	initialize()

	// start api
	go api.Start(input.InputHandler)
	<-time.After(1 * time.Second)
	rtn := m.Run()

	// clean test dir
	os.RemoveAll("/tmp/syslogTest")

	os.Exit(rtn)
}

// test post logs
func TestPostLogs(t *testing.T) {
	body, err := rest("POST", "/logs", "{\"id\":\"log-test\",\"type\":\"app\",\"message\":\"test log\"}")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(body) != "success!\n" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
	// pause for travis
	time.Sleep(500 * time.Millisecond)
	body, err = rest("POST", "/logs", "another test log")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(body) != "success!\n" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
	// boltdb seems to take some time committing the record (probably the speed/immediate commit tradeoff)
	time.Sleep(time.Second)
}


// test get logs
func TestGetLogs(t *testing.T) {
	body, err := rest("GET", "/logs?type=app", "")
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

	if len(msg) != 2 || msg[0].Content != "test log" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
	if msg[1].Content != "another test log" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
}

// hit api and return response body
func rest(method, route, data string) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))

	req, _ := http.NewRequest(method, fmt.Sprintf("http://%s%s", config.ListenHttp, route), body)

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
	config.ListenHttp = "0.0.0.0:4234"
	config.DbAddress = "boltdb:///tmp/syslogTest/log_agg.bolt"
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
