package main

import (
	"image/png"
	"os"
	"time"

	"github.com/schidstorm/s3tool/pkg/e2e"
)

func main() {
	screens := map[string]func(e *e2e.Emulator){
		"startPage": func(e *e2e.Emulator) {
			time.Sleep(100 * time.Millisecond)
		},
		"bucketsPage": func(e *e2e.Emulator) {
			e.Send(e2e.KeyArrowDown + e2e.KeyArrowDown)
			e.Send(e2e.KeyEnter)
			time.Sleep(1 * time.Second)
		},
	}

	for name, handler := range screens {
		outFile := "screens/" + name + ".png"
		if err := generateScreen(outFile, handler); err != nil {
			panic(err)
		}
	}
}

func generateScreen(outFile string, handler func(e *e2e.Emulator)) error {
	emulator := &e2e.Emulator{}
	defer emulator.Close()
	err := emulator.Run("./build/s3tool", 80, 24)
	if err != nil {
		panic(err)
	}

	handler(emulator)

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	img := emulator.ContentImage()
	return png.Encode(f, img)
}
