package compare

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LongStrings can be used to see changes in JetBrain's comparator as it validates file over 800 bytes frames
func LongStrings(t *testing.T, expected, actual []byte) {
	expectedReader := bytes.NewReader(expected)
	actualReader := bytes.NewReader(actual)
	for {
		expectedSlice := make([]byte, 800)
		actualSlice := make([]byte, 800)

		expLen, expErr := expectedReader.Read(expectedSlice)
		actLen, actErr := actualReader.Read(actualSlice)

		if expErr == io.EOF && actErr == io.EOF {
			return
		}

		expectedSlice = expectedSlice[:expLen]
		actualSlice = actualSlice[:actLen]

		require.NoError(t, actErr)
		require.NoError(t, expErr)

		assert.Equal(t, string(expectedSlice), string(actualSlice))
	}
}
