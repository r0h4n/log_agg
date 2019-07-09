// log_agg_test tests the transform output functionality (adding, writing, removing)
package log_agg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
)

func TestMain(m *testing.M) {
	// manually configure
	err := initialize()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// Test adding and writing to a output
func TestAddOutput(t *testing.T) {
	// create a buffer to "output" to
	buf := &bytes.Buffer{}

	// add buffer output
	log_agg.AddOutput("test", writeOutput(buf))

	// create test message
	msg := log_agg.Message{
		Time:     time.Now(),
		UTime:    time.Now().UnixNano(),
		Id:       "myhost",
		Tag:      []string{"test[outputs]"},
		Type:     "app",
		Priority: 4,
		Content:  "This is quite important",
	}

	// write test message
	log_agg.WriteMessage(msg)
	time.Sleep(time.Millisecond)

	// ensure write succeeded
	// read from buffer
	r, err := buf.ReadBytes('\n')
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// convert readbytes to Message
	rMsg := log_agg.Message{}
	err = json.Unmarshal(r, &rMsg)
	if err != nil {
		t.Error(fmt.Errorf("Failed to unmarshal - %s", err))
		t.FailNow()
	}

	// compare "outputed" message to original
	if rMsg.Content != msg.Content {
		t.Errorf("%q doesn't match expected out", rMsg.Content)
		t.FailNow()
	}
}

// Test removing a output
func TestRemoveOutput(t *testing.T) {
	tag := "null"
	output := func(msg log_agg.Message) {
		return
	}
	log_agg.AddOutput(tag, output)
	log_agg.RemoveOutput(tag)
}

// Test closing the log_agg instance
func TestClose(t *testing.T) {
	log_agg.Close()
	time.Sleep(time.Second)
}

// writeOutput creates a output from an io.Writer
func writeOutput(writer io.Writer) log_agg.OutputFunc {
	return func(msg log_agg.Message) {
		data, err := json.Marshal(msg)
		if err != nil {
			config.Log.Error("writeOutput failed to marshal message")
			return
		}
		writer.Write(append(data, '\n'))
	}
}

// manually configure and start internals
func initialize() error {
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("ERROR"))

	// initialize log_agg
	return log_agg.Init()
}
