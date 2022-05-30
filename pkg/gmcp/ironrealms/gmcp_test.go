package ironrealms_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tobiassjosten/nogfx/pkg/gmcp"
	"github.com/tobiassjosten/nogfx/pkg/gmcp/ironrealms"

	"github.com/stretchr/testify/assert"
)

func makeGMCP(id string, data interface{}) string {
	jsondata, _ := json.Marshal(data)
	return fmt.Sprintf("%s %s", id, string(jsondata))
}

func TestParse(t *testing.T) {
	tcs := map[string]struct {
		data string
		msg  gmcp.Message
		err  string
	}{
		"Char.Status": {
			data: "Char.Status {}",
			msg: &ironrealms.CharStatus{
				CharStatus: &gmcp.CharStatus{},
			},
		},

		"Char.Vitals": {
			data: "Char.Vitals {}",
			msg:  &ironrealms.CharVitals{},
		},

		"non-existant": {
			data: "Non.Existant",
			err:  "unknown message 'Non.Existant'",
		},

		"invalid JSON": {
			data: "Char.Status asdf",
			err:  "couldn't unmarshal *ironrealms.CharStatus: invalid character 'a' looking for beginning of value",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			msg, err := ironrealms.Parse([]byte(tc.data))

			if tc.err != "" {
				if assert.NotNil(t, err) {
					assert.Equal(t, tc.err, err.Error())
				}
				return
			} else if err != nil {
				assert.Equal(t, "", err.Error())
				return
			}

			assert.Equal(t, tc.msg, msg)
		})
	}
}
