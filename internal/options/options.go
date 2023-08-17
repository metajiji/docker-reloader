package options

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/jessevdk/go-flags"
)

// options are command-line options that are provided by the user.
type Options struct {
	Config flags.Filename `long:"config" description:"Path to config file, required parameter." required:"true" default:"config.yml"`

	Version bool `short:"V" long:"version" description:"display the version and exit"`
}

func buildVersion(version, commit, date string) string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	result = fmt.Sprintf("%s\ngoos: %s\ngoarch: %s", result, runtime.GOOS, runtime.GOARCH)
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	result = fmt.Sprintf("%s\ngoversion: %s", result, runtime.Version())
	return result
}

func ParseArgs(version string, commit string, date string) Options {
	// Parse command line arguments
	var (
		args Options

		parser = flags.NewParser(&args, flags.Default)
	)

	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	if args.Version {
		println(buildVersion(version, commit, date))
		os.Exit(0)
	}

	return args
}
