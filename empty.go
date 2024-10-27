package toolbox

func IsEmpty[T comparable](val T) bool {
	var empty T
	if val == empty {
		return true
	}

	return false
}

func Coalesce[T comparable](in ...T) T {
	var empty T

	for _, v := range in {
		if v != empty {
			return v
		}
	}

	return empty
}

func DefaultValue[T any]() T {
	var t T
	return t
}
