package tui_test

import (
	"log"
	"strings"
	"testing"

	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
	"github.com/tobiassjosten/nogfx/pkg/tui"
)

func TestRenderMap(t *testing.T) {
	rowToString := func(row tui.Row) (str string) {
		for _, c := range row {
			str += string(c.Content)
		}
		return
	}

	rowsToStrings := func(rows tui.Rows) (strs []string) {
		for _, row := range rows {
			strs = append(strs, rowToString(row))
		}
		return
	}

	stringsToMap := func(strs []string) string {
		return strings.Join(strs, "\n")
	}

	tcs := map[string]struct {
		room   *navigation.Room
		width  int
		height int
		visual []string
	}{
		"single room no content": {
			room:   &navigation.Room{},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [ ] `,
				`     `,
			},
		},

		"single room player": {
			room: &navigation.Room{
				HasPlayer: true,
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [+] `,
				`     `,
			},
		},

		"unknown exit": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"asdf": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [ ] `,
				`     `,
			},
		},

		"up exit no player": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"u": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [^] `,
				`     `,
			},
		},

		"up exit with player": {
			room: &navigation.Room{
				HasPlayer: true,
				Exits: map[string]*navigation.Room{
					"u": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [+] `,
				`     `,
			},
		},

		"down exit no player": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"d": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [v] `,
				`     `,
			},
		},

		"down exit with player": {
			room: &navigation.Room{
				HasPlayer: true,
				Exits: map[string]*navigation.Room{
					"d": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [+] `,
				`     `,
			},
		},

		"up down exit no player": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"u": &navigation.Room{},
					"d": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [=] `,
				`     `,
			},
		},

		"up down exit with player": {
			room: &navigation.Room{
				HasPlayer: true,
				Exits: map[string]*navigation.Room{
					"u": &navigation.Room{},
					"d": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [+] `,
				`     `,
			},
		},

		"in exit": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"in": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` [ } `,
				`     `,
			},
		},

		"out exit": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"out": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`     `,
				` { ] `,
				`     `,
			},
		},

		"all exits": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"n":   &navigation.Room{},
					"ne":  &navigation.Room{},
					"e":   &navigation.Room{},
					"se":  &navigation.Room{},
					"s":   &navigation.Room{},
					"sw":  &navigation.Room{},
					"w":   &navigation.Room{},
					"nw":  &navigation.Room{},
					"u":   &navigation.Room{},
					"d":   &navigation.Room{},
					"in":  &navigation.Room{},
					"out": &navigation.Room{},
				},
			},
			width:  5,
			height: 3,
			visual: []string{
				`\ | /`,
				`-{=}-`,
				`/ | \`,
			},
		},

		"all adjacent rooms": {
			room: &navigation.Room{
				Exits: map[string]*navigation.Room{
					"n":   &navigation.Room{},
					"ne":  &navigation.Room{},
					"e":   &navigation.Room{},
					"se":  &navigation.Room{},
					"s":   &navigation.Room{},
					"sw":  &navigation.Room{},
					"w":   &navigation.Room{},
					"nw":  &navigation.Room{},
					"u":   &navigation.Room{},
					"d":   &navigation.Room{},
					"in":  &navigation.Room{},
					"out": &navigation.Room{},
				},
			},
			width:  13,
			height: 7,
			visual: []string{
				`             `,
				` [ ] [ ] [ ] `,
				`    \ | /    `,
				` [ ]-{=}-[ ] `,
				`    / | \    `,
				` [ ] [ ] [ ] `,
				`             `,
			},
		},

		"criss-cross": {
			room: &navigation.Room{
				HasPlayer: true,
				Exits: map[string]*navigation.Room{
					"n": &navigation.Room{
						Exits: map[string]*navigation.Room{
							"sw": &navigation.Room{},
						},
					},
					"w": &navigation.Room{
						Exits: map[string]*navigation.Room{
							"ne": &navigation.Room{},
						},
					},
					"nw": &navigation.Room{},
				},
			},
			width:  13,
			height: 7,
			visual: []string{
				`             `,
				` [ ] [ ]     `,
				`    X |      `,
				` [ ]-[+]     `,
				`             `,
				`             `,
				`             `,
			},
		},

		"long distance": {
			room: &navigation.Room{
				ID:        1,
				HasPlayer: true,
				X:         gox.NewInt(3),
				Y:         gox.NewInt(3),
				Exits: map[string]*navigation.Room{
					"n": &navigation.Room{
						ID: 2,
						X:  gox.NewInt(3),
						Y:  gox.NewInt(5),
						Exits: map[string]*navigation.Room{
							"sw": &navigation.Room{
								ID: 3,
								X:  gox.NewInt(1),
								Y:  gox.NewInt(3),
							},
						},
					},
					"w": &navigation.Room{
						ID: 3,
						X:  gox.NewInt(1),
						Y:  gox.NewInt(3),
						Exits: map[string]*navigation.Room{
							"ne": &navigation.Room{
								ID: 2,
								X:  gox.NewInt(3),
								Y:  gox.NewInt(5),
							},
						},
					},
				},
			},
			width:  21,
			height: 11,
			visual: []string{
				`                     `,
				`         [ ]         `,
				`        / |          `,
				`      /   |          `,
				`    /     |          `,
				` [ ]- - -[+]         `,
				`                     `,
				`                     `,
				`                     `,
				`                     `,
				`                     `,
			},
		},

		"exits outwards only": {
			room: &navigation.Room{
				ID: 1,
				Exits: map[string]*navigation.Room{
					"in": &navigation.Room{
						ID: 2,
						Exits: map[string]*navigation.Room{
							"w": &navigation.Room{ID: 1},
						},
					},
					"out": &navigation.Room{
						ID: 3,
						Exits: map[string]*navigation.Room{
							"e": &navigation.Room{ID: 1},
						},
					},
				},
			},
			width:  13,
			height: 3,
			visual: []string{
				`             `,
				` [ ] { } [ ] `,
				`             `,
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			rows := tui.RenderMap(tc.room, tc.width, tc.height)
			if !assert.Equal(t,
				stringsToMap(tc.visual),
				stringsToMap(rowsToStrings(rows)),
			) {
				log.Println("RENDITION")
				for _, str := range rowsToStrings(rows) {
					log.Println(str)
				}
			}
		})
	}
}
