
package main

import (
    "fmt"
    "log"
    "github.com/segmentio/go-loggly"
)

func main(){


    logToLoggly := loggly.New("YOUR_LOGGLY_TOKEN_HERE")

    err := logToLoggly.Info("A test message from Golang to Loggly!")

    if err != nil {
        log.Fatal(err)
    }

    logToLoggly.Flush()

    fmt.Println("done.")

}
