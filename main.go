// main.go

package main

import (
	"fmt"
	"mira/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Mira REPL!\n", user.Username)
	fmt.Printf("You can get started by typing some commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}
