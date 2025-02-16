package utilities

import (
	"strconv"
	"strings"
)

func SplitBillingDetailIds(ids string) []int {
	if ids == "" {
		return []int{}
	}

	strIds := strings.Split(ids, ",")
	intIds := make([]int, 0, len(strIds))

	for _, id := range strIds {
		if trimmedID := strings.TrimSpace(id); trimmedID != "" {
			if intID, err := strconv.Atoi(trimmedID); err == nil {
				intIds = append(intIds, intID)
			}
		}
	}

	return intIds
}
