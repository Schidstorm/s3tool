package e2e

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	vt "github.com/hinshun/vt10x"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
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

type Emulator struct {
	term vt.View
	ptmx *os.File
}

func (e *Emulator) Run(binary string, cols, rows int) error {
	cmd := exec.Command(binary)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	cmd.Stderr = os.Stderr

	// Start PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	_ = pty.Setsize(ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})

	var buf bytes.Buffer
	out := io.TeeReader(ptmx, &buf)

	// Create emulator
	term := vt.New(vt.WithSize(cols, rows))
	// Pump PTY output -> emulator
	go func() {
		r := bufio.NewReader(out)
		for {
			if err := term.Parse(r); err != nil {
				if err == io.EOF {
					return
				}
				// Parse may return quickly; small sleep avoids tight loop on EAGAIN
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	e.term = term
	e.ptmx = ptmx

	return nil
}

func (e *Emulator) Resize(cols, rows int) error {
	if e.ptmx == nil {
		return io.ErrClosedPipe
	}
	if err := pty.Setsize(e.ptmx, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)}); err != nil {
		return err
	}
	e.term.Resize(cols, rows)
	return nil
}

func (e *Emulator) Close() error {
	if e.ptmx != nil {
		if err := e.ptmx.Close(); err != nil {
			return err
		}
		e.ptmx = nil
	}
	return nil
}

func (e *Emulator) Send(content string) error {
	if e.ptmx == nil {
		return io.ErrClosedPipe
	}
	time.Sleep(100 * time.Millisecond)
	_, err := e.ptmx.Write([]byte(content))
	return err
}

func (e *Emulator) ContentString() string {
	glyphs := e.Content()
	var data []rune
	for _, row := range glyphs {
		for _, g := range row {
			data = append(data, g.Char)
		}
		data = append(data, '\n')
	}

	return string(data)
}

func (e *Emulator) ContentHtml() string {
	_, rows := e.term.Size()
	content := e.Content()

	var lines [][]coloredString
	for y := range rows {
		lines = append(lines, compactGlyphs(content[y]))
	}

	return glyphsToHtml(lines)
}

func (e *Emulator) Content() [][]vt.Glyph {
	if e.term == nil {
		return nil
	}

	e.term.Lock()
	defer e.term.Unlock()

	cols, rows := e.term.Size()
	content := make([][]vt.Glyph, rows)
	for y := 0; y < rows; y++ {
		row := make([]vt.Glyph, cols)
		for x := 0; x < cols; x++ {
			row[x] = e.term.Cell(x, y)
		}
		content[y] = row
	}
	return content
}

type coloredString struct {
	s  string
	fg vt.Color
	bg vt.Color
}

func compactGlyphs(glyphs []vt.Glyph) []coloredString {
	var result []coloredString
	for _, g := range glyphs {
		if len(result) == 0 || result[len(result)-1].fg != g.FG || result[len(result)-1].bg != g.BG {
			result = append(result, coloredString{s: string([]rune{g.Char}), fg: g.FG, bg: g.BG})
		} else {
			result[len(result)-1].s += string([]rune{g.Char})
		}
	}
	return result
}

func colorToCssColor(c vt.Color) string {
	if c == vt.DefaultFG {
		return "inherit"
	}
	if c == vt.DefaultBG {
		return "transparent"
	}

	switch c {
	case vt.Black:
		return "black"
	case vt.Red:
		return "#800000"
	case vt.Green:
		return "#008000"
	case vt.Yellow:
		return "#008000"
	case vt.Blue:
		return "#000080"
	case vt.Magenta:
		return "#000080"
	case vt.Cyan:
		return "#008080"
	case vt.LightGrey:
		return "#c0c0c0"
	case vt.DarkGrey:
		return "#808080"
	case vt.LightRed:
		return "#ff0000"
	case vt.LightGreen:
		return "#00ff00"
	case vt.LightYellow:
		return "#ffff00"
	case vt.LightBlue:
		return "#0000ff"
	case vt.LightMagenta:
		return "#ff00ff"
	case vt.LightCyan:
		return "#00ffff"
	case vt.White:
		return "white"
	default:
		return "black" // Fallback for unknown colors
	}
}

func glyphsToHtml(glyphs [][]coloredString) string {
	var html string
	for _, row := range glyphs {
		html += "<div style='white-space: pre;>"
		for _, g := range row {
			fgColor := colorToCssColor(g.fg)
			bgColor := colorToCssColor(g.bg)
			html += "<span style='color:" + fgColor + "; background-color:" + bgColor + ";>" + g.s + "</span>"
		}
		html += "</div>"
	}
	return "<div style='font-family: monospace, monospace;'>" + html + "</div>"
}

func (e *Emulator) ContentImage() image.Image {
	cellW, cellH := 8, 16 // Fixed width and height for basic font
	cols, rows := e.term.Size()
	img := image.NewRGBA(image.Rect(0, 0, cols*cellW, rows*cellH))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.White,        // simple: draw all glyphs white
		Face: basicfont.Face7x13, // fixed width font
	}

	// very simple palette mapping: you can map vt10x Color -> RGBA here
	bg := color.RGBA{0, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	glyphs := e.Content()
	for y := range rows {
		for x := range cols {
			g := glyphs[y][x]
			if g.Char == 0 || g.Char == ' ' {
				continue
			}
			// TODO: map g.FG/g.BG to actual colors; handle bold/attr, double-width runes, etc.
			drawer.Dot = fixed.P((x * cellW), (y*cellH)+int(drawer.Face.Metrics().Ascent.Ceil()))
			drawer.DrawString(string(g.Char))
		}
	}

	return img
}
