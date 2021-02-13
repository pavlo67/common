package selectors

// to be continued... -----------------------------------------

type DNF [][]Proposition

type Proposition struct {
	Assertion string
	IsNot     bool
}
