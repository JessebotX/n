package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	UserConfigDirBasename   = "n"
	LocalConfigFileBasename = "n.yml"
	UserConfigFileBasename  = "config.yml"
	UsageInfo               = `
USAGE
=====
n <command> [options...]
`
	Version = "0.1.0"
)

var DefaultConfig = UserConfig{
	NoteEntryIndexFileName: "README.org",
	Editor:                 "vim",
	DefaultNotesDir:        "~/Documents/notes",
}

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
		} else if (arg == "-r" || arg == "--reference") && (len(args)-1) > i {
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
	configPath := filepath.Join(userConfigDir, UserConfigDirBasename, UserConfigFileBasename)
	config := DefaultConfig
	if len(opts["-d"]) > 0 {
		config.DefaultNotesDir = opts["-d"]
	}

	err = UnmarshalConfig(&config, configPath)
	if err != nil {
		exitWithMessage(1, err.Error())
	}

	fmt.Printf("CONFIG: %v\n", config)
	fmt.Printf("OPTS: %v\n", opts)
	fmt.Printf("ARGS: %v\n", nonOpts)

	// get a new note entry dir path
	entryDir, err := getNewNoteEntryDir(config.DefaultNotesDir)
	if err != nil {
		exitWithMessage(1, err.Error())
	}

	// create dir and create index file
	err = os.MkdirAll(entryDir, os.ModePerm)
	if err != nil {
		exitWithMessage(1, err.Error())
	}

	entryPath := filepath.Join(entryDir, config.NoteEntryIndexFileName)
	fmt.Println("CREATING", entryPath)

	entry, err := os.Create(entryPath)
	if err != nil {
		exitWithMessage(1, err.Error())
	}
	defer entry.Close()

	metadata := ""
	content := ""
	title := "No title provided..."
	filetags := []string{}
	if len(nonOpts) > 0 {
		title = nonOpts[0]
	}
	metadata = fmt.Sprintf("#+title: %s\n", title)

	// check if a reference link has been provided
	_, optsRefExist := opts["-r"]
	if optsRefExist {
		metadata = metadata + fmt.Sprintf("#+ref: %s\n", opts["-r"])
		filetags = append(filetags, "ref")
		content = content + fmt.Sprintf("<%s>\n", opts["-r"])
	}

	// check if an embed object is provided
	_, optsEmbedExist := opts["-e"]
	if optsEmbedExist {
		embedBaseFileName := filepath.Base(opts["-e"])
		copyDestination := filepath.Join(entryDir, embedBaseFileName)
		err = os.Link(opts["-e"], copyDestination)
		if err != nil {
			exitWithMessage(1, err.Error())
		}
		filetags = append(filetags, "embed")

		metadata = metadata + fmt.Sprintf("#+embed: ./%s\n", embedBaseFileName)
		content = content + fmt.Sprintf("\n* Attachment\n[[./%s]]\n", embedBaseFileName)
	}

	// add filetags
	if len(filetags) > 0 {
		metadata = metadata + "#+filetags: "
		for _, tag := range filetags {
			metadata = metadata + ":" + tag
		}
		metadata = metadata + ":\n"
	}

	// write to the file
	entry.Write([]byte(metadata + "\n" + content + "\n"))
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
	fmt.Fprintf(os.Stderr, "ERROR: "+message+"\n", args...)

	os.Exit(exitCode)
}
