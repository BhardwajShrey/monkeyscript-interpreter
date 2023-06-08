package main

import (
    "fmt"
    "os"
    "os/user"
    "monkey/repl"
)

func main() {
    user, err := user.Current()

    if err != nil {
        panic(err)
    }

    fmt.Printf("Hello %s. This is the Monkey Programming Language!\n", user.Username)
    fmt.Printf("Start typing away\n")
    repl.Start(os.Stdin, os.Stdout)
}
