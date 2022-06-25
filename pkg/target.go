package pkg

import (
	"strings"

	"golang.org/x/exp/slices"
)

type TargetSetter func(string, *Target)

type Target struct {
	Name   string
	Health int

	setter TargetSetter

	// A list of possible targets, e.g. "manticore".
	candidates []string

	// A list of NPCs present, e.g. "a ferocious manticore".
	present []string
}

func NewTarget(setter TargetSetter) *Target {
	return &Target{
		Health: -1,
		setter: setter,
	}
}

func (tgt *Target) Set(name string) {
	if name == tgt.Name {
		return
	}

	tgt.setter(name, tgt)
}

func (tgt *Target) SetCandidates(names []string) {
	oldCandidates := tgt.candidates
	tgt.candidates = names

	isOldCandidate := slices.Contains(oldCandidates, tgt.Name)
	isNewCandidate := slices.Contains(tgt.candidates, tgt.Name)

	if isOldCandidate && !isNewCandidate {
		tgt.Set("")
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
		tgt.retarget()
	}
}

func (tgt *Target) Queue() int {
	queue := 0
	for _, present := range tgt.present {
		if tgt.Name != "" && strings.Contains(present, tgt.Name) {
			queue++
			continue
		}

		for _, candidate := range tgt.candidates {
			if strings.Contains(present, candidate) {
				queue++
				break
			}
		}
	}

	return queue
}

func (tgt *Target) retarget() {
	if tgt.Name != "" && !slices.Contains(tgt.candidates, tgt.Name) {
		return
	}

	var newTarget string

outer:
	for _, candidate := range tgt.candidates {
		for _, present := range tgt.present {
			if strings.Contains(present, candidate) {
				newTarget = candidate
				break outer
			}
		}
	}

	if newTarget != "" && newTarget != tgt.Name {
		tgt.Set(newTarget)
	}
}
