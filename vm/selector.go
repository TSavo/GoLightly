package vm

import (
	"sort"
)

type Selector interface {
	Select(*SolutionList) *SolutionList
}

type MultiSelector []Selector

func NewMultiSelector(selectors... Selector) MultiSelector {
	return MultiSelector(selectors)
}

func (multi MultiSelector) Select(s *SolutionList) *SolutionList {
	solutions := make(SolutionList, 0)
	for _, x := range multi {
		solutions = append(solutions, *(x).Select(s)...)
	}
	return &solutions
}

func All(selectors... Selector) MultiSelector {
	return MultiSelector(selectors)
}

func (multi *MultiSelector) AddSelector(s *Selector) {
	(*multi) = append(*multi, *s)
}

type TopXSelector struct {
	Keep int
}

func (topx TopXSelector) Select(s *SolutionList) *SolutionList {
	sort.Sort(s)
	x := (*s)[:topx.Keep]
	return &x
}

type StochasticUniversalSelector struct {
	Keep int
}

func NewStochasticUniversalSelector(keep int) *StochasticUniversalSelector {
	return &StochasticUniversalSelector{keep}
}

func (sel StochasticUniversalSelector) Select(s *SolutionList) *SolutionList {
	sort.Sort(sort.Reverse(*s))
	f := int64(0)
	for _, x := range *s {
		f += x.Reward
	}
	n := sel.Keep
	p := int(f / int64(n))
	start := rng.Int()%p + 1
	pointers := make([]int, n)
	for i := 0; i < n; i++ {
		pointers[i] = start + i*p
	}
	ret := RWS(*s, pointers)
	return &ret
}

func RWS(solutions SolutionList, pointers []int) SolutionList {
	keep := make(SolutionList, len(pointers))
	i := 0
	for _, p := range pointers {
		for int(solutions[i].Reward) < p {
			i++
		}
		keep = append(keep, solutions[i])
	}
	return keep
}

type TournamentSelector struct {
	Keep int
}

func (t TournamentSelector) Select(solutions *SolutionList) *SolutionList {
	keepers := make(SolutionList, t.Keep)
	for x := 0; x < t.Keep; x++ {
		keepers = append(keepers, Tournament((*solutions)[rng.Int()%len(*solutions)], (*solutions)[rng.Int()%len(*solutions)]))
	}
	return &keepers
}

func Tournament(warrior1 *Solution, warrior2 *Solution) *Solution {
	var highest, lowest *Solution
	if warrior1.Reward >= warrior2.Reward {
		highest, lowest = warrior1, warrior2
	} else {
		highest, lowest = warrior2, warrior1
	}
	if rng.Int63()%highest.Reward > lowest.Reward/2 {
		return highest
	} else {
		return lowest
	}
}
