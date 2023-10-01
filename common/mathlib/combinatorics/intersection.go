package combinatorics

func Intersection(from1, to1, from2, to2 int) (from, to, intersectionLength int) {
	if from1 < from2 {
		if to1 < from2 {
			return from1, to2, 0
		} else if to1 <= to2 {
			return from1, to2, to1 - from2
		}
		return from1, to1, to2 - from2

	}

	// else from2 <= from1

	if to2 < from1 {
		return from2, to1, 0
	} else if to2 <= to1 {
		return from2, to1, to2 - from1
	}
	return from2, to2, to1 - from1

}
