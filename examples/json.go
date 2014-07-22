package main

import "github.com/segmentio/go-loggly"
import "time"
import "os"

func main() {
	client := loggly.New("8bad16f2-6c0e-4d90-944e-5668779f8d47")
	client.Writer = os.Stderr

	for {
		client.Info("connecting")

		client.Critical("failed to connect", map[string]interface{}{
			"host": "some-address-here",
			"port": "some-port-here",
		})

		time.Sleep(50 * time.Millisecond)
	}
}
