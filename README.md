[![Go Reference](https://pkg.go.dev/badge/github.com/jordan-bonecutter/sort.svg)](https://pkg.go.dev/github.com/jordan-bonecutter/sort)

# sort
A faster implementation of Go's sort package for slices.

# How
Go now has generics, but the sort implementation uses a `lessSwap` wrapper struct with the following form:

```go
// lessSwap is a pair of Less and Swap function for use with the
// auto-generated func-optimized variant of sort.go in
// zfuncversion.go.
type lessSwap struct {
	Less func(i, j int) bool
	Swap func(i, j int)
}
```

Instead of directy operating on slices, swap (`data[i], data[j] = data[j], data[i]`) and less (`data[i] < data[j]`) are passed as function literals.

This means each swap, instead of being a quick array access, is a function call then a swap. A similar statement can be made
regarding `Less`. We make this faster be implementing these operations directly on slices, saving us time in the pointer deref and the function call!
