package bivariate

// Order function.
//
// Return value is meant to be interpreted in the following way:
// -1: deg1<deg2; 0: deg1==deg2; 1: deg1>deg2
type Order func(deg1, deg2 [2]uint) int

// degCompare is an internal function used for the (weighted) degree (reverse)
// lexicographical orderings.
func degCompare(xWeight, yWeight uint, tiebreak Order) Order {
	return func(deg1, deg2 [2]uint) (out int) {
		switch {
		case deg1 == deg2:
			out = 0
		case deg1[0]*xWeight+deg1[1]*yWeight > deg2[0]*xWeight+deg2[1]*yWeight:
			out = 1
		case deg1[0]*xWeight+deg1[1]*yWeight < deg2[0]*xWeight+deg2[1]*yWeight:
			out = -1
		case deg1[0]*xWeight+deg1[1]*yWeight == deg2[0]*xWeight+deg2[1]*yWeight:
			// Fall back to lexicographical ordering
			out = tiebreak(deg1, deg2)
		}
		return out
	}
}

// WDegLex returns the weighted degree lexicographical ordering.
//
// The resulting order will break ties using the lexicographical ordering. The
// boolean xGtY indicates whether X is greater than Y.
func WDegLex(xWeight, yWeight uint, xGtY bool) Order {
	return degCompare(
		xWeight,
		yWeight,
		Lex(xGtY),
	)
}

// WDegRevLex returns the weighted degree reverse lexicographical ordering.
//
// The resulting order will break ties using the lexicographical ordering, but
// the smaller degree with respect to the lexicographical ordering is considered
// as the larger degree with respect to the reverse ordering. In addition, the
// order of the variables is reversed when performing the lexicographical
// comparison.
//
// The boolean xGtY indicates whether X is greater than Y.
func WDegRevLex(xWeight, yWeight uint, xGtY bool) Order {
	return degCompare(
		xWeight,
		yWeight,
		func(deg1, deg2 [2]uint) int {
			return -1 * Lex(!xGtY)(deg1, deg2)
		},
	)
}

// DegLex returns the total degree ordering.
//
// The resulting order will break ties using the lexicographical ordering. The
// boolean xGtY indicates whether X is greater than Y.
func DegLex(xGtY bool) Order {
	return WDegLex(1, 1, xGtY)
}

// DegRevLex returns the total degree reverse ordering.
//
// The resulting order will break ties using the lexicographical ordering, but
// the smaller degree with respect to the lexicographical ordering is considered
// as the larger degree with respect to the reverse ordering. In addition, the
// order of the variables is reversed when performing the lexicographical
// comparison.
//
// The boolean xGtY indicates whether X is greater than Y.
func DegRevLex(xGtY bool) Order {
	return WDegRevLex(1, 1, xGtY)
}

// Lex returns the lexicographical ordering.
//
// xGtY indicates whether X is greater than Y.
func Lex(xGtY bool) Order {
	f := func(deg1, deg2 [2]uint) (out int) {
		switch {
		case deg1 == deg2:
			out = 0
		case deg1[0] > deg2[0]:
			out = 1
		case deg1[0] < deg2[0]:
			out = -1
		case deg1[0] == deg2[0] && deg1[1] > deg2[1]:
			out = 1
		case deg1[0] == deg2[0] && deg1[1] < deg2[1]:
			out = -1
		}
		return out
	}
	if xGtY {
		return f
	}
	return func(deg1, deg2 [2]uint) int { return f(swap(deg1), swap(deg2)) }
}

func swap(deg [2]uint) [2]uint {
	return [2]uint{deg[1], deg[0]}
}

/* Copyright 2019 René Bødker Christensen
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
