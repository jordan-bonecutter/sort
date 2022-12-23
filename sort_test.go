package sort_test

import (
  "testing"
  fastsort "github.com/jordan-bonecutter/sort"
  "math/rand"
  "runtime"
  "sort"
  "time"
)

const InputSize = 10000
const BenchIters = 100

func NewRandArray[T any](length int, rvg func() T) []T {
  ret := make([]T, length)
  for idx := 0; idx < length; idx++ {
    ret[idx] = rvg()
  }
  return ret
}

func TestSortCapability(t *testing.T) {
  data := NewRandArray(InputSize, rand.Int)
  fastsort.Ordered(data)
  if !sort.IntsAreSorted(data) {
    t.Errorf("Failed sorting array!")
  }
}

func TestSortSpeed(t *testing.T) {
  avgFastSorterDur := getAvgSpeed[int](fastsort.Ordered[int], rand.Int)
  avgDefaultSorterDur := getAvgSpeed[int](sort.Ints, rand.Int)

  if avgFastSorterDur >= avgDefaultSorterDur {
    t.Errorf("Default implementation(%v) should be slower than generic implementation(%v).", avgDefaultSorterDur, avgFastSorterDur)
  }
}

func getAvgSpeed[T any](sorter func([]T), rvg func() T) time.Duration {
  var total time.Duration

  for iter := 0; iter < BenchIters; iter++ {
    data := NewRandArray(InputSize, rvg)
    runtime.GC()

    start := time.Now()
    sorter(data)
    end := time.Now()
    total += end.Sub(start)
  }

  return total / BenchIters
}

