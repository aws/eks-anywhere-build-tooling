package slices

func Contains(s []string, str string) bool {
	for _, elem := range s {
		if elem == str {
			return true
		}
	}
	return false
}

func IsSubset(s1, s2 []string) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !Contains(s2, e) {
			return false
		}
	}
	return true
}
