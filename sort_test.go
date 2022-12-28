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
const BenchIters = 1000

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
  t.Logf("Fast implementation took %v", avgFastSorterDur)
  avgFastLessSorterDur := getAvgSpeed[int](func(data []int) {
    fastsort.LessSort[int](data, func(a, b *int) bool { return *a < *b })
  }, rand.Int)
  t.Logf("Fast LessSort implementation took %v", avgFastLessSorterDur)
  avgDefaultSorterDur := getAvgSpeed[int](sort.Ints, rand.Int)
  t.Logf("Default implementation took %v", avgDefaultSorterDur)

  if avgFastSorterDur >= avgDefaultSorterDur {
    t.Errorf("Default implementation(%v) should be slower than generic implementation(%v).", avgDefaultSorterDur, avgFastSorterDur)
  }

  if avgFastSorterDur >= avgFastLessSorterDur {
    t.Errorf("Ordered implementation(%v) should be faster than less implementation(%v).", avgFastSorterDur, avgFastLessSorterDur)
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

type Person struct {
  Name string
  Age int
}

func PersonLess(a, b *Person) bool {
  return a.Age < b.Age
}

func RandomPerson() Person {
  return Person{ Age: rand.Int() }
}

func TestSortLessCapability(t *testing.T) {
  data := NewRandArray(InputSize, RandomPerson)
  fastsort.LessSort(data, PersonLess)
  if !sort.SliceIsSorted(data, func(i, j int) bool {
    return PersonLess(&data[i], &data[j])
  }) {
    t.Errorf("Failed sorting array!")
  }
}

func TestLessSortSpeed(t *testing.T) {
  avgFastSorterDur := getAvgSpeed[Person](func(data []Person) {
    fastsort.LessSort(data, PersonLess)
  }, RandomPerson)
  t.Logf("Less sort implementation took %v", avgFastSorterDur)

  avgFieldSorterDur := getAvgSpeed[Person](func(data []Person) {
    proto := Person{}
    fastsort.StructField[Person, int](&proto, &proto.Age, data)
  }, RandomPerson)
  t.Logf("Struct field implementation took %v", avgFieldSorterDur)

  avgDefaultSorterDur := getAvgSpeed[Person](func(data []Person) {
    sort.Slice(data, func(i, j int) bool {
      return PersonLess(&data[i], &data[j])
    })
  }, RandomPerson)
  t.Logf("Default implementation took %v", avgDefaultSorterDur)

  if avgFastSorterDur >= avgDefaultSorterDur {
    t.Errorf("Default implementation(%v) should be slower than generic implementation(%v).", avgDefaultSorterDur, avgFastSorterDur)
  }

  if avgFieldSorterDur >= avgFastSorterDur {
    t.Errorf("Less function implementation(%v) should be slower than struct field implementation(%v).", avgFastSorterDur, avgFieldSorterDur)
  }
}

