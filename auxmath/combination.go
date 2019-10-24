package auxmath

// CombinIter is an iterator for combinations.
//
// It will iterate over all possible ways to choose a given number of elements
// from n elements. The combinations are represented by sorted slices of
// indices, and the iterator will produce them in lexicographically increasing
// order. For instance, it will generate the sequence [0,1,2], [0,1,3],...,
// [3,4,5] if defined with n=6 and k=3.
type CombinIter struct {
	n     int
	slice []int
	atEnd bool
}

// NewCombinIter return a new iterator for combinations.
func NewCombinIter(n, k int) *CombinIter {
	s := make([]int, k, k)
	for i := range s {
		s[i] = i
	}
	return &CombinIter{
		n:     n,
		slice: s,
	}
}

// Current returns the current combination.
func (ci *CombinIter) Current() []int {
	return ci.slice
}

// Active returns a boolean describing whether all combinations have been
// considered.
func (ci *CombinIter) Active() bool {
	return !ci.atEnd
}

// Next increments the iterator to the next combination with respect to the
// lexicographical ordering.
func (ci *CombinIter) Next() {
	if ci.atEnd {
		return
	}
	for i := range ci.slice {
		// Step back through the slice...
		j := len(ci.slice) - 1 - i
		if ci.slice[j] < (ci.n - i - 1) {
			// ...until we find an index that can be incremented
			ci.slice[j]++
			for l := 1; l <= i; l++ {
				// Update the following combinations to preserve the
				// sorting and the lexicographical ordering
				ci.slice[j+l] = ci.slice[j] + l
			}
			return
		}
	}
	ci.atEnd = true
}