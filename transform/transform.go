// Package log_agg handles the adding, removing, and writing to outputs. It also
// defines the common types used accross log_agg.
package log_agg

import (
	"io"
	"sync"
	"time"

	"github.com/r0h4n/log_agg/config"
)

type (
	// Logger is a simple interface that's designed to be intentionally generic to
	// allow many different types of Logger's to satisfy its interface
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	// Message defines the structure of a log message
	Message struct {
		Time     time.Time `json:"time"`
		UTime    int64     `json:"utime"`
		Id       string    `json:"id"`   // ignoreifempty? // If setting multiple tags in id (syslog), set hostname first
		Tag      []string  `json:"tag"`  // ignoreifempty?
		Type     string    `json:"type"` // Can be set if logs are submitted via http (deploy logs)
		Priority int       `json:"priority"`
		Content  string    `json:"message"`
		Raw      []byte    `json:"raw,omitempty"`
	}

	// Log_agg defines the structure for the default log_agg object
	Log_agg struct {
		outputs map[string]outputChannels
	}

	// Output defines a third party log output endpoint (generally, only raw logs get outputed)
	Output struct {
		Type       string `json:"type"`             // type of service ("papertrail")
		URI        string `json:"endpoint"`         // uri of endpoint "log6.papertrailapp.com:199900"
		ID				 string	`json:"id"`								// id to identify this app with external logger
		AuthKey    string `json:"key,omitempty"`    // key or user for authentication
		AuthSecret string `json:"secret,omitempty"` // password or secret for authentication
	}

	// OutputFunc is a function that "outputs a Message"
	OutputFunc func(Message)

	outputChannels struct {
		send chan Message
		done chan bool
	}
)

// Vac is the default log_agg object
var Vac Log_agg

// Initializes a log_agg object
func Init() error {
	Vac = Log_agg{
		outputs: make(map[string]outputChannels),
	}
	config.Log.Debug("Log_agg initialized")
	return nil
}

// Close log_agg and remove all outputs
func Close() {
	Vac.close()
}

func (l *Log_agg) close() {
	for tag := range l.outputs {
		l.removeOutput(tag)
	}
}

// AddOutput adds a output to the listeners and sets its logger
func AddOutput(tag string, output OutputFunc) {
	Vac.addOutput(tag, output)
}

func (l *Log_agg) addOutput(tag string, output OutputFunc) {
	channels := outputChannels{
		done: make(chan bool),
		send: make(chan Message),
	}

	go func() {
		for {
			select {
			case <-channels.done:
				return
			case msg := <-channels.send:
				// don't goroutine to preserve log order
				output(msg)
			}
		}
	}()

	l.outputs[tag] = channels
}

// RemoveOutput drops a output
func RemoveOutput(tag string) {
	Vac.removeOutput(tag)
}

func (l *Log_agg) removeOutput(tag string) {
	_, ok := l.outputs[tag]
	if ok {
		close(l.outputs[tag].done)
		delete(l.outputs, tag)
	}
}

// WriteMessage broadcasts to all outputs in seperate go routines
// Returns once all outputs have received the message, but may not have processed
// the message yet
func WriteMessage(msg Message) {
	Vac.writeMessage(msg)
}

func (l *Log_agg) writeMessage(msg Message) {
	// config.Log.Trace("Writing message - %s...", msg)
	group := sync.WaitGroup{}
	for _, output := range l.outputs {
		group.Add(1)
		go func(myOutput outputChannels) {
			select {
			case <-myOutput.done:
			case myOutput.send <- msg:
			}
			group.Done()
		}(output)
	}
	group.Wait()
}

func (m Message) eof() bool {
	return len(m.Raw) == 0
}

func (m *Message) readByte() byte {
	// this function assumes that eof() check was done before
	b := m.Raw[0]
	m.Raw = m.Raw[1:]
	return b
}

func (m *Message) Read(p []byte) (n int, err error) {
	if m.eof() {
		err = io.EOF
		return
	}

	if c := cap(p); c > 0 {
		for n < c {
			p[n] = m.readByte()
			n++
			if m.eof() {
				break
			}
		}
	}
	return
}

