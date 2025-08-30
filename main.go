package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/HubertasVin/findstr/mappers"
	"github.com/HubertasVin/findstr/models"
	"github.com/HubertasVin/findstr/utils"
	"github.com/spf13/pflag"
)

func main() {
	exdir, exfile, threadc, context, root, style, pattern, json, err := parseFlags()
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
		ContextSize: context,
		Root:        *root,
		Style:       *style,
		Json:        json,
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

	styleVal, err := utils.ParseStyle(*style)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	matches, err := utils.SearchMatchLines(flags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if json {
		matchesArr := mappers.MapChanToJsonFile(matches)
		out, err := utils.BuildJson(matchesArr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(out)
	} else {
		utils.PrintMatches(matches, styleVal)
	}
}

func parseFlags() (*string, *string, int, int, *string, *string, string, bool, error) {
	exdir := pflag.StringP("exclude-dir", "e", "", "relative paths to ignore")
	exfile := pflag.StringP(
		"exclude-file",
		"x",
		"",
		"bash-style glob patterns of files to ignore (comma-separated).\nPattern \"noext\" can be used for files with no extension.",
	)
	threadc := pflag.IntP("thread-count", "t", 1, "thread count to use for file parsing.")
	context := pflag.IntP("context-size", "c", 2, "number of context lines to show around a matched line.")
	root := pflag.StringP("root", "r", "./", "root directory to walk")
	style := pflag.StringP("style", "", "", "custom style in valid json format for highlighting.")
	json := pflag.BoolP("json", "", false, "print result in json format.")

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: findstr [flags] <pattern>")
		fmt.Fprintln(os.Stderr, "Search for file content matching <pattern> under the given root.")
		fmt.Fprintln(os.Stderr)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if args := pflag.Args(); len(args) == 0 {
		return nil, nil, 0, -1, nil, nil, "", false, errors.New("you must provide a <pattern> to search for")
	} else {
		return exdir, exfile, *threadc, *context, root, style, args[0], *json, nil
	}
}
