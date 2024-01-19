package main

import (
	"fmt"
	"os"
)

const Version = "0.1.0"
const UsageInfo = `
USAGE
=====
n <command> [options...]
`

func main() {
	if len(os.Args) < 2 {
		exitWithMessage(1, "Missing arguments. See command 'help' for more information.")
	}

	command := os.Args[1]

	if command == "help" {
		fmt.Println(UsageInfo)
	} else if command == "version" {
		fmt.Println("n v" + Version)
	} else if command == "new" {
		if len(os.Args) >= 3 {
			commandNew(os.Args[2:])
		} else {
			commandNew([]string{})
		}
	} else {
		exitWithMessage(1, "Command '%s' does not exist. See command 'help' for more information.", command)
	}
}

func commandNew(args []string) {
	_ = ReadConfig()

	fmt.Println("Create new command.")
}

func exitWithMessage(exitCode int, message string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: " + message + "\n", args...)

	os.Exit(exitCode)
}
