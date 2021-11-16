package generics

import (
	"github.com/cheekybits/genny/generic"
	uuid2 "github.com/satori/go.uuid"
)

type ItemType generic.Type
type ItemTypeSlice []ItemType
type ItemTypeSlicePtrs []*ItemType

var (
	_ uuid2.UUID
)

func (slice ItemTypeSlice) Unique() (out ItemTypeSlice) {
	umap := map[ItemType]int{}
	for _, item := range slice {
		umap[item]++
	}

	for item, _ := range umap {
		out = append(out, item)
	}
	return out
}

func (slice ItemTypeSlice) UniqueBy(keyer func(ItemType) interface{}) ItemTypeSlice {
	u := map[interface{}]ItemType{}

	for _, item := range slice {
		key := keyer(item)
		u[key] = item
	}

	out := []ItemType{}
	for _, item := range u {
		out = append(out, item)
	}
	return out
}

func (slice ItemTypeSlice) FilterBy(predicate func(ItemType) bool) ItemTypeSlice {
	out := []ItemType{}
	for _, item := range slice {
		if predicate(item) {
			out = append(out, item)
		}
	}

	return out
}

func (slice ItemTypeSlice) Contains(t ItemType) bool {
	for _, item := range slice {
		if item == t {
			return true
		}
	}

	return false
}

func (slice ItemTypeSlice) IsEmpty() bool {
	return len(slice) == 0
}

// from https://stackoverflow.com/a/36000696
func (slice ItemTypeSlice) SameAs(other ItemTypeSlice) bool {
	if len(slice) != len(other) {
		return false
	}
	// create a map of string -> int
	diff := make(map[ItemType]int, len(slice))
	for _, _x := range slice {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}

	for _, _y := range other {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}

	return len(diff) == 0
}

//
// Slice of pointers
//

func (slice ItemTypeSlicePtrs) UniqueBy(keyer func(*ItemType) interface{}) ItemTypeSlicePtrs {
	u := map[interface{}]*ItemType{}

	for _, item := range slice {
		key := keyer(item)
		u[key] = item
	}

	out := []*ItemType{}
	for _, item := range u {
		out = append(out, item)
	}
	return out
}

func (slice ItemTypeSlicePtrs) FilterBy(predicate func(*ItemType) bool) ItemTypeSlicePtrs {
	out := []*ItemType{}
	for _, item := range slice {
		if predicate(item) {
			out = append(out, item)
		}
	}

	return out
}

func (slice ItemTypeSlicePtrs) Contains(t *ItemType) bool {
	for _, item := range slice {
		if item == t {
			return true
		}
	}

	return false
}

func (slice ItemTypeSlicePtrs) FilterNotNil() ItemTypeSlicePtrs {
	return slice.FilterBy(func(t *ItemType) bool {
		return t != nil
	})
}

func (slice ItemTypeSlicePtrs) IsEmpty() bool {
	return len(slice) == 0
}

// from https://stackoverflow.com/a/36000696
func (slice ItemTypeSlicePtrs) SameAs(other ItemTypeSlicePtrs) bool {
	if len(slice) != len(other) {
		return false
	}
	// create a map of string -> int
	diff := make(map[*ItemType]int, len(slice))
	for _, _x := range slice {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}

	for _, _y := range other {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}

	return len(diff) == 0
}
