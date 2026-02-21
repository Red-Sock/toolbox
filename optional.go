package toolbox

import (
	"bytes"
	"encoding/json"
)

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

func (o *Optional[T]) UnmarshalJSON(v []byte) error {
	if bytes.Equal(bytes.TrimSpace(v), []byte("null")) {
		var zero T
		o.Value = zero
		o.Valid = false
		return nil
	}

	err := json.Unmarshal(v, &o.Value)

	o.Valid = err == nil

	return err
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(o.Value)
}
