package internal

import (
    "fmt"
    "os"
)

func Fatal(msg string) {
    fmt.Fprintln(os.Stderr, msg)
    os.Exit(1)
}

func FatalIf(err error) {
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}