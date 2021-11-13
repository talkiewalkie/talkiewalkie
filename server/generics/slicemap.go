package generics

import "github.com/cheekybits/genny/generic"

type MapTarget generic.Type

func (slice ItemTypeSlice) MapToMapTarget(f func(ItemType) MapTarget) []MapTarget {
	out := []MapTarget{}
	for _, item := range slice {
		out = append(out, f(item))
	}

	return out
}

func (slice ItemTypeSlicePtrs) MapToMapTarget(f func(*ItemType) MapTarget) []MapTarget {
	out := []MapTarget{}
	for _, item := range slice {
		out = append(out, f(item))
	}

	return out
}
