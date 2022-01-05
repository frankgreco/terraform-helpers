package overlap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntSlice(t *testing.T) {
	for _, test := range []struct {
		name  string
		slice []int
		err   string
	}{
		{
			name:  "no overlapping happy path",
			slice: []int{1, 2, 3, 4},
		},
		{
			name:  "overlapping happy path (1/2)",
			slice: []int{1, 1, 3, 4},
			err:   "The element 1 is supplied by more than one range.",
		},
		{
			name:  "overlapping happy path (2/2)",
			slice: []int{1, 2, 3, 4, 4},
			err:   "The element 4 is supplied by more than one range.",
		},
		{
			name:  "overlapping happy path (start with unsorted)",
			slice: []int{4, 3, 3, 1},
			err:   "The element 3 is supplied by more than one range.",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := IntSlice(test.slice); test.err == "" {
				require.NoError(t, err, test.name)
			} else {
				require.NotNil(t, err, test.name)
				require.Equal(t, test.err, err.Error(), test.name)
			}
		})
	}
}
