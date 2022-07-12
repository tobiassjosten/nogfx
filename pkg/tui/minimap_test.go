package tui

import (
	"log"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/icza/gox/gox"
	"github.com/stretchr/testify/assert"
	"github.com/tobiassjosten/nogfx/pkg/mock"
	"github.com/tobiassjosten/nogfx/pkg/navigation"
)

func TestRenderTUIMap(t *testing.T) {
	ui := NewTUI(&mock.ScreenMock{
		HideCursorFunc:     func() {},
		SetCursorStyleFunc: func(_ tcell.CursorStyle) {},
		SetStyleFunc:       func(_ tcell.Style) {},
	})

	ui.SetRoom(&navigation.Room{})
	rows := ui.RenderMap(5, 3)

	assert.Equal(t,
		[]string{"     ", " [ ] ", "     "},
		rows.Strings(),
	)

	ui.room = &navigation.Room{Exits: map[string]*navigation.Room{
		"n": {},
	}}
	rows = ui.RenderMap(5, 3)

	// Cache exists and rendition is the same as before.
	assert.Equal(t,
		[]string{"     ", " [ ] ", "     "},
		rows.Strings(),
	)

	ui.setCache(paneMap, nil)
	rows = ui.RenderMap(5, 3)

	// Cache is cleared and a new map is rendered.
	assert.Equal(t,
		[]string{"  |  ", " [ ] ", "     "},
		rows.Strings(),
	)
}

func TestRenderMap(t *testing.T) {
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

		"single room too little space": {
			room:   &navigation.Room{},
			width:  4,
			height: 2,
			visual: []string{"    ", "    "},
		},

		"single room no space": {
			room:   &navigation.Room{},
			width:  0,
			height: 0,
			visual: []string{},
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
			rows := RenderMap(tc.room, tc.width, tc.height)
			if !assert.Equal(t,
				strings.Join(tc.visual, "\n"),
				strings.Join(rows.Strings(), "\n"),
			) {
				log.Println("RENDITION")
				for _, str := range rows.Strings() {
					log.Println(str)
				}
			}
		})
	}
}
