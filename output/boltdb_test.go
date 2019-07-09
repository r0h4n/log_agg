// output_test tests the archiver output
// functionality (writing, reading, cleaning)
package output_test

import (
	"os"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
	"github.com/r0h4n/log_agg/output"
)

func TestMain(m *testing.M) {
	// clean test dir
	os.RemoveAll("/tmp/boltdbTest")

	// manually configure
	err := initialize()
	if err != nil {
		os.Exit(1)
	}

	rtn := m.Run()

	// clean test dir
	os.RemoveAll("/tmp/boltdbTest")

	os.Exit(rtn)
}

// Test writing and getting data
func TestWrite(t *testing.T) {
	// create test messages
	messages := []log_agg.Message{
		log_agg.Message{
			Time:     time.Now(),
			UTime:    time.Now().UnixNano(),
			Id:       "myhost",
			Tag:      []string{"test[bolt]"},
			Type:     "app",
			Priority: 4,
			Content:  "This is a test message",
		},
		log_agg.Message{
			Time:     time.Now(),
			UTime:    time.Now().UnixNano(),
			Id:       "myhost",
			Tag:      []string{"test[expire]"},
			Type:     "deploy",
			Priority: 4,
			Content:  "This is another test message",
		},
	}
	// write test messages
	output.Archiver.Write(messages[0])
	output.Archiver.Write(messages[1])

	// test successful write
	appMsgs, err := output.Archiver.Slice("app", "", []string{""}, 0, 0, 100, 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// compare written message to original
	if len(appMsgs) != 1 || appMsgs[0].Content != messages[0].Content {
		t.Errorf("%q doesn't match expected out", appMsgs)
		t.FailNow()
	}
}

// Test expiring/cleanup of data
func TestExpire(t *testing.T) {
	go output.Archiver.Expire()
	time.Sleep(2 * time.Second)

	// finish expire loop
	output.Archiver.(*output.BoltArchive).Done <- true

	// test successful clean
	appMsgs, err := output.Archiver.Slice("app", "", []string{""}, 0, 0, 100, 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// compare written message to original
	if len(appMsgs) != 0 {
		t.Errorf("%q doesn't match expected out", appMsgs)
		t.FailNow()
	}

	// test successful clean
	depMsgs, err := output.Archiver.Slice("deploy", "", []string{""}, 0, 0, 100, 0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// compare written message to original
	if len(depMsgs) != 0 {
		t.Errorf("%q doesn't match expected out", depMsgs)
		t.FailNow()
	}

	output.Archiver.(*output.BoltArchive).Close()

}

// manually configure and start internals
func initialize() error {
	var err error
	config.CleanFreq = 1
	config.LogKeep = `{"app": "1s", "deploy":0}`
	config.LogKeep = `{"app": "1s", "deploy":0, "a":"1m", "aa":"1h", "b":"1d", "c":"1w", "d":"1y", "e":"1"}`
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("ERROR"))

	// initialize log_agg
	log_agg.Init()

	// initialize archiver
	// Doing broke db
	config.DbAddress = "~!@#$%^&*()"
	output.Init()

	// Doing file db
	config.DbAddress = "file:///tmp/boltdbTest/log_agg.bolt"
	output.Init()
	output.Archiver.(*output.BoltArchive).Close()

	// Doing no db
	config.DbAddress = "/tmp/boltdbTest/log_agg.bolt"
	output.Init()
	output.Archiver.(*output.BoltArchive).Close()

	// Doing bolt db
	config.DbAddress = "boltdb:///tmp/boltdbTest/log_agg.bolt"
	output.Init()

	return err
}
