package main

import (
	"fmt"
	"os"

	"github.com/schidstorm/s3tool/internal/cli"
	"github.com/schidstorm/s3tool/internal/s3lib"
	"github.com/schidstorm/s3tool/internal/terminal"
)

func main() {
	err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing CLI arguments: %v\n", err)
		os.Exit(1)
	}

	app := terminal.NewApp(nil, loaders()...)
	app.Run()

}

func loaders() []s3lib.ConnectorLoader {
	var loaders []s3lib.ConnectorLoader
	if cli.Config.Loaders.Aws {
		loaders = append(loaders, &s3lib.AwsLoader{})
	}
	if cli.Config.Loaders.S3Tool {
		loaders = append(loaders, &s3lib.S3ToolLoader{})
	}
	if cli.Config.Loaders.Memory {
		loaders = append(loaders, &s3lib.MemoryLoader{})
	}
	return loaders
}
