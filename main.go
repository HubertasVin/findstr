package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/HubertasVin/findstr/models"
	"github.com/HubertasVin/findstr/utils"
	"github.com/spf13/pflag"
)

func main() {
	exdir, exfile, threadc, root, pattern, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		pflag.Usage()
		os.Exit(1)
	}

	flags := models.ProgramFlags{
		ExcludeDir:  *exdir,
		ExcludeFile: *exfile,
		ThreadCount: threadc,
		Root:        *root,
		Pattern:     pattern,
	}

	// TODO: Implement exclude-file-regex flag
	// TODO: Parallelise the program
	if *exfile != "" {
	    fmt.Println("Info: Exclude-file flag is work in-progress, skipping flag for now.")
	}
	if threadc <= 0 {
        log.Panicln("Error: You must select a thread count that is greater than 0.")
	}
	matches, err := utils.SearchMatchLines(flags)
	if err != nil {
		log.Fatal(err)
	}

	utils.PrintMatches(matches)
}

func parseFlags() (*string, *string, int, *string, string, error) {
	exdir := pflag.StringP("exclude-dir", "e", "", "relative paths to ignore (comming soon)")
	exfile := pflag.StringP(
		"exclude-file",
		"x",
		"",
		"regex match which files to ignore (comming soon)",
	)
	threadc := pflag.IntP(
		"thread-count",
		"t",
		1,
		"thread count to use for file parsing (comming soon)",
	)
	root := pflag.StringP("root", "r", "./", "root directory to walk")

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: findstr [flags] <pattern>")
		fmt.Fprintln(os.Stderr, "Search for file content matching <pattern> under the given root.")
		fmt.Fprintln(os.Stderr)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if args := pflag.Args(); len(args) == 0 {
		return nil, nil, 0, nil, "", errors.New("you must provide a <pattern> to search for")
	} else {
		return exdir, exfile, *threadc, root, args[0], nil
	}
}
