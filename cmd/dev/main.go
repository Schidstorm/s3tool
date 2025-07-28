package main

import (
	"fmt"
	"os"

	"github.com/schidstorm/s3tool/pkg/boxes"
	"github.com/schidstorm/s3tool/pkg/cli"
)

func main() {
	err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing CLI arguments: %v\n", err)
		os.Exit(1)
	}

	app := boxes.NewApp(nil)

	app.Run()
}
