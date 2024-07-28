package toolbox

func IsEmpty[T comparable](val T) bool {
	var empty T
	if val == empty {
		return true
	}

	return false
}
