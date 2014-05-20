# loggly-go

  Loggly client for Go.

## Installation

    $ go get github.com/segmentio/loggly-go

## Example

```go
package main

import "github.com/segmentio/loggly-go"

func main() {
  log := loggly.New("8bad16f2-6c0e-4d90-944e-51238379f8d47")

  log.Error("boom")

  log.Info("connecting", loggly.Message{
    "some": "details",
    "here": 123,
  }})
}
```

## Options

  By default the client will flush every __100__ messages _or_ every __5__ seconds. A `.timestamp` property is also provided per log, and a map of overridable properties is provided, but defaults to only `.hostname`.

 - `.Debugging` (bool) enable debug logs
 - `.BufferSize` (int) size of the buffer [100]
 - `.FlushInterval` (time.Duration) flush interval [5 seconds]
 - `.Token` (string) loggly api token
 - `.Endpoint` (string) loggly api url
 - `.Defaults` (loggly.Message) default properties
 - `.Level` (loggly.Level) log level [loggly.Info]

## Levels

 Syslog level methods are provided, as well
 as a base `.Log()` call.

## License

 MIT
