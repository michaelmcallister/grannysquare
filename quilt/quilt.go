package quilt

import (
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

type Quilt struct {
	width, height   int
	colorsPerSquare int
	colors          []color.RGBA
	squareSet       map[GrannySquare]bool
	grid            [][]GrannySquare
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
	q.squareSet = make(map[GrannySquare]bool)
	return q
}

// Combinations returns the number of distinct non-repeating combinations from the supplied colors.
func (q *Quilt) Combinations() int {
	return factorial(len(q.colors)) / factorial(len(q.colors)-q.colorsPerSquare)
}

func (q *Quilt) Size() int {
	return q.width * q.height
}

func (q *Quilt) UseLimit() int {
	a := float64(float64(q.width*q.height) / float64(len(q.colors))) // 53.333
	b := float64(a * float64(q.colorsPerSquare))                     // 319.99
	return int(b)
}

func (q *Quilt) PlaceSquares() {
	squares := q.squareSet
	placed := make(map[GrannySquare]int)

	for x := 0; x < q.width; x++ {
		for y := 0; y < q.height; y++ {
			var match GrannySquare
			foundMatch := false
			for sq := range squares {
				match = sq
				foundMatch = true
				break
			}
		}

		if !foundMatch {
			fmt.Printf("couldn't find a match for x=%d,y=%d\n", x, y)
			// q.Draw("/tmp/granny.png")
			panic("nope")
			//q.PlaceSquares()
		}

		q.grid[x][y] = match
		placed[match]++
		if placed[match] == 3 {
			delete(squares, match)
		}
	}
}

// func (q *Quilt) PlaceSquares() {
// 	squares := q.squareSet
// 	placed := make(map[GrannySquare]int)

// 	for x := 0; x < q.width; x++ {
// 		for y := 0; y < q.height; y++ {
// 			var match GrannySquare
// 			foundMatch := false
// 			for sq := range squares {
// 				if q.PassesValidation(x, y, sq) {
// 					match = sq
// 					foundMatch = true
// 					break
// 				}
// 			}

// 			if !foundMatch {
// 				fmt.Printf("couldn't find a match for x=%d,y=%d\n", x, y)
// 				// q.Draw("/tmp/granny.png")
// 				panic("nope")
// 				//q.PlaceSquares()
// 			}

// 			q.grid[x][y] = match
// 			placed[match]++
// 			if placed[match] == 3 {
// 				delete(squares, match)
// 			}
// 		}
// 	}
// }

func (q *Quilt) PassesValidation(x, y int, proposed GrannySquare) bool {
	for _, r := range GetRules() {
		if !r.Validates(x, y, proposed, *q) {
			return false
		}
	}
	return true
}

// GenerateSquares pre-generates all combinations of squares available.
func (q *Quilt) GenerateSquares() map[GrannySquare]bool {
	set := make(map[GrannySquare]bool) // New empty set

	for i := 0; i < q.Combinations(); i++ {
		for {
			sq := q.GetUniqueSquare()
			if ok := set[sq]; !ok {
				set[sq] = true
				break
			}
		}
	}
	return set
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
	// If this a unique square return it.
	if ok := q.squareSet[gen]; !ok {
		q.squareSet[gen] = true
		return gen
	}

	// Or try again.
	return q.GetUniqueSquare()
}

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

// func (q *Quilt) Used() string {
// 	var s strings.Builder
// 	m := make(map[color.RGBA]int)
// 	s.WriteString(fmt.Sprintf("Used the following colours:\n"))
// 	for square := range q.squareSet {
// 		m[square.inner]++
// 		m[square.middle]++
// 		m[square.outer]++
// 	}
// 	s.WriteString(fmt.Sprintf("%v\n", m))
// 	return s.String()
// }
