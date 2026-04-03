package terminal

import "strings"

func limitedItemsAsString[T any](items []T, itemToString func(T) string) string {
	return strings.Join(func() []string {
		const limit = 3
		keys := make([]string, 0, limit)
		more, limitedItems := limitSlice(items, limit)
		for _, item := range limitedItems {
			keys = append(keys, itemToString(item))
		}
		if more {
			keys = append(keys, "...")
		}
		return keys
	}(), "\n")
}

func limitSlice[T any](slice []T, limit int) (more bool, items []T) {
	if len(slice) > limit {
		return true, slice[:limit]
	}
	return false, slice
}
