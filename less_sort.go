package sort

import (
  "math/bits"
)

func LessSort[T any](data []T, less func(a, b *T) bool) {
  pdqsort_func_less(data, 0, len(data), bits.Len(uint(len(data))), less)
}

// insertionSort_func sorts data[a:b] using insertion sort.
func insertionSort_func_less[T any](data []T, a, b int, less func(a, b *T) bool) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(&data[j], &data[j-1]); j-- {
      data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

// siftDown_func implements the heap property on data[lo:hi].
// first is an offset into the array where the root of the heap lies.
func siftDown_func_less[T any](data []T, lo, hi, first int, less func(a, b *T) bool) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(&data[first+child], &data[first+child+1]) {
			child++
		}
		if !less(&data[first+root], &data[first+child]) {
			return
		}
    data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func heapSort_func_less[T any](data []T, a, b int, less func(a, b *T) bool) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown_func_less(data, i, hi, first, less)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
    data[first], data[first+i] = data[first+i], data[first]
		siftDown_func_less(data, lo, i, first, less)
	}
}

// pdqsort_func sorts data[a:b].
// The algorithm based on pattern-defeating quicksort(pdqsort), but without the optimizations from BlockQuicksort.
// pdqsort paper: https://arxiv.org/pdf/2106.05123.pdf
// C++ implementation: https://github.com/orlp/pdqsort
// Rust implementation: https://docs.rs/pdqsort/latest/pdqsort/
// limit is the number of allowed bad (very unbalanced) pivots before falling back to heapsort.
func pdqsort_func_less[T any](data []T, a, b, limit int, less func(a, b *T) bool) {
	const maxInsertion = 12

	var (
		wasBalanced    = true // whether the last partitioning was reasonably balanced
		wasPartitioned = true // whether the slice was already partitioned
	)

	for {
		length := b - a

		if length <= maxInsertion {
			insertionSort_func_less(data, a, b, less)
			return
		}

		// Fall back to heapsort if too many bad choices were made.
		if limit == 0 {
			heapSort_func_less(data, a, b, less)
			return
		}

		// If the last partitioning was imbalanced, we need to breaking patterns.
		if !wasBalanced {
			breakPatterns_func(data, a, b)
			limit--
		}

		pivot, hint := choosePivot_func_less(data, a, b, less)
		if hint == decreasingHint {
			reverseRange_func(data, a, b)
			// The chosen pivot was pivot-a elements after the start of the array.
			// After reversing it is pivot-a elements before the end of the array.
			// The idea came from Rust's implementation.
			pivot = (b - 1) - (pivot - a)
			hint = increasingHint
		}

		// The slice is likely already sorted.
		if wasBalanced && wasPartitioned && hint == increasingHint {
			if partialInsertionSort_func_less(data, a, b, less) {
				return
			}
		}

		// Probably the slice contains many duplicate elements, partition the slice into
		// elements equal to and elements greater than the pivot.
		if a > 0 && !less(&data[a-1], &data[pivot]) {
			mid := partitionEqual_func_less(data, a, b, pivot, less)
			a = mid
			continue
		}

		mid, alreadyPartitioned := partition_func_less(data, a, b, pivot, less)
		wasPartitioned = alreadyPartitioned

		leftLen, rightLen := mid-a, b-mid
		balanceThreshold := length / 8
		if leftLen < rightLen {
			wasBalanced = leftLen >= balanceThreshold
			pdqsort_func_less(data, a, mid, limit, less)
			a = mid + 1
		} else {
			wasBalanced = rightLen >= balanceThreshold
			pdqsort_func_less(data, mid+1, b, limit, less)
			b = mid
		}
	}
}

// partition_func does one quicksort partition.
// Let p = data[pivot]
// Moves elements in data[a:b] around, so that data[i]<p and data[j]>=p for i<newpivot and j>newpivot.
// On return, data[newpivot] = p
func partition_func_less[T any](data []T, a, b, pivot int, less func(a, b *T) bool) (newpivot int, alreadyPartitioned bool) {
  data[a], data[pivot] = data[pivot], data[a]
	i, j := a+1, b-1 // i and j are inclusive of the elements remaining to be partitioned

	for i <= j && less(&data[i], &data[a]) {
		i++
	}
	for i <= j && !less(&data[j], &data[a]) {
		j--
	}
	if i > j {
    data[j], data[a] = data[a], data[j]
		return j, true
	}
  data[i], data[j] = data[j], data[i]
	i++
	j--

	for {
		for i <= j && less(&data[i], &data[a]) {
			i++
		}
		for i <= j && !less(&data[j], &data[a]) {
			j--
		}
		if i > j {
			break
		}
    data[i], data[j] = data[j], data[i]
		i++
		j--
	}
  data[j], data[a] = data[a], data[j]
	return j, false
}

// partitionEqual_func partitions data[a:b] into elements equal to data[pivot] followed by elements greater than data[pivot].
// It assumed that data[a:b] does not contain elements smaller than the data[pivot].
func partitionEqual_func_less[T any](data []T, a, b, pivot int, less func(a, b *T) bool) (newpivot int) {
  data[a], data[pivot] = data[pivot], data[a]
	i, j := a+1, b-1 // i and j are inclusive of the elements remaining to be partitioned

	for {
		for i <= j && !less(&data[a], &data[i]) {
			i++
		}
		for i <= j && less(&data[a], &data[j]) {
			j--
		}
		if i > j {
			break
		}
    data[i], data[j] = data[j], data[i]
		i++
		j--
	}
	return i
}

// partialInsertionSort_func partially sorts a slice, returns true if the slice is sorted at the end.
func partialInsertionSort_func_less[T any](data []T, a, b int, less func(a, b *T) bool) bool {
	const (
		maxSteps         = 5  // maximum number of adjacent out-of-order pairs that will get shifted
		shortestShifting = 50 // don't shift any elements on short arrays
	)
	i := a + 1
	for j := 0; j < maxSteps; j++ {
		for i < b && !less(&data[i], &data[i-1]) {
			i++
		}

		if i == b {
			return true
		}

		if b-a < shortestShifting {
			return false
		}

    data[i], data[i-1] = data[i-1], data[i]

		// Shift the smaller one to the left.
		if i-a >= 2 {
			for j := i - 1; j >= 1; j-- {
				if !less(&data[j], &data[j-1]) {
					break
				}
        data[j], data[j-1] = data[j-1], data[j]
			}
		}
		// Shift the greater one to the right.
		if b-i >= 2 {
			for j := i + 1; j < b; j++ {
				if !less(&data[j], &data[j-1]) {
					break
				}
        data[j], data[j-1] = data[j-1], data[j]
			}
		}
	}
	return false
}

// choosePivot_func chooses a pivot in data[a:b].
//
// [0,8): chooses a static pivot.
// [8,shortestNinther): uses the simple median-of-three method.
// [shortestNinther,∞): uses the Tukey ninther method.
func choosePivot_func_less[T any](data []T, a, b int, less func(a, b *T) bool) (pivot int, hint sortedHint) {
	const (
		shortestNinther = 50
		maxSwaps        = 4 * 3
	)

	l := b - a

	var (
		swaps int
		i     = a + l/4*1
		j     = a + l/4*2
		k     = a + l/4*3
	)

	if l >= 8 {
		if l >= shortestNinther {
			// Tukey ninther method, the idea came from Rust's implementation.
			i = medianAdjacent_func_less(data, i, &swaps, less)
			j = medianAdjacent_func_less(data, j, &swaps, less)
			k = medianAdjacent_func_less(data, k, &swaps, less)
		}
		// Find the median among i, j, k and stores it into j.
		j = median_func_less(data, i, j, k, &swaps, less)
	}

	switch swaps {
	case 0:
		return j, increasingHint
	case maxSwaps:
		return j, decreasingHint
	default:
		return j, unknownHint
	}
}

// order2_func returns x,y where data[x] <= data[y], where x,y=a,b or x,y=b,a.
func order2_func_less[T any](data []T, a, b int, swaps *int, less func(a, b *T) bool) (int, int) {
	if less(&data[b], &data[a]) {
		*swaps++
		return b, a
	}
	return a, b
}

// median_func returns x where data[x] is the median of data[a],data[b],data[c], where x is a, b, or c.
func median_func_less[T any](data []T, a, b, c int, swaps *int, less func(a, b *T) bool) int {
	a, b = order2_func_less(data, a, b, swaps, less)
	b, c = order2_func_less(data, b, c, swaps, less)
	a, b = order2_func_less(data, a, b, swaps, less)
	return b
}

// medianAdjacent_func finds the median of data[a - 1], data[a], data[a + 1] and stores the index into a.
func medianAdjacent_func_less[T any](data []T, a int, swaps *int, less func(a, b *T) bool) int {
	return median_func_less(data, a-1, a, a+1, swaps, less)
}

