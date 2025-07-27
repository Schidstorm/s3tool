package main

import (
	"github.com/schidstorm/s3tool/pkg/boxes"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

func main() {
	client, err := s3lib.New()
	if err != nil {
		panic(err)
	}

	box := boxes.NewBucketsBox(client)
	app := boxes.NewApp()
	app.OpenPage(box)

	app.Run()
}
