package quilt

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"math/rand"
	"os"
)

var defaultColor = color.RGBA{0, 0, 0, 0}

func clamp(x, max int) int {
	switch {
	case x < 0:
		return 0
	case x > max:
		return max
	default:
		return x
	}
}

// factorial returns the factorial of x.
func factorial(x int) int {
	if x == 0 {
		return 1
	}
	return x * factorial(x-1)
}

type GrannySquare struct {
	inner  color.RGBA
	middle color.RGBA
	outer  color.RGBA
}

func defaultSquare() GrannySquare {
	return GrannySquare{
		inner:  defaultColor,
		middle: defaultColor,
		outer:  defaultColor,
	}
}

type Quilt struct {
	width, height   int
	colorsPerSquare int
	colors          []color.RGBA
	squareSet       map[GrannySquare]bool
	grid            []GrannySquare
	frames          []image.Image
}

func (q *Quilt) set(x, y int, v GrannySquare) {
	x = clamp(x, q.width-1)
	y = clamp(y, q.height-1)
	q.grid[y*q.width+x] = v
	q.addFrame()
}

func (q *Quilt) get(x, y int) GrannySquare {
	x = clamp(x, q.width-1)
	y = clamp(y, q.height-1)
	return q.grid[y*q.width+x]
}

func New(width, height, colorsPerSquare int, colors []color.RGBA) *Quilt {
	q := &Quilt{
		width:           width,
		height:          height,
		colorsPerSquare: colorsPerSquare,
		colors:          colors,
		grid:            make([]GrannySquare, width*height),
	}
	// squareSet is a map of all available squares (combinations) computed ahead
	// of time.
	q.squareSet = q.generateAvailableSquares()
	return q
}

// Combinations returns the number of distinct non-repeating combinations from the supplied colors.
func (q *Quilt) Combinations() int {
	return factorial(len(q.colors)) / factorial(len(q.colors)-q.colorsPerSquare)
}

// Size returns the amount of squares in the grid.
func (q *Quilt) Size() int {
	return q.width * q.height
}

func (q *Quilt) generateAvailableSquares() map[GrannySquare]bool {
	set := make(map[GrannySquare]bool)
	for len(set) < q.Combinations() {
		sq := q.getUniqueSquare()
		if ok := set[sq]; !ok {
			set[sq] = true
		}
	}
	return set
}

func (q *Quilt) findSuitableSquare(x, y int) GrannySquare {
	for k := range q.squareSet {
		if q.PassesValidation(x, y, k) {
			return k
		}
	}
	return defaultSquare()
}

func (q *Quilt) GenerateQuilt() {
	// pre-seed the quilt with some colors. This for the most part will just place
	// random squares, as there are no neighbours to violate rules.
	for x := 0; x < q.width; x++ {
		for y := 0; y < q.height; y++ {
			sq := q.findSuitableSquare(x, y)
			q.set(x, y, sq)
		}
	}
	// This can run infinitely if the rules are too restrictive.
	q.findRuleOffenders()
}

// Loop through a pre-generated grid and replace squares that violate any rules.
// If there are no suitable squares for the exact co-ordinate then remove all
// neighours and try again. This will run infinitely until a grid that satisfies
// the rules completes.
func (q *Quilt) findRuleOffenders() {
	for x := 0; x < q.width; x++ {
		for y := 0; y < q.height; y++ {
			if !q.PassesValidation(x, y, q.get(x, y)) {
				sq := q.findSuitableSquare(x, y)
				q.set(x, y, sq)
				q.findRuleOffenders()
			}
		}
	}
}

// PassesValidation aggregates the rules and returns false if any of them fail,
// or true if otherwise.
func (q *Quilt) PassesValidation(x, y int, proposed GrannySquare) bool {
	if proposed == defaultSquare() {
		return false
	}
	if q.UsedMoreThanNTimes(3, proposed) {
		return false
	}
	if !q.NoSidesMatch(x, y, proposed) {
		return false
	}
	if !q.NoSameMiddleAndInner(x, y, proposed) {
		return false
	}

	if !q.NoInners(x, y, proposed) {
		return false
	}
	return true
}

func (q *Quilt) UsedMoreThanNTimes(n int, g GrannySquare) bool {
	seen := 0
	for x := 0; x < q.width; x++ {
		for y := 0; y < q.height; y++ {
			if q.get(x, y) == g {
				seen++
			}
		}
	}
	return seen > n
}

func (q *Quilt) NoSidesMatch(x, y int, proposed GrannySquare) bool {
	// positions contains the relative distance between the current square and
	// it's neighbours.
	var positions = [8][2]int{
		{-1, 1}, {0, 1}, {1, 1},
		{-1, 0}, {1, 0},
		{-1, -1}, {0, -1}, {1, -1},
	}

	for _, pos := range positions {
		x2, y2 := clamp(x+pos[0], q.width-1), clamp(y+pos[1], q.height-1)

		// we are checking ourselves...
		if x2 == x && y2 == y {
			continue
		}

		c := q.get(x2, y2)
		if proposed.outer == c.outer {
			return false
		}
	}
	return true
}

func (q *Quilt) NoSameMiddleAndInner(x, y int, proposed GrannySquare) bool {
	// positions contains the relative distance between the current square and
	// it's neighbours.
	var positions = [12][2]int{
		{-1, 1}, {0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
		{1, -1},
		{-1, -1},
		{1, 1},
		{0, 2},
		{-2, 0}, {2, 0},
		{0, -2},
	}

	for _, pos := range positions {
		x2, y2 := clamp(x+pos[0], q.width-1), clamp(y+pos[1], q.height-1)

		// we are checking ourselves...
		if x2 == x && y2 == y {
			continue
		}

		c := q.get(x2, y2)
		if proposed.middle == c.middle && proposed.inner == c.inner {
			return false
		}
	}
	return true
}

func (q *Quilt) NoInners(x, y int, proposed GrannySquare) bool {
	// positions contains the relative distance between the current square and
	// it's neighbours.
	var positions = [8][2]int{
		{0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
	}

	for _, pos := range positions {
		x2, y2 := clamp(x+pos[0], q.width-1), clamp(y+pos[1], q.height-1)

		// we are checking ourselves...
		if x2 == x && y2 == y {
			continue
		}

		c := q.get(x2, y2)
		if proposed.inner == c.inner {
			return false
		}
	}
	return true
}

func (q *Quilt) NoSameThreecolors(x, y int, proposed GrannySquare) bool {
	// positions contains the relative distance between the current square and
	// it's neighbours.
	var positions = [8][2]int{
		{0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
	}
	for _, pos := range positions {
		x2, y2 := clamp(x+pos[0], q.width-1), clamp(y+pos[1], q.height-1)

		// we are checking ourselves...
		if x2 == x && y2 == y {
			continue
		}

		c := q.get(x2, y2)
		if proposed.outer == c.inner || proposed.outer == c.middle || proposed.outer == c.outer {
			return false
		}
		if proposed.middle == c.inner || proposed.middle == c.middle || proposed.middle == c.outer {
			return false
		}
		if proposed.inner == c.inner || proposed.inner == c.middle || proposed.inner == c.outer {
			return false
		}
	}
	return true
}

// getUniqueSquare returns a GrannySquare that is guaranteed to have 3 distinct
// colors.
func (q *Quilt) getUniqueSquare() GrannySquare {
	cache := map[color.RGBA]bool{}
	gen := make([]color.RGBA, 0, q.colorsPerSquare)

	// Generate enough distinct colors per square.
	for len(gen) < q.colorsPerSquare {
		n := rand.Intn(len(q.colors))
		color := q.colors[n]
		if ok := cache[color]; !ok {
			cache[color] = true
			gen = append(gen, color)
		}
	}

	sq := GrannySquare{
		inner:  gen[0],
		middle: gen[1],
		outer:  gen[2],
	}

	return sq
}

// addFrame renders the current grid and appends it to the frames field in Quilt.
func (q *Quilt) addFrame() {
	squareSize := 40
	img := image.NewRGBA(image.Rect(0, 0, q.width*squareSize, q.height*squareSize))

	posX := 0
	for x := 0; x < q.width; x++ {
		posY := 0
		for y := 0; y < q.height; y++ {
			sq := q.get(x, y)
			squares := []color.RGBA{sq.outer, sq.middle, sq.inner}
			// Loop through each square and draw a smaller square inside it.
			for i := 0; i < 3; i++ {
				d := i * 5 // the distance between the outer edge.
				draw.Draw(img, image.Rect(posX+d, posY+d, posX+squareSize-d, posY+squareSize-d),
					&image.Uniform{squares[i]}, image.Point{}, draw.Src)
			}
			posY += squareSize
		}
		posX += squareSize
	}
	q.frames = append(q.frames, img)
}

// PNG takes a filename and generates a PNG representation of the latest frame.
func (q *Quilt) PNG(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	png.Encode(f, q.frames[len(q.frames)-1])
}

// GIF renders each frame (a snapshot every time set is called) as a GIF, saved
// to filename.
func (q *Quilt) GIF(filename string) {
	g := &gif.GIF{LoopCount: -1}
	// set the palette to the supplied colors as well as the default color, which
	// serves as the background.
	palette := []color.Color{defaultColor}
	for _, c := range q.colors {
		palette = append(palette, c)
	}
	for _, img := range q.frames {
		p := image.NewPaletted(img.Bounds(), palette)
		draw.Draw(p, p.Rect, img, img.Bounds().Min, draw.Over)
		g.Image = append(g.Image, p)
		g.Delay = append(g.Delay, 0)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	gif.EncodeAll(f, g)
}
