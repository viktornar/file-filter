package slice

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
