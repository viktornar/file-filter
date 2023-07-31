package slice

func Reverse[T comparable](ss []T) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func Remove[T comparable](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func IndexOf[T comparable](collection []T, el T) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}
