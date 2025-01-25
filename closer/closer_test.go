package closer

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CloserHTTP(t *testing.T) {
	closer := Closer{}
	closer.Add(func() error {
		return nil
	})

	closer.Add(func() error {
		return errors.New("some error")
	})

	err := closer.Close()
	require.Error(t, err)

	w := httptest.NewRecorder()

	closer.ServeHTTP(w, nil)
	body, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	var resp []State
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)

	expected := []State{
		{
			Name:     "go.redsock.ru/toolbox/closer.Test_CloserHTTP.func1",
			IsClosed: true,
		},
		{
			Name:  "go.redsock.ru/toolbox/closer.Test_CloserHTTP.func2",
			Error: "some error",
		},
	}

	require.Equal(t, expected, resp)

}
