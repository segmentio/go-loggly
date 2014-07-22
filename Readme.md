# go-loggly

  Loggly client for Go.

  View the [docs](http://godoc.org/github.com/segmentio/go-loggly).

## Installation

    $ go get github.com/segmentio/go-loggly

## Example

```go
package main

import "github.com/segmentio/go-loggly"
import "time"
import "os"

func main() {
  client := loggly.New("api-token-here-whoop")
  client.Writer = os.Stderr

  for {
    client.Info("something here")
    time.Sleep(15 * time.Millisecond)
  }
}
```

## Debug

 Enable verbose debugging output using the __DEBUG__ environment variable, for exmaple `DEBUG=loggly`.

## License

 MIT
