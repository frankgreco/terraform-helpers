package overlap

import (
	"github.com/frankgreco/terraform-helpers/internal/utils"
)

func OrderedPairs(items []utils.OrderedPair) error {
	if len(items) < 2 {
		return nil
	}
	return utils.Overlaps(items)
}
