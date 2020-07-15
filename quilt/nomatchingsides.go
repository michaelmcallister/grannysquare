package quilt

type NoMatchingSides struct{}

func (NoMatchingSides) Name() string { return "No sides with the same colour" }
func (NoMatchingSides) Validates(x, y int, proposed GrannySquare, quilt Quilt) bool {
	var positions = [8][2]int{
		{0, 1}, {1, 1}, {1, 0},
		{-1, -1}, {0, -1}, {1, -1},
	}
	for _, pos := range positions {
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

		neigh := quilt.grid[xCheck][yCheck]
		if proposed.outer == neigh.outer {
			return false
		}
	}
	return true
}
