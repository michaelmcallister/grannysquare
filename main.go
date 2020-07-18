package main

import (
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

	yarns := []color.RGBA{red, white, blue, green, purple, yellow}

	q := quilt.New(16, 20, 3, yarns)
	q.GenerateQuilt()
	q.GIF("/tmp/granny.gif")
}
