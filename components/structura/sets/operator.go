package sets

import "github.com/pavlo67/constructor/components/basis"

type Condition func(Something) bool

// TODO!!! test order conditions
type Compare func(Something, Something) int64

type Something interface {
	IsEqual(Something) bool
	DifferenceFrom(Something) []Something
	Variations(basis.Info) []Something

	Configure(basis.Info)
}

type Set interface {
	Something
	Contains(Something) bool
	Elements() []Something

	Cardinality() int64
	Height() int64

	Sort(Compare) Set
	Factor(Compare) Set

	Tuple(...Set) Set
	Choice(...Set) Set
	Union(...Set) Set
	Intersection(...Set) Set
	Subtraction(Set) Set
	Delta(Set) Set
	Map(Set) Set

	SubSet(Condition) Set
	SubPower(Condition, Set) Set
	SubBoolean(Condition) Set
	SubProduct(Condition, ...Set) Set
}

type Operator interface {
	Empty() Set
	Set(...Something) Set
	Range(from, to uint64) Set
}
