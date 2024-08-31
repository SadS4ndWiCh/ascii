package main

import (
	"fmt"
	"os"

	"github.com/SadS4ndWiCh/ascii/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ascii <image | video> -i <file> [-w <width>] [-h <height>]")
		os.Exit(1)
	}

	if err := commands.Root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
