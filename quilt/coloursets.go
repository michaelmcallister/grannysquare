package quilt

type ColourSets struct{}

func (ColourSets) Name() string { return "no blah blah" }
func (ColourSets) Validates(x, y int, proposed GrannySquare, quilt Quilt) bool {
	var noSameColourSets = [4][2]int{
		{0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
	}

	for _, pos := range noSameColourSets {
		xCheck := x + pos[0]
		yCheck := y + pos[1]
		if xCheck < 0 {
			xCheck = 0
		}
		if yCheck < 0 {
			yCheck = 0
		}
		if yCheck > 19 {
			yCheck = 19
		}
		if xCheck > 15 {
			xCheck = 15
		}

		if xCheck == 0 || yCheck == 0 {
			return true
		}

		neigh := quilt.grid[xCheck][yCheck]
		if proposed.middle == neigh.middle && proposed.inner == neigh.inner {
			return false
		}
	}
	return true
}
