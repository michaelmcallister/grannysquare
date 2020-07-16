package quilt

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
)

type GrannySquare struct {
	inner  color.RGBA
	middle color.RGBA
	outer  color.RGBA
}

func (g GrannySquare) String() string {
	inner := fmt.Sprintf("#%.2x%.2x%.2x", g.inner.R, g.inner.G, g.inner.B)
	middle := fmt.Sprintf("#%.2x%.2x%.2x", g.middle.R, g.middle.G, g.middle.B)
	outer := fmt.Sprintf("#%.2x%.2x%.2x", g.outer.R, g.outer.G, g.outer.B)
	return fmt.Sprintf("{outer:%s, middle:%s, inner:%s}", outer, middle, inner)
}

type Quilt struct {
	width, height   int
	colorsPerSquare int
	colors          []color.RGBA
	squareSet       map[GrannySquare]bool
	grid            [][]GrannySquare
}

func getCoord(x, y int) (int, int) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if y > 19 {
		y = 19
	}
	if x > 15 {
		x = 15
	}

	return x, y
}

func defaultSquare() GrannySquare {
	return GrannySquare{
		inner:  color.RGBA{0, 0, 0, 0},
		middle: color.RGBA{0, 0, 0, 0},
		outer:  color.RGBA{0, 0, 0, 0},
	}
}

// factorial returns the factorial of x.
func factorial(x int) int {
	if x == 0 {
		return 1
	}
	return x * factorial(x-1)
}

func New(width, height, colorsPerSquare int, colors []color.RGBA) *Quilt {
	q := &Quilt{
		width:           width,
		height:          height,
		colorsPerSquare: colorsPerSquare,
		colors:          colors,
	}
	q.grid = make([][]GrannySquare, q.width)
	for x := 0; x < q.width; x++ {
		q.grid[x] = make([]GrannySquare, q.height)
	}
	// squareSet is a map of all available squares and whether they are available
	// for use.
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
		sq := q.GetUniqueSquare()
		if ok := set[sq]; !ok {
			set[sq] = true
		}
	}
	return set
}

func (q *Quilt) findSuitableSquare(x, y int) (*GrannySquare, error) {
	for k, ok := range q.squareSet {
		if !ok {
			// not available for use, move on.
			continue
		}
		if q.PassesValidation(x, y, k) {
			return &k, nil
		}
	}
	return nil, errors.New("unable to find a suitable square")
}

func (q *Quilt) GenerateQuilt() {
	// pre-seed the quilt with some colours. This for the most part will just place
	// random squares, as there are no neighbours to violate rules.
	for x := 0; x < q.width; x++ {
		for y := 0; y < q.height; y++ {
			sq, err := q.findSuitableSquare(x, y)
			if err == nil {
				q.grid[x][y] = *sq
			}
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
			if !q.PassesValidation(x, y, q.grid[x][y]) {
				sq, err := q.findSuitableSquare(x, y)
				if err == nil {
					q.grid[x][y] = *sq
				} else {
					q.DeleteNeighbours(x, y)
				}
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
	// if !q.NoSameThreeColours(x, y, proposed) {
	// 	return false
	// }
	return true
}

func (q *Quilt) UsedMoreThanNTimes(n int, g GrannySquare) bool {
	seen := 0
	for x := 0; x < q.width; x++ {
		for y := 0; y < q.height; y++ {
			if q.grid[x][y] == g {
				seen++
			}
		}
	}
	return seen > n
}

func (q *Quilt) DeleteNeighbours(x, y int) {
	var positions = [8][2]int{
		{-1, 1}, {0, 1}, {1, 1},
		{-1, 0}, {1, 0},
		{-1, -1}, {0, -1}, {1, -1},
	}

	for _, pos := range positions {
		x2, y2 := getCoord(x+pos[0], y+pos[1])
		// we are checking ourselves...
		if x2 == x && y2 == y {
			continue
		}

		q.grid[x2][y2] = defaultSquare()
	}
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
		x2, y2 := getCoord(x+pos[0], y+pos[1])

		// we are checking ourselves...
		if x2 == x && y2 == y {
			return true
		}

		if proposed.outer == q.grid[x2][y2].outer {
			return false
		}
	}
	return true
}

func (q *Quilt) NoSameMiddleAndInner(x, y int, proposed GrannySquare) bool {
	// positions contains the relative distance between the current square and
	// it's neighbours.
	var positions = [4][2]int{
		{0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
	}

	for _, pos := range positions {
		x2, y2 := getCoord(x+pos[0], y+pos[1])

		// we are checking ourselves...
		if x2 == x && y2 == y {
			return true
		}

		if proposed.middle == q.grid[x2][y2].middle && proposed.inner == q.grid[x2][y2].inner {
			return false
		}
	}
	return true
}

func (q *Quilt) NoSameThreeColours(x, y int, proposed GrannySquare) bool {
	// positions contains the relative distance between the current square and
	// it's neighbours.
	var positions = [4][2]int{
		{0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
	}
	for _, pos := range positions {
		x2, y2 := getCoord(x+pos[0], y+pos[1])

		// we are checking ourselves...
		if x2 == x && y2 == y {
			return true
		}

		c := q.grid[x2][y2]
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

func (q *Quilt) GetUniqueSquare() GrannySquare {
	cache := map[color.RGBA]bool{}
	generated := make([]color.RGBA, 0, q.colorsPerSquare)

	// Generate enough distinct colours per square.
	for len(generated) < q.colorsPerSquare {
		n := rand.Intn(len(q.colors))
		color := q.colors[n]
		if ok := cache[color]; !ok {
			cache[color] = true
			generated = append(generated, color)
		}
	}

	gen := GrannySquare{
		inner:  generated[0],
		middle: generated[1],
		outer:  generated[2],
	}

	return gen
}

// Draw takes a filename and generates a PNG representation of the current grid.
func (q *Quilt) Draw(filename string) {
	squareSize := 40
	quiltImage := image.NewRGBA(image.Rect(0, 0, q.width*squareSize, q.height*squareSize))

	loc_x := 0
	for x := 0; x < q.width; x++ {
		loc_y := 0
		for y := 0; y < q.height; y++ {
			outer := q.grid[x][y].outer
			middle := q.grid[x][y].middle
			inner := q.grid[x][y].inner

			draw.Draw(quiltImage, image.Rect(loc_x, loc_y, loc_x+squareSize, loc_y+squareSize),
				&image.Uniform{outer}, image.ZP, draw.Src)

			draw.Draw(quiltImage, image.Rect(loc_x+5, loc_y+5, loc_x+squareSize-5, loc_y+squareSize-5),
				&image.Uniform{middle}, image.ZP, draw.Src)

			draw.Draw(quiltImage, image.Rect(loc_x+10, loc_y+10, loc_x+squareSize-10, loc_y+squareSize-10),
				&image.Uniform{inner}, image.ZP, draw.Src)
			loc_y += squareSize
		}
		loc_x += squareSize
	}
	myfile, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	defer myfile.Close()
	png.Encode(myfile, quiltImage)
}
