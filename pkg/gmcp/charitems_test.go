package gmcp_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tobiassjosten/nogfx/pkg/gmcp"
)

func TestCharItemsMessages(t *testing.T) {
	tcs := map[string]struct {
		msg         gmcp.Message
		data        string
		unmarshaled gmcp.Message
		marshaled   string
		err         string
	}{
		"Char.Items.Contents one": {
			msg:  &gmcp.CharItemsContents{},
			data: "Char.Items.Contents 1",
			unmarshaled: &gmcp.CharItemsContents{
				Container: 1,
			},
			marshaled: "Char.Items.Contents 1",
		},

		"Char.Items.Contents invalid JSON": {
			msg:  &gmcp.CharItemsContents{},
			data: "Char.Items.Contents asdf",
			err:  "invalid character 'a' looking for beginning of value",
		},

		"Char.Items.Inv": {
			msg:         &gmcp.CharItemsInv{},
			data:        "Char.Items.Inv",
			unmarshaled: &gmcp.CharItemsInv{},
			marshaled:   "Char.Items.Inv",
		},

		"Char.Items.Room": {
			msg:         &gmcp.CharItemsRoom{},
			data:        "Char.Items.Room",
			unmarshaled: &gmcp.CharItemsRoom{},
			marshaled:   "Char.Items.Room",
		},

		"Char.Items.List empty": {
			msg:         &gmcp.CharItemsList{},
			data:        "Char.Items.List {}",
			unmarshaled: &gmcp.CharItemsList{},
			marshaled: makeGMCP("Char.Items.List", map[string]interface{}{
				"location": "",
				"items":    []string{},
			}),
		},

		"Char.Items.List hydrated": {
			msg: &gmcp.CharItemsList{},
			data: makeGMCP("Char.Items.List", map[string]interface{}{
				"location": "room",
				"items": []map[string]interface{}{
					{
						"id":     1234,
						"name":   "an item",
						"attrib": "wWlLgcrfemdtx",
						"icon":   "item",
					},
				},
			}),
			unmarshaled: &gmcp.CharItemsList{
				Location: "room",
				Items: []gmcp.CharItem{
					{
						ID:   1234,
						Name: "an item",
						Attributes: gmcp.CharItemAttributes{
							Container:    true,
							Dangerous:    true,
							Dead:         true,
							Edible:       true,
							Fluid:        true,
							Groupable:    true,
							Monster:      true,
							Riftable:     true,
							Takeable:     true,
							Wearable:     true,
							Worn:         true,
							WieldedLeft:  true,
							WieldedRight: true,
						},
						Icon: "item",
					},
				},
			},
			marshaled: makeGMCP("Char.Items.List", map[string]interface{}{
				"location": "room",
				"items": []map[string]interface{}{
					{
						"id":     1234,
						"name":   "an item",
						"attrib": "cdefglLmrtwWx",
						"icon":   "item",
					},
				},
			}),
		},

		"Char.Items.List invalid attribute": {
			msg: &gmcp.CharItemsList{},
			data: makeGMCP("Char.Items.List", map[string]interface{}{
				"location": "room",
				"items": []map[string]interface{}{
					{
						"id":     1234,
						"name":   "an item",
						"attrib": "z",
						"icon":   "item",
					},
				},
			}),
			err: "unknown attribute 'z'",
		},

		"Char.Items.Add empty": {
			msg:         &gmcp.CharItemsAdd{},
			data:        "Char.Items.Add {}",
			unmarshaled: &gmcp.CharItemsAdd{},
			marshaled: makeGMCP("Char.Items.Add", map[string]interface{}{
				"location": "",
				"item": map[string]interface{}{
					"id":     0,
					"name":   "",
					"attrib": "",
					"icon":   "",
				},
			}),
		},

		"Char.Items.Add hydrated": {
			msg: &gmcp.CharItemsAdd{},
			data: makeGMCP("Char.Items.Add", map[string]interface{}{
				"location": "room",
				"item": map[string]interface{}{
					"id":     1234,
					"name":   "an item",
					"attrib": "wWlLgcrfemdtx",
					"icon":   "item",
				},
			}),
			unmarshaled: &gmcp.CharItemsAdd{
				Location: "room",
				Item: gmcp.CharItem{
					ID:   1234,
					Name: "an item",
					Attributes: gmcp.CharItemAttributes{
						Container:    true,
						Dangerous:    true,
						Dead:         true,
						Edible:       true,
						Fluid:        true,
						Groupable:    true,
						Monster:      true,
						Riftable:     true,
						Takeable:     true,
						Wearable:     true,
						Worn:         true,
						WieldedLeft:  true,
						WieldedRight: true,
					},
					Icon: "item",
				},
			},
			marshaled: makeGMCP("Char.Items.Add", map[string]interface{}{
				"location": "room",
				"item": map[string]interface{}{
					"id":     1234,
					"name":   "an item",
					"attrib": "cdefglLmrtwWx",
					"icon":   "item",
				},
			}),
		},

		"Char.Items.Remove empty": {
			msg:         &gmcp.CharItemsRemove{},
			data:        "Char.Items.Remove {}",
			unmarshaled: &gmcp.CharItemsRemove{},
			marshaled: makeGMCP("Char.Items.Remove", map[string]interface{}{
				"location": "",
				"item": map[string]interface{}{
					"id":     0,
					"name":   "",
					"attrib": "",
					"icon":   "",
				},
			}),
		},

		"Char.Items.Remove hydrated": {
			msg: &gmcp.CharItemsRemove{},
			data: makeGMCP("Char.Items.Remove", map[string]interface{}{
				"location": "room",
				"item": map[string]interface{}{
					"id":     1234,
					"name":   "an item",
					"attrib": "wWlLgcrfemdtx",
					"icon":   "item",
				},
			}),
			unmarshaled: &gmcp.CharItemsRemove{
				Location: "room",
				Item: gmcp.CharItem{
					ID:   1234,
					Name: "an item",
					Attributes: gmcp.CharItemAttributes{
						Container:    true,
						Dangerous:    true,
						Dead:         true,
						Edible:       true,
						Fluid:        true,
						Groupable:    true,
						Monster:      true,
						Riftable:     true,
						Takeable:     true,
						Wearable:     true,
						Worn:         true,
						WieldedLeft:  true,
						WieldedRight: true,
					},
					Icon: "item",
				},
			},
			marshaled: makeGMCP("Char.Items.Remove", map[string]interface{}{
				"location": "room",
				"item": map[string]interface{}{
					"id":     1234,
					"name":   "an item",
					"attrib": "cdefglLmrtwWx",
					"icon":   "item",
				},
			}),
		},

		"Char.Items.Update empty": {
			msg:         &gmcp.CharItemsUpdate{},
			data:        "Char.Items.Update {}",
			unmarshaled: &gmcp.CharItemsUpdate{},
			marshaled: makeGMCP("Char.Items.Update", map[string]interface{}{
				"location": "",
				"item": map[string]interface{}{
					"id":     0,
					"name":   "",
					"attrib": "",
					"icon":   "",
				},
			}),
		},

		"Char.Items.Update hydrated": {
			msg: &gmcp.CharItemsUpdate{},
			data: makeGMCP("Char.Items.Update", map[string]interface{}{
				"location": "room",
				"item": map[string]interface{}{
					"id":     1234,
					"name":   "an item",
					"attrib": "wWlLgcrfemdtx",
					"icon":   "item",
				},
			}),
			unmarshaled: &gmcp.CharItemsUpdate{
				Location: "room",
				Item: gmcp.CharItem{
					ID:   1234,
					Name: "an item",
					Attributes: gmcp.CharItemAttributes{
						Container:    true,
						Dangerous:    true,
						Dead:         true,
						Edible:       true,
						Fluid:        true,
						Groupable:    true,
						Monster:      true,
						Riftable:     true,
						Takeable:     true,
						Wearable:     true,
						Worn:         true,
						WieldedLeft:  true,
						WieldedRight: true,
					},
					Icon: "item",
				},
			},
			marshaled: makeGMCP("Char.Items.Update", map[string]interface{}{
				"location": "room",
				"item": map[string]interface{}{
					"id":     1234,
					"name":   "an item",
					"attrib": "cdefglLmrtwWx",
					"icon":   "item",
				},
			}),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.Unmarshal([]byte(tc.data))

			if tc.err != "" {
				require.NotNil(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			} else if err != nil {
				require.Equal(t, "", err.Error())
			}

			require.Equal(t, tc.unmarshaled, tc.msg, "unmarshaling hydrates message")

			if tc.marshaled == "" {
				return
			}

			marshaled := tc.msg.Marshal()
			data := strings.TrimSpace(strings.TrimPrefix(marshaled, tc.msg.ID()))
			tcdata := strings.TrimSpace(strings.TrimPrefix(tc.marshaled, tc.msg.ID()))

			assert.NotEqual(t, marshaled, data, "marshaled data has ID prefix")
			assert.NotEqual(t, tc.marshaled, tcdata, "marshaled data has ID prefix")

			if tcdata == "" {
				assert.Equal(t, tcdata, data)
				return
			}

			assert.JSONEq(t, tcdata, data, "marshaling maintains data integrity")

			require.Equal(t, tc.unmarshaled, tc.msg, "marshaling doesn't mutate")
		})
	}
}
