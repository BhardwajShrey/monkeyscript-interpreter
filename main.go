package main

import (
    "fmt"
    "os"
    "os/user"
    "monkey/repl"
)

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func main() {
    user, err := user.Current()

    if err != nil {
        panic(err)
    }

    fmt.Printf(MONKEY_FACE)
    fmt.Printf("Hello %s. This is the Monkey Programming Language!\n", user.Username)
    fmt.Printf("Start typing away\n")
    repl.Start(os.Stdin, os.Stdout)
}
