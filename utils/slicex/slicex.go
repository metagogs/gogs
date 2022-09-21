package slicex

func InSlice[T comparable](target T, data []T) bool {
	for _, s := range data {
		if s == target {
			return true
		}
	}

	return false
}

func RemoveSliceItem[T comparable](list []T, item T) []T {
	for i, v := range list {
		if v == item {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	return list
}
