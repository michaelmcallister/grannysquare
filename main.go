package main

import (
	"fmt"
	"image/color"

	"github.com/michaelmcallister/grannysquare/quilt"
)

func main() {
	red := color.RGBA{0xFF, 0, 0, 0xFF}
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	blue := color.RGBA{0, 0, 0xFF, 0xFF}
	green := color.RGBA{0, 0xFF, 0, 0xFF}
	purple := color.RGBA{80, 0, 80, 0xFF}
	yellow := color.RGBA{0xFF, 0xFF, 0, 0xFF}
	/*pink := color.RGBA{0xFF, 0xC0, 0xCB, 0xFF}
	grey := color.RGBA{0x80, 0x80, 0x80, 0xFF}
	frestgreen := color.RGBA{0x22, 0x8B, 0x22, 0xFF}
	brown := color.RGBA{0xD2, 0x69, 0x1E, 0xFF}
	orange := color.RGBA{0xFF, 0xA5, 0x00, 0xFF}*/

	yarns := []color.RGBA{red, white, blue, green, purple, yellow}

	q := quilt.New(16, 20, 3, yarns)
	fmt.Printf("quilt has %d combinations...\n", q.Combinations())
	fmt.Printf("generating all permuations...\n")
	out := q.GenerateSquares()
	fmt.Printf("generated %d squares\n", len(out))
	q.GenerateQuilt()
	q.Draw("/tmp/granny.png")
}
