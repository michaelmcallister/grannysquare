package quilt

// Rule defines the contract for implementing a rule that determines if a proposed grannySquare is legal.
type Rule interface {
	Name() string
	Validates(x, y int, proposed GrannySquare, quilt Quilt) bool
}

func GetRules() []Rule {
	return []Rule{
		NoMatchingSides{},
		ColourSets{},
		NoRepeats{},
	}
}
