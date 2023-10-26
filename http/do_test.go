package http

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	ar := []int{200, 1, 101, 3, 0, 2, 0, 6, 4}
	sort.Slice(ar, func(i, j int) bool {
		if ar[i] > 100 {
			return false
		}
		if ar[j] > 100 {
			return true
		}
		if ar[i] == 0 {
			return false
		}
		if ar[j] == 0 {
			return true
		}
		return ar[i] < ar[j]
	})
	t.Log(ar)
}
