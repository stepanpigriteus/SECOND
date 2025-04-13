package pkg

import (
	"flag"
	"fmt"
	"os"
)

func Flags() int {
	help := flag.Bool("help", false, "help panel")
	port := flag.Int("port", 8083, "Port number")
	flag.Parse()

	if !ValidatePort(*port) {
		*port = 8083
	}

	if *help {
		showHelp()
		os.Exit(0)
	}
	return *port
}

func showHelp() {
	helpText := `

Usage:
 1337b04rd [--port <N>] 
 1337b04rd --help

Options:
 --help Show this screen.
 --port N Port number.
`
	fmt.Println(helpText)
}

func ValidatePort(port int) bool {
	if port == 5432 || port == 5073 || port < 1000 || port > 9999 {
		return false
	}
	return true
}
