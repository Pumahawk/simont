package main

import (
	"flag"
	"fmt"
)

var Commands = []*Command{
	LsCommand,
	InspectCommand,
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		for _, c := range Commands {
			if args[0] == c.Name {
				c.Run(c, args[1:])
				return
			}
		}
	}
	PrintAllCommands()
}

func PrintAllCommands() {
	for _, c := range Commands {
		fmt.Printf("%s - %s\n", c.Name, c.Descr)
	}
}
