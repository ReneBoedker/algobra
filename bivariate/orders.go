package bivariate

// Order function. Return value is meant to be interpreted as:
// -1: deg1<deg2; 0: deg1==deg2; 1: deg1>deg2
type order func(deg1, deg2 [2]uint) int

func WDegLex(xWeight, yWeight uint) order {
	return func(deg1, deg2 [2]uint) int {
		switch {
		case deg1 == deg2:
			return 0
		case deg1[0]*xWeight+deg1[1]*yWeight > deg2[0]*xWeight+deg2[1]*yWeight:
			return 1
		case deg1[0]*xWeight+deg1[1]*yWeight < deg2[0]*xWeight+deg2[1]*yWeight:
			return -1
		case deg1[0]*xWeight+deg1[1]*yWeight == deg2[0]*xWeight+deg2[1]*yWeight:
			if deg1[1] > deg2[1] {
				return 1
			}
			return -1
		default:
			panic("WDegLex: Comparison failed")
		}
	}
}

func Lex(xGtY bool) order {
	f := func(deg1, deg2 [2]uint) int {
		switch {
		case deg1 == deg2:
			return 0
		case deg1[0] > deg2[0]:
			return 1
		case deg1[0] < deg2[0]:
			return -1
		case deg1[0] == deg2[0] && deg1[1] > deg2[1]:
			return 1
		case deg1[0] == deg2[0] && deg1[1] < deg2[1]:
			return -1
		default:
			panic("Lex: Comparison failed")
		}
	}
	if xGtY {
		return f
	}
	return func(deg1, deg2 [2]uint) int { return -1 * f(deg1, deg2) }
}
