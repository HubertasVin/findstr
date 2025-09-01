package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/HubertasVin/findstr/mappers"
	"github.com/HubertasVin/findstr/models"
	"github.com/HubertasVin/findstr/utils"
	"github.com/spf13/pflag"
)

func main() {
	showVersion, exdir, exfile, threadc, context, root, config, pattern, jsonOut, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		pflag.Usage()
		os.Exit(1)
	}

	if showVersion {
		printVersion()
		return
	}

	flags := models.ProgramFlags{
		ExcludeDir:  *exdir,
		ExcludeFile: *exfile,
		ThreadCount: threadc,
		ContextSize: context,
		Root:        *root,
		Json:        jsonOut,
		Pattern:     pattern,
	}

	if threadc <= 0 {
		fmt.Println("Error: Thread count must be greater than 0")
		os.Exit(1)
	}
	if context < 0 {
		fmt.Println("Error: Context size must be greater than or equal to 0")
		os.Exit(1)
	}

	cl, matchStyle, err := utils.ParseConfig(*config)
	if err != nil {
		fmt.Println("Error: While parsing json: " + err.Error())
		os.Exit(1)
	}

	matches, err := utils.SearchMatchLines(flags)
	if err != nil {
		fmt.Println("Error: While searching for matches: " + err.Error())
		os.Exit(1)
	}

	if jsonOut {
		matchesArr := mappers.MapChanToJsonFile(matches)
		out, err := utils.BuildJson(matchesArr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(out)
		return
	}

	utils.PrintMatches(matches, cl, matchStyle, context)
}

func parseFlags() (bool, *string, *string, int, int, *string, *string, string, bool, error) {
	showVersion := pflag.BoolP("version", "v", false, "print version information")
	exdir := pflag.StringP("exclude-dir", "e", "", "relative paths to ignore")
	exfile := pflag.StringP(
		"exclude-file",
		"x",
		"",
		"bash-style glob patterns of files to ignore (comma-separated).\nPattern \"noext\" can be used for files with no extension",
	)
	threadc := pflag.IntP("thread-count", "t", 1, "thread count to use for file parsing")
	context := pflag.IntP("context-size", "c", 2, "number of context lines to show around a matched line")
	root := pflag.StringP("root", "r", "./", "root directory to walk")
	config := pflag.String("config", "", "JSON for layout+theme")
	jsonOut := pflag.Bool("json", false, "print result in json format")

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: findstr [flags] <pattern>")
		fmt.Fprintln(os.Stderr, "Search for file content matching <pattern> under the given root.")
		fmt.Fprintln(os.Stderr)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if args := pflag.Args(); len(args) == 0 {
		if *showVersion {
			return true, nil, nil, 0, 0, nil, nil, "", false, nil
		}
		return false, nil, nil, 0, -1, nil, nil, "", false, errors.New(
			"you must provide a <pattern> to search for",
		)
	} else {
		return *showVersion, exdir, exfile, *threadc, *context, root, config, args[0], *jsonOut, nil
	}
}

func printVersion() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("Unable to determine version information.")
		return
	}
	if buildInfo.Main.Version != "" {
		fmt.Printf("Version: %s\n", buildInfo.Main.Version)
	} else {
		fmt.Println("Version: unknown")
	}
}
