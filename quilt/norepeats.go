package quilt

type NoRepeats struct{}

func (NoRepeats) Name() string { return "no blah blah" }
func (NoRepeats) Validates(x, y int, proposed GrannySquare, quilt Quilt) bool {
	var noRepeatsCoord = [8][2]int{
		{0, 1},
		{-1, 0}, {1, 0},
		{0, -1},
		{0, 2},
		{-2, 0}, {2, 0},
		{0, -2},
	}

	for _, pos := range noRepeatsCoord {
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
		if proposed.middle == neigh.middle || proposed.inner == neigh.inner {
			return false
		}
	}
	return true
}
