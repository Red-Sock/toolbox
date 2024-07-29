package toolbox

func Contains[T comparable](a []T, b T) bool {
	for _, item := range a {
		if item == b {
			return true
		}
	}

	return false
}
