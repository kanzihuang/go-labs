package sort_lab

import (
	"math/rand"
	"sort"
	"testing"
)

func testSortBy(t *testing.T, data []int, sortName string, sortFunc func(data []int)) {
	got := make([]int, len(data))
	copy(got, data)
	sortFunc(got)
	if !sort.IntsAreSorted(got) {
		t.Errorf("%s(%v) = %v, not sorted\n", sortName, data, got)
	}
}
func TestSort(t *testing.T) {
	cases := [][]int{
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
		{3, 7, 2, 7, 11},
		{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586},
	}
	for _, data := range cases {
		testSortBy(t, data, "insertionSort", insertionSort)
		testSortBy(t, data, "quickSort", quickSort)
	}
}

func TestSortGreat(t *testing.T) {
	num := 100000
	data := generateRands(num)
	sortedData := generateRands(num)
	sort.Ints(sortedData)
	reversedData := generateRands(num)
	sort.Sort(sort.Reverse(sort.IntSlice(reversedData)))

	tests := []struct {
		name string
		sort func([]int)
	}{
		{name: "insertionSort", sort: insertionSort},
		{name: "quickSort", sort: quickSort},
		{name: "sort.Ints", sort: sort.Ints},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSortBy(t, data, tt.name, tt.sort)
		})

		t.Run(tt.name+"_sorted", func(t *testing.T) {
			testSortBy(t, sortedData, tt.name, tt.sort)
		})

		t.Run(tt.name+"_reversed", func(t *testing.T) {
			testSortBy(t, reversedData, tt.name, tt.sort)
		})
	}
}

func generateRands(num int) []int {
	data := make([]int, num)
	for i := range data {
		data[i] = rand.Int()
	}
	return data
}
