package zfc

import "github.com/pavlo67/constructor/basis"

type Condition func(Something) bool

type Something interface {
	IsEqual(Something) bool
	DifferenceFrom(Something) []Something
	Variations(basis.Info) []Something
}

type Set interface {
	Something
	Contains(Something) bool
	Elements() []Something

	Cardinality() int64
	Height() int64

	Tuple(...Set) Set
	Boolean() Set
	Union(...Set) Set
	Selection(Condition) Set
	Choice(...Set) Set

	Intersection(...Set) Set
	Subtraction(Set) Set
	Delta(Set) Set
	Carthesian(Set) Set
	Map(Set) Set
}

type Operator interface {
	Empty() Set
	Set(...Something) Set
}
