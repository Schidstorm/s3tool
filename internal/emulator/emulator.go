package emulator

import (
	"errors"
	"image"
	"image/color"

	"golang.org/x/image/draw"

	_ "embed"

	"github.com/gdamore/tcell/v2"
	"github.com/schidstorm/s3tool/internal/s3lib"
	"github.com/schidstorm/s3tool/internal/terminal"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	KeyArrowUp    string = "\x1b[A"
	KeyArrowDown  string = "\x1b[B"
	KeyArrowLeft  string = "\x1b[D"
	KeyArrowRight string = "\x1b[C"
	KeyEnter      string = "\r\n"
	KeyEscape     string = "\x1b"
)

var ttfs = map[tcell.AttrMask][]byte{
	tcell.AttrBold:                    gomonobold.TTF,
	tcell.AttrNone:                    gomono.TTF,
	tcell.AttrItalic:                  gomonoitalic.TTF,
	tcell.AttrBold | tcell.AttrItalic: gomonobolditalic.TTF,
}

type Emulator struct {
	app terminal.SimulatedApp
}

func NewEmulator(loaders ...s3lib.ConnectorLoader) *Emulator {
	return &Emulator{
		app: terminal.NewSimulatedApp(nil, loaders...),
	}
}

func (e *Emulator) Run(cols, rows int) error {
	e.app.GetScreen().SetSize(cols, rows)
	return e.app.Run()
}

func (e *Emulator) Close() {
	e.app.App.Application.Stop()
}

func (e *Emulator) Send(content string) error {
	if e.app.GetScreen().InjectKeyBytes([]byte(content)) {
		return nil
	} else {
		return errors.New("not all keys were fully understood")
	}
}

func (e *Emulator) ContentString() string {
	cells, width, height := e.app.GetScreen().GetContents()
	data := make([]rune, 0, width*height)
	var x, y int
	for _, cell := range cells {
		if x == width {
			x = 0
			y++
			data = append(data, '\n')
		}

		data = append(data, cell.Runes...)
	}

	return string(data)
}

func (e *Emulator) ContentImage(finalImageWidth int) image.Image {
	cells, cols, rows := e.app.GetScreen().GetContents()
	return generateImage(cells, cols, rows, finalImageWidth)
}

func generateImage(cells []tcell.SimCell, cols, rows, finalImageWidth int) image.Image {
	cellPixelWidth, cellPixelHeight := getCellSize()
	imageWidth := cols * cellPixelWidth
	imageHeight := rows * cellPixelHeight

	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	for y := range rows {
		for x := range cols {
			g := cells[y*cols+x]

			drawer := createDrawer(img, x, y, g)
			drawCellBackground(img, drawer, g)
			drawer.DrawString(string(g.Runes))
		}
	}

	if finalImageWidth > 0 {
		scaledImg := image.NewRGBA(image.Rect(0, 0, finalImageWidth, finalImageWidth*imageHeight/imageWidth))
		draw.NearestNeighbor.Scale(scaledImg, scaledImg.Bounds(), img, img.Bounds(), draw.Over, nil)
		return scaledImg
	}

	return img
}

func createDrawer(img *image.RGBA, x, y int, cell tcell.SimCell) *font.Drawer {
	cellWidth, cellHeight := getCellSize()
	fg, _, attr := cell.Style.Decompose()
	fgColor := tcellColorToColor(fg)

	drawer := &font.Drawer{
		Src:  image.NewUniform(fgColor),
		Face: getFace(attr),
	}

	drawer.Dot.X = fixed.I(x * cellWidth)
	cellY := y*cellHeight + cellHeight
	drawer.Dot.Y = fixed.I(cellY)
	drawer.Dst = img
	return drawer
}

func drawCellBackground(img *image.RGBA, drawer *font.Drawer, cell tcell.SimCell) {
	_, bg, _ := cell.Style.Decompose()
	if bg != tcell.ColorDefault {
		cellX := drawer.Dot.X.Round()
		cellY := drawer.Dot.Y.Round() - (drawer.Face.Metrics().Height - drawer.Face.Metrics().Descent).Round()
		cellWidth, cellHeight := getCellSize()
		bgColor := tcellColorToColor(bg)
		draw.Draw(
			img,
			image.Rect(
				cellX,
				cellY,
				cellX+cellWidth,
				cellY+cellHeight,
			),
			&image.Uniform{C: bgColor},
			image.Point{},
			draw.Src,
		)
	}
}

func getCellSize() (width, height int) {
	face := getFace(tcell.AttrNone)
	width = font.MeasureString(face, " ").Round()
	height = (face.Metrics().Ascent + face.Metrics().Descent).Round()
	return
}

func tcellColorToColor(c tcell.Color) color.RGBA {
	fgR, fgG, fgB := c.TrueColor().RGB()
	return color.RGBA{
		R: uint8(fgR),
		G: uint8(fgG),
		B: uint8(fgB),
		A: 255,
	}
}

func getFace(attr tcell.AttrMask) font.Face {
	attr = attr & (tcell.AttrBold | tcell.AttrItalic)
	font, err := opentype.Parse(ttfs[attr])
	if err != nil {
		panic(err) // This should never happen, as the TTFs are embedded
	}
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 34, // 15 and 34 works well
		DPI:  72,
	})
	if err != nil {
		panic(err) // This should never happen, as the TTFs are embedded
	}
	return face
}
