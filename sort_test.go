package main

import (
	"log"
	"sort"
	"testing"
)

// type sort.Interface interface {
// 	Len() int
// 	Less(i, j int) bool
// 	Swap(i, j int)
// }

// https://golang.org/pkg/sort/
func TestSort(t *testing.T) {
	s := []string{"a", "z", "b"}
	f := []float64{0.1, 100.1, -0.9}
	// destructive!
	sort.Strings(s)
	sort.Float64s(f)

	i := []int{10, 99, -1}
	log.Println(i)
	sort.Ints(i)
	log.Println(i)
	// sort.IntSlice(i) returns []int implementing sort.Interface
	// sort.Reverse returns the reverse order for sort.Interface
	sort.Sort(sort.Reverse(sort.IntSlice(i)))
	log.Println(i)
}

type Data struct {
	I int
	J int
}
type DataSlice []Data

func (d DataSlice) Len() int           { return len(d) }
func (d DataSlice) Less(i, j int) bool { return d[i].I*d[i].J < d[j].I*d[j].J }
func (d DataSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func TestMyData(t *testing.T) {
	d := []Data{
		Data{9, 10},
		Data{2, 5},
		Data{10, 9},
	}
	sort.Sort(DataSlice(d))
	log.Println(d)
}
