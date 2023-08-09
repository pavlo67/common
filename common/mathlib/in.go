package mathlib

func InInt(ints []int, val int) bool {
	for _, n := range ints {
		if n == val {
			return true
		}
	}
	return false
}

//func Index(strs []string, str string) int {
//	for i, s := range strs {
//		if s == str {
//			return i
//		}
//	}
//	return -1
//}
//
