package slicex

func InSlice(target string, data []string) bool {
	for _, s := range data {
		if s == target {
			return true
		}
	}

	return false
}
