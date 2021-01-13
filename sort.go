package scoretable

import "sort"

type lessFunc func(p1, p2 *Score) bool

var (
	NBW = func(s1, s2 *Score) bool {
		return s1.NBW > s2.NBW
	}

	SOS = func(s1, s2 *Score) bool {
		return s1.SOS > s2.SOS
	}

	SOSOS = func(s1, s2 *Score) bool {
		return s1.SOSOS > s2.SOSOS
	}

	PlayerId = func(s1, s2 *Score) bool {
		return s1.PlayerId < s2.PlayerId
	}
)

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	scores Scores
	less   []lessFunc
}

func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

func (ms *multiSorter) Len() int {
	return len(ms.scores)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.scores[i], ms.scores[j] = ms.scores[j], ms.scores[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.scores[i], ms.scores[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		if less(p, q) {
			return true
		} else if less(q, p) {
			return false
		}
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(scores Scores) {
	ms.scores = scores
	sort.Sort(ms)
}
