package overlap

import (
	"strconv"

	"github.com/frankgreco/terraform-helpers/internal/utils"
)

// IntSlice ensures there is no overlap.
// Of course there is an easier way to do this
// But i'm utilizing an existing helper.
func IntSlice(items []int) error {
	if len(items) < 2 {
		return nil
	}

	var pairs []utils.OrderedPair
	{
		for _, item := range items {
			pairs = append(pairs, utils.NewOrderedPair(strconv.Itoa(item), strconv.Itoa(item)))
		}
	}

	return utils.Overlaps(pairs)
}
