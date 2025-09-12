package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/HubertasVin/findstr/mappers"
	"github.com/HubertasVin/findstr/models"
	"github.com/HubertasVin/findstr/utils"
	"github.com/spf13/pflag"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	signal.Ignore(syscall.SIGPIPE)

	showVersion, exdir, exfile, threadc, contextSize, root, pattern, jsonOut, createConfig, err := parseFlags()
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

	if createConfig {
		path, err := utils.CreateDefaultConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created config at " + path)
		return
	}

	flags := models.ProgramFlags{
		ExcludeDir:  *exdir,
		ExcludeFile: *exfile,
		ThreadCount: threadc,
		ContextSize: contextSize,
		Root:        *root,
		Json:        jsonOut,
		Pattern:     pattern,
	}

	if threadc <= 0 {
		fmt.Println("Error: Thread count must be greater than 0")
		os.Exit(1)
	}
	if contextSize < 0 {
		fmt.Println("Error: Context size must be greater than or equal to 0")
		os.Exit(1)
	}

	cl, theme, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("Error: While loading config: " + err.Error())
		os.Exit(1)
	}

	matches, err := utils.SearchMatchLines(ctx, flags)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			os.Exit(130) // interrupted
		}
		fmt.Println("Error: While searching for matches: " + err.Error())
		os.Exit(1)
	}

	if jsonOut {
		matchesArr := mappers.MapChanToJsonFile(ctx, matches)
		out, err := utils.BuildJson(matchesArr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(out)
		return
	}

	utils.PrintMatches(ctx, matches, cl, theme, contextSize)

	if ctx.Err() != nil {
		// ensure styles are reset and we end on a fresh line
		fmt.Fprint(os.Stdout, "\x1b[0m\x1b[K\n")
		os.Exit(130)
	}
}

func parseFlags() (bool, *string, *string, int, int, *string, string, bool, bool, error) {
	showVersion := pflag.BoolP("version", "v", false, "print version information")
	exdir := pflag.StringP("exclude-dir", "e", "", "relative paths to ignore")
	exfile := pflag.StringP(
		"exclude-file",
		"x",
		"",
		"bash-style glob patterns of files to ignore (comma-separated).\nPattern \"noext\" can be used for files with no extension",
	)
	threadc := pflag.IntP("thread", "t", 1, "thread count to use for file parsing")
	context := pflag.IntP("context", "c", 2, "number of context lines to show around a matched line")
	root := pflag.StringP("root", "r", "./", "root directory to walk")
	jsonOut := pflag.Bool("json", false, "print result in json format")
	createConfig := pflag.Bool("create-config", false, "create default config at $HOME/.config/findstr.toml and exit")

	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: findstr [flags] <pattern>")
		fmt.Fprintln(os.Stderr, "Search for file content matching <pattern> under the given root.")
		fmt.Fprintln(os.Stderr)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if args := pflag.Args(); len(args) == 0 {
		if *showVersion {
			return *showVersion, nil, nil, 0, 0, nil, "", false, false, nil
		}
		if *createConfig {
			return false, nil, nil, 0, 0, nil, "", false, *createConfig, nil
		}
		return false, nil, nil, 0, -1, nil, "", false, *createConfig, errors.New(
			"you must provide a <pattern> to search for",
		)
	} else {
		return *showVersion, exdir, exfile, *threadc, *context, root, args[0], *jsonOut, *createConfig, nil
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
