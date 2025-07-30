package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/HubertasVin/findstr/utils"
	"github.com/spf13/pflag"
)

func main() {
	root, exdir, pattern, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		pflag.Usage()
		os.Exit(1)
	}

    if *exdir != "" {
        fmt.Println("Info: Exclude-dir flag is work in-progress, skipping flag for now.")
    }
	matches, err := utils.SearchMatchLines(*root, pattern)
	if err != nil {
		log.Fatal(err)
	}

	utils.PrintMatches(matches)
}

func parseFlags() (*string, *string, string, error) {
	exdir := pflag.StringP("exclude-dir", "e", "", "relative paths to ignore (comming soon)")
	root := pflag.StringP("root", "r", "./", "root directory to walk")

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: findstr [flags] <pattern>")
		fmt.Fprintln(os.Stderr, "Search for file content matching <pattern> under the given root.")
		fmt.Fprintln(os.Stderr)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if args := pflag.Args(); len(args) == 0 {
		return nil, nil, "", errors.New("you must provide a <pattern> to search for")
	} else {
		return root, exdir, args[0], nil
	}
}
