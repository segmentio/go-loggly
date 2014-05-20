package loggly

//
// dependencies
//

import . "encoding/json"
import "io/ioutil"
import "net/http"
import "strings"
import "bytes"
import "time"
import "log"
import "os"

//
// Library version
//

const Version = "0.0.2"

//
// Default API end-point
//

const api = "https://logs-01.loggly.com/bulk/{token}/tag/bulk"

//
// Message type
//

type Message map[string]interface{}

//
// Newline delim
//

var nl = []byte{'\n'}

//
// Log levels
//

type Level int

const (
	Debug Level = iota
	Info
	Notice
	Warning
	Error
	Critical
	Alert
	Emergency
)

//
// Loggly client
//

type Client struct {
	Level         Level
	Debugging     bool
	BufferSize    int
	FlushInterval time.Duration
	Endpoint      string
	Token         string
	buffer        [][]byte
	Defaults      Message
}

//
// Return a new Loggly client with the given token.
//

func New(token string) (c *Client) {
	host, err := os.Hostname()

	if err != nil {
		log.Fatal("failed to get hostname")
	}

	defaults := Message{
		"hostname": host,
	}

	defer func() {
		go func() {
			for {
				time.Sleep(c.FlushInterval)
				c.log("interval %v reached", c.FlushInterval)
				go c.flush()
			}
		}()
	}()

	return &Client{
		Level:         Info,
		Debugging:     false,
		BufferSize:    100,
		FlushInterval: 5 * time.Second,
		Token:         token,
		Endpoint:      strings.Replace(api, "{token}", token, 1),
		buffer:        make([][]byte, 0),
		Defaults:      defaults,
	}
}

//
// Buffer a log message.
//

func (c *Client) Send(msg Message) error {
	msg["timestamp"] = int32(time.Now().Unix())
	merge(msg, c.Defaults)

	json, err := Marshal(msg)

	if err != nil {
		return err
	}

	c.buffer = append(c.buffer, json)

	c.log("buffer (%d/%d) %v", len(c.buffer), c.BufferSize, msg)

	if len(c.buffer) >= c.BufferSize {
		go c.flush()
	}

	return nil
}

//
// Debug log.
//

func (c *Client) Debug(t string, props ...Message) error {
	if c.Level > Debug {
		return nil
	}
	msg := Message{"level": "debug", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Info log.
//

func (c *Client) Info(t string, props ...Message) error {
	if c.Level > Info {
		return nil
	}
	msg := Message{"level": "info", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Notice log.
//

func (c *Client) Notice(t string, props ...Message) error {
	if c.Level > Notice {
		return nil
	}
	msg := Message{"level": "notice", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Warning log.
//

func (c *Client) Warn(t string, props ...Message) error {
	if c.Level > Warning {
		return nil
	}
	msg := Message{"level": "warning", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Error log.
//

func (c *Client) Error(t string, props ...Message) error {
	if c.Level > Error {
		return nil
	}
	msg := Message{"level": "error", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Critical log.
//

func (c *Client) Critical(t string, props ...Message) error {
	if c.Level > Critical {
		return nil
	}
	msg := Message{"level": "critical", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Alert log.
//

func (c *Client) Alert(t string, props ...Message) error {
	if c.Level > Alert {
		return nil
	}
	msg := Message{"level": "alert", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Emergency log.
//

func (c *Client) Emergency(t string, props ...Message) error {
	if c.Level > Emergency {
		return nil
	}
	msg := Message{"level": "emergency", "type": t}
	merge(msg, props...)
	return c.Send(msg)
}

//
// Log in debug mode.
//

func (c *Client) log(format string, v ...interface{}) {
	if c.Debugging {
		log.Printf(format, v...)
	}
}

//
// Merge others into a.
//

func merge(a Message, others ...Message) {
	for _, msg := range others {
		for k, v := range msg {
			a[k] = v
		}
	}
}

//
// Flush the buffered messages.
//

func (c *Client) flush() error {
	if len(c.buffer) == 0 {
		c.log("no messages to flush")
		return nil
	}

	c.log("flushing %d messages", len(c.buffer))
	body := bytes.Join(c.buffer, nl)

	c.buffer = nil

	client := &http.Client{}
	c.log("POST %s with %d bytes", c.Endpoint, len(body))
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(body))

	if err != nil {
		c.log("error: %v", err)
		return err
	}

	req.Header.Add("User-Agent", "loggly-go (version: "+Version+")")
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Content-Length", string(len(body)))

	res, err := client.Do(req)

	if err != nil {
		c.log("error: %v", err)
		return err
	}

	c.log("%d response", res.StatusCode)

	if res.StatusCode >= 400 {
		resp, _ := ioutil.ReadAll(res.Body)
		c.log("error: %s", string(resp))
	}

	return err
}
