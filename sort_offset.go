package sort

import (
  "golang.org/x/exp/constraints"
  "math/bits"
  u "unsafe"
)

// This code is modified go source.
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
//Redistribution and use in source and binary forms, with or without
//modification, are permitted provided that the following conditions are
//met:
//
//   * Redistributions of source code must retain the above copyright
//notice, this list of conditions and the following disclaimer.
//   * Redistributions in binary form must reproduce the above
//copyright notice, this list of conditions and the following disclaimer
//in the documentation and/or other materials provided with the
//distribution.
//   * Neither the name of Google Inc. nor the names of its
//contributors may be used to endorse or promote products derived from
//this software without specific prior written permission.
//
//THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

func StructField[Struct any, Field constraints.Ordered](base *Struct, field *Field, data []Struct) {
  offset := uintptr(u.Pointer(field)) - uintptr(u.Pointer(base))
  pdqsort_func_offset[Struct, Field](data, 0, len(data), bits.Len(uint(len(data))), int(offset))
}

// insertionSort_func sorts data[a:b] using insertion sort.
func insertionSort_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b int, offset int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && *(*T)(u.Add(u.Pointer(&data[j]), offset)) < *(*T)(u.Add(u.Pointer(&data[j-1]), offset)); j-- {
      data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

// siftDown_func implements the heap property on data[lo:hi].
// first is an offset into the array where the root of the heap lies.
func siftDown_func_offset[Struct any, T constraints.Ordered](data []Struct, lo, hi, first int, offset int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
    if child+1 < hi && *(*T)(u.Add(u.Pointer(&data[first+child]), offset)) < *(*T)(u.Add(u.Pointer(&data[first+child+1]), offset)) {
			child++
		}
    if !(*(*T)(u.Add(u.Pointer(&data[first+root]), offset)) < *(*T)(u.Add(u.Pointer(&data[first+child]), offset))) {
			return
		}
    data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func heapSort_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, offset int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown_func_offset[Struct, T](data, i, hi, first, offset)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
    data[first], data[first+i] = data[first+i], data[first]
		siftDown_func_offset[Struct, T](data, lo, i, first, offset)
	}
}

// pdqsort_func sorts data[a:b].
// The algorithm based on pattern-defeating quicksort(pdqsort), but without the optimizations from BlockQuicksort.
// pdqsort paper: https://arxiv.org/pdf/2106.05123.pdf
// C++ implementation: https://github.com/orlp/pdqsort
// Rust implementation: https://docs.rs/pdqsort/latest/pdqsort/
// limit is the number of allowed bad (very unbalanced) pivots before falling back to heapsort.
func pdqsort_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, limit, offset int) {
	const maxInsertion = 12

	var (
		wasBalanced    = true // whether the last partitioning was reasonably balanced
		wasPartitioned = true // whether the slice was already partitioned
	)

	for {
		length := b - a

		if length <= maxInsertion {
			insertionSort_func_offset[Struct, T](data, a, b, offset)
			return
		}

		// Fall back to heapsort if too many bad choices were made.
		if limit == 0 {
			heapSort_func_offset[Struct, T](data, a, b, offset)
			return
		}

		// If the last partitioning was imbalanced, we need to breaking patterns.
		if !wasBalanced {
			breakPatterns_func(data, a, b)
			limit--
		}

		pivot, hint := choosePivot_func_offset[Struct, T](data, a, b, offset)
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
			if partialInsertionSort_func_offset[Struct, T](data, a, b, offset) {
				return
			}
		}

		// Probably the slice contains many duplicate elements, partition the slice into
		// elements equal to and elements greater than the pivot.
		if a > 0 && !(*(*T)(u.Add(u.Pointer(&data[a-1]), offset)) < *(*T)(u.Add(u.Pointer(&data[pivot]), offset))) {
			mid := partitionEqual_func_offset[Struct, T](data, a, b, pivot, offset)
			a = mid
			continue
		}

		mid, alreadyPartitioned := partition_func_offset[Struct, T](data, a, b, pivot, offset)
		wasPartitioned = alreadyPartitioned

		leftLen, rightLen := mid-a, b-mid
		balanceThreshold := length / 8
		if leftLen < rightLen {
			wasBalanced = leftLen >= balanceThreshold
			pdqsort_func_offset[Struct, T](data, a, mid, limit, offset)
			a = mid + 1
		} else {
			wasBalanced = rightLen >= balanceThreshold
			pdqsort_func_offset[Struct, T](data, mid+1, b, limit, offset)
			b = mid
		}
	}
}

// partition_func does one quicksort partition.
// Let p = data[pivot]
// Moves elements in data[a:b] around, so that data[i]<p and data[j]>=p for i<newpivot and j>newpivot.
// On return, data[newpivot] = p
func partition_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, pivot, offset int) (newpivot int, alreadyPartitioned bool) {
  data[a], data[pivot] = data[pivot], data[a]
	i, j := a+1, b-1 // i and j are inclusive of the elements remaining to be partitioned

  for i <= j && *(*T)(u.Add(u.Pointer(&data[i]), offset)) < *(*T)(u.Add(u.Pointer(&data[a]), offset)) {
		i++
	}
	for i <= j && !(*(*T)(u.Add(u.Pointer(&data[j]), offset)) < *(*T)(u.Add(u.Pointer(&data[a]), offset))) {
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
		for i <= j && *(*T)(u.Add(u.Pointer(&data[i]), offset)) < *(*T)(u.Add(u.Pointer(&data[a]), offset)) {
			i++
		}
		for i <= j && !(*(*T)(u.Add(u.Pointer(&data[j]), offset)) < *(*T)(u.Add(u.Pointer(&data[a]), offset))) {
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
func partitionEqual_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, pivot, offset int) (newpivot int) {
  data[a], data[pivot] = data[pivot], data[a]
	i, j := a+1, b-1 // i and j are inclusive of the elements remaining to be partitioned

	for {
		for i <= j && !(*(*T)(u.Add(u.Pointer(&data[a]), offset)) < *(*T)(u.Add(u.Pointer(&data[i]), offset))) {
			i++
		}
		for i <= j && *(*T)(u.Add(u.Pointer(&data[a]), offset)) < *(*T)(u.Add(u.Pointer(&data[j]), offset)) {
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
func partialInsertionSort_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, offset int) bool {
	const (
		maxSteps         = 5  // maximum number of adjacent out-of-order pairs that will get shifted
		shortestShifting = 50 // don't shift any elements on short arrays
	)
	i := a + 1
	for j := 0; j < maxSteps; j++ {
		for i < b && !(*(*T)(u.Add(u.Pointer(&data[i]), offset)) < *(*T)(u.Add(u.Pointer(&data[i-1]), offset))) {
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
				if !(*(*T)(u.Add(u.Pointer(&data[j]), offset)) < *(*T)(u.Add(u.Pointer(&data[j-1]), offset))) {
					break
				}
        data[j], data[j-1] = data[j-1], data[j]
			}
		}
		// Shift the greater one to the right.
		if b-i >= 2 {
			for j := i + 1; j < b; j++ {
				if !(*(*T)(u.Add(u.Pointer(&data[j]), offset)) < *(*T)(u.Add(u.Pointer(&data[j-1]), offset))) {
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
func choosePivot_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, offset int) (pivot int, hint sortedHint) {
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
			i = medianAdjacent_func_offset[Struct, T](data, i, &swaps, offset)
			j = medianAdjacent_func_offset[Struct, T](data, j, &swaps, offset)
			k = medianAdjacent_func_offset[Struct, T](data, k, &swaps, offset)
		}
		// Find the median among i, j, k and stores it into j.
		j = median_func_offset[Struct, T](data, i, j, k, &swaps, offset)
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
func order2_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b int, swaps *int, offset int) (int, int) {
	if *(*T)(u.Add(u.Pointer(&data[b]), offset)) < *(*T)(u.Add(u.Pointer(&data[a]), offset)) {
		*swaps++
		return b, a
	}
	return a, b
}

// median_func returns x where data[x] is the median of data[a],data[b],data[c], where x is a, b, or c.
func median_func_offset[Struct any, T constraints.Ordered](data []Struct, a, b, c int, swaps *int, offset int) int {
	a, b = order2_func_offset[Struct, T](data, a, b, swaps, offset)
	b, c = order2_func_offset[Struct, T](data, b, c, swaps, offset)
	a, b = order2_func_offset[Struct, T](data, a, b, swaps, offset)
	return b
}

// medianAdjacent_func finds the median of data[a - 1], data[a], data[a + 1] and stores the index into a.
func medianAdjacent_func_offset[Struct any, T constraints.Ordered](data []Struct, a int, swaps *int, offset int) int {
	return median_func_offset[Struct, T](data, a-1, a, a+1, swaps, offset)
}
