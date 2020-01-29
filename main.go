package main

import (
	"flag"
	"fmt"
	"monkey/bin"
	"monkey/evaluator"
	"monkey/repl"
	"os"
	user2 "os/user"
)

func init() {
	evaluator.InitBuiltins()
}

func main() {
	user, err := user2.Current()
	if err != nil {
		panic(err)
	}

	var code, path string
	flag.StringVar(&code, "c", "", "the code should run")
	flag.StringVar(&path, "f", "", "the source code file path")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	} else if code != "" {
		bin.RunCode(code)
	} else if path != "" {
		bin.RunFile(path)
	} else if len(os.Args) == 2 {
		bin.RunFile(os.Args[1])
	} else {
		fmt.Println("Parsing command line parameter failed!")
	}
}
