package pkg

import (
	"strings"

	"golang.org/x/exp/slices"
)

type TargetSetter func(string, *Target)

type Target struct {
	Name   string
	Health int

	set TargetSetter

	// A list of possible targets, e.g. "manticore".
	candidates []string

	// A list of NPCs present, e.g. "a ferocious manticore".
	present []string
}

func NewTarget(set TargetSetter) *Target {
	return &Target{set: set}
}

func (tgt *Target) Set(name string) {
	if name == tgt.Name {
		return
	}

	tgt.set(name, tgt)
}

func (tgt *Target) SetCandidates(names []string) {
	oldCandidates := tgt.candidates
	tgt.candidates = names

	// Don't proceed with autotargeting if there's a current target which
	// isn't one of the previous candidates. I.e. a non-autotarget.
	if tgt.Name != "" && !slices.Contains(oldCandidates, tgt.Name) {
		return
	}

	tgt.retarget()
}

func (tgt *Target) SetPresent(names []string) {
	tgt.present = names
	tgt.retarget()
}

func (tgt *Target) AddPresent(name string) {
	tgt.present = append(tgt.present, name)
	tgt.retarget()
}

func (tgt *Target) RemovePresent(name string) {
	if i := slices.Index(tgt.present, name); i >= 0 {
		tgt.present = append(tgt.present[:i], tgt.present[i+1:]...)
	}
	tgt.retarget()
}

func (tgt *Target) retarget() {
	var shouldTarget string

outer:
	for _, candidate := range tgt.candidates {
		for _, present := range tgt.present {
			if strings.Contains(present, candidate) {
				shouldTarget = candidate
				break outer
			}
		}
	}

	if (shouldTarget != "" || len(tgt.candidates) == 0) && shouldTarget != tgt.Name {
		tgt.Set(shouldTarget)
	}
}
