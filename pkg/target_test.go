package pkg_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"

	"github.com/stretchr/testify/assert"
)

func TestTarget(t *testing.T) {
	tcs := map[string]struct {
		startName     string
		startHealth   int
		candidates    []string
		present       []string
		enters        []string
		leaves        []string
		newcandidates []string
		sets          []string
		name          string
		health        int
	}{
		"one present": {
			candidates: []string{"one", "two"},
			present:    []string{"a one thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"two present": {
			candidates: []string{"one", "two"},
			present:    []string{"a two thing"},
			sets:       []string{"two"},
			name:       "two",
		},

		"three present": {
			candidates: []string{"one", "two"},
			present:    []string{"a three thing"},
			name:       "",
		},

		"one one present": {
			candidates: []string{"one", "two"},
			present:    []string{"a one thing", "another one thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"two one present": {
			candidates: []string{"one", "two"},
			present:    []string{"a two thing", "a one thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"one enters": {
			candidates: []string{"one", "two"},
			enters:     []string{"a one thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"two enters": {
			candidates: []string{"one", "two"},
			enters:     []string{"a two thing"},
			sets:       []string{"two"},
			name:       "two",
		},

		"three enters": {
			candidates: []string{"one", "two"},
			enters:     []string{"a three thing"},
			name:       "",
		},

		"two present one enters": {
			candidates: []string{"one", "two"},
			present:    []string{"a two thing"},
			enters:     []string{"a one thing"},
			sets:       []string{"two", "one"},
			name:       "one",
		},

		"one present two enters": {
			candidates: []string{"one", "two"},
			present:    []string{"a one thing"},
			enters:     []string{"a two thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"one two present one leaves": {
			candidates: []string{"one", "two"},
			present:    []string{"a one thing", "a two thing"},
			leaves:     []string{"a one thing"},
			sets:       []string{"one", "two"},
			name:       "two",
		},

		"one two present two leaves": {
			candidates: []string{"one", "two"},
			present:    []string{"a one thing", "a two thing"},
			leaves:     []string{"a two thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"one present one leaves": {
			candidates: []string{"one", "two"},
			present:    []string{"a one thing"},
			leaves:     []string{"a one thing"},
			sets:       []string{"one"},
			name:       "one",
		},

		"two present two leaves": {
			candidates: []string{"one", "two"},
			present:    []string{"a two thing"},
			leaves:     []string{"a two thing"},
			sets:       []string{"two"},
			name:       "two",
		},

		"clear empty new candidates": {
			candidates:    []string{"one", "two"},
			present:       []string{"a one thing"},
			newcandidates: []string{},
			sets:          []string{"one", ""},
			name:          "",
		},

		"retarget new candidates": {
			candidates:    []string{"one", "two"},
			present:       []string{"a three thing"},
			newcandidates: []string{"three"},
			sets:          []string{"three"},
			name:          "three",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var sets []string
			tgt := pkg.NewTarget(func(name string, tgt *pkg.Target) {
				sets = append(sets, name)
				tgt.Name = name
			})

			if tc.candidates != nil {
				tgt.SetCandidates(tc.candidates)
			}

			if tc.present != nil {
				tgt.SetPresent(tc.present)
			}

			if tc.enters != nil {
				for _, name := range tc.enters {
					tgt.AddPresent(name)
				}
			}

			if tc.leaves != nil {
				for _, name := range tc.leaves {
					tgt.RemovePresent(name)
				}
			}

			if tc.newcandidates != nil {
				tgt.SetCandidates(tc.newcandidates)
			}

			assert.Equal(t, tc.sets, sets)
			assert.Equal(t, tc.name, tgt.Name)
			assert.Equal(t, tc.health, tgt.Health)
		})
	}
}
