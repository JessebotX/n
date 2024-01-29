package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	UserConfigDirBasename = "n"
	LocalConfigFileBasename = "n.yml"
	UserConfigFileBasename = "config.yml"
	UsageInfo = `
USAGE
=====
n <command> [options...]
`
	Version = "0.1.0"
)

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
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		exitWithMessage(1, err.Error())
	}

	// parse opts
	opts := make(map[string]string)
	nonOpts := make([]string, 0)
	for i, arg := range args {
		if (arg == "-e" || arg == "--embed") && (len(args)-1) > i {
			opts["-e"] = args[i+1]
			continue
		} else if (arg == "-r" || arg == "--reference") && (len(args)-1) > i  {
			opts["-r"] = args[i+1]
		} else if (arg == "-d" || arg == "--directory") && (len(args)-1) > i {
			opts["-d"] = args[i+1]
		} else {
			if isNonOptArg(arg, opts) {
				nonOpts = append(nonOpts, arg)
			}
		}
	}

	// get config
	config := UserConfig{
		Editor:          "vim",
		DefaultNotesDir: "~/Documents/notes",
	}
	if len(opts["-d"]) > 0 {
		config.DefaultNotesDir = opts["-d"]
	}

	configPath := filepath.Join(userConfigDir, UserConfigDirBasename, UserConfigFileBasename)
	err = UnmarshalConfig(&config, configPath)
	if err != nil {
		exitWithMessage(1, err.Error())
	}

	fmt.Printf("CONFIG: %v\n", config)
	fmt.Printf("OPTS: %v\n", opts)
	fmt.Printf("ARGS: %v\n", nonOpts)

	//fmt.Println("TODO Create new command.")
	entry, err := getNewNoteEntryDir(config.DefaultNotesDir)
	if err != nil {
		exitWithMessage(1, err.Error())
	}

	fmt.Println(entry)
}

func getNewNoteEntryDir(notesDir string) (string, error) {
	i := 1
	entryDir := filepath.Join(notesDir, strconv.Itoa(i))
	for {
		_, err := os.Stat(entryDir)

		if os.IsNotExist(err) {
			break
		}

		i++
		entryDir = filepath.Join(notesDir, strconv.Itoa(i))
	}

	return entryDir, nil
}

func isNonOptArg(arg string, currentOpts map[string]string) bool {
	for _, v := range currentOpts {
		if v == arg {
			return false
		}
	}

	return true
}

func exitWithMessage(exitCode int, message string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: " + message + "\n", args...)

	os.Exit(exitCode)
}
