package zfc

type Condition func(Set) bool

type Set interface {
	AreEqual(Set) bool
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

	// Unwind() (Set, bool)
}

type Operator interface {
	Empty() Set
	Set(interface{}) Set
}
