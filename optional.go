package toolbox

type Optional[T any] struct {
	Value T
	Valid bool
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		Value: value,
		Valid: true,
	}
}
