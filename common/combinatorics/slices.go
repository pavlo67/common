package combinatorics

func Multiple[T interface{}](data [][]T) (dataMultipled [][]T) {

	for _, optionsSet := range data {
		var optionsMultipledNew [][]T

		for _, optionValue := range optionsSet {
			for _, opts := range dataMultipled {
				optionsMultipledNew = append(optionsMultipledNew, append(opts, optionValue))
			}
		}

		dataMultipled = optionsMultipledNew
	}

	return dataMultipled
}

func Intersects[T comparable](set0, set1 []T) bool {

	for _, v0 := range set0 {
		for _, v1 := range set1 {
			if v0 == v1 {
				return true
			}
		}
	}

	return false
}
