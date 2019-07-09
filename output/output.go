// Package output handles the storing of logs.
package output

import (
	"fmt"
	"net/url"
	"github.com/r0h4n/log_agg/config"
	"github.com/r0h4n/log_agg/transform"
)

type (
	// defines a storage type output
	Output interface {
		// Init initializes the output interface
		Init() error
		// Slice returns a slice of logs based on the name, offset, limit, and log-level
		Slice(name, host string, tag []string, offset, end, limit int64, level int) ([]log_agg.Message, error)
		// Write writes the message to file/database
		Write(msg log_agg.Message)
		// Expire cleans up old logs
		Expire()
	}

)

var Archiver Output             // default archive output

// Init initializes the archiver output if configured
func Init() error {
	// initialize archiver
	err := archiveInit()
	if err != nil {
		return fmt.Errorf("Failed to initialize archiver - %s", err)
	}
	config.Log.Info("Archiving output '%s' initialized", config.DbAddress)

	return nil
}

func archiveInit() error {
	u, err := url.Parse(config.DbAddress)
	if err != nil {
		u, err = url.Parse("boltdb://" + config.DbAddress)
		if err != nil {
			return fmt.Errorf("Failed to parse db connection - %s", err)
		}
	}


	switch u.Scheme {
	case "boltdb":
		Archiver, err = NewBoltArchive(u.Path)
		if err != nil {
			return err
		}
	case "file":
		Archiver, err = NewBoltArchive(u.Path)
		if err != nil {
			return err
		}
	default:
		Archiver, err = NewBoltArchive(u.Path)
		if err != nil {
			return err
		}
	}
	// initialize Archiver
	err = Archiver.Init()
	if err != nil {
		return err
	}
	// start cleanup goroutine
	go Archiver.Expire()
	return nil
}
