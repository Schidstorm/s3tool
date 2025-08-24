package main

import (
	"image/png"
	"os"
	"time"

	"github.com/schidstorm/s3tool/pkg/cli"
	"github.com/schidstorm/s3tool/pkg/emulator"
)

func init() {
	cli.Config.Loaders.Aws = false
	cli.Config.Loaders.S3Tool = false
	cli.Config.Loaders.Memory = true
}

var screenColumns = 160
var screenRows = 48
var imageWidth = 1600

var screens = map[string]func(e *emulator.Emulator){
	"start_page": func(e *emulator.Emulator) {
	},
	"buckets_page": func(e *emulator.Emulator) {
		e.Send(emulator.KeyEnter)
	},
	"objects_page": func(e *emulator.Emulator) {
		e.Send(emulator.KeyEnter)
		e.Send(emulator.KeyEnter)
	},
}

func main() {
	for name, handler := range screens {
		outFile := "screens/" + name + ".png"
		if err := generateScreen(outFile, handler); err != nil {
			panic(err)
		}
	}
}

func generateScreen(outFile string, handler func(e *emulator.Emulator)) error {
	emulator := emulator.NewEmulator(&ScreenLoader{})
	defer emulator.Close()
	go func() {
		if err := emulator.Run(screenColumns, screenRows); err != nil {
			panic(err)
		}
	}()

	handler(emulator)
	time.Sleep(100 * time.Millisecond) // Allow time for the screen to render

	img := emulator.ContentImage(imageWidth)
	pngFile, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer pngFile.Close()
	if err := png.Encode(pngFile, img); err != nil {
		return err
	}
	return nil
}
