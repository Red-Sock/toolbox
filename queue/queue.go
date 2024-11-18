package queue

type Queue[T any] struct {
	value *T

	prev *Queue[T]
}
