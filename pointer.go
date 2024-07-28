package toolbox

func ToPtr[T any](in T) *T {
	return &in
}

func FromPtr[T any](in *T) T {
	if in == nil {
		var empty T
		return empty
	}

	return *in
}
