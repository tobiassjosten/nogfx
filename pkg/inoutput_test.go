package pkg_test

import (
	"testing"

	"github.com/tobiassjosten/nogfx/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInoutput(t *testing.T) {
	inout := pkg.NewInoutput(
		[][]byte{[]byte("qwer")},
		[][]byte{[]byte("asdf")},
	)

	{
		inoutreplace := inout.ReplaceInput(0, []byte("zxcv"))
		assert.Equal(t, pkg.Inoutput{
			Input: pkg.Exput{
				pkg.Line{Text: []byte("zxcv")},
			},
			Output: pkg.Exput{
				pkg.Line{Text: []byte("asdf")},
			},
		}, inoutreplace)
		require.NotEqual(t, inout, inoutreplace, "operation shouldn't mutate")
	}

	{
		inoutreplace := inout.ReplaceOutput(0, []byte("zxcv"))
		assert.Equal(t, pkg.Inoutput{
			Input: pkg.Exput{
				pkg.Line{Text: []byte("qwer")},
			},
			Output: pkg.Exput{
				pkg.Line{Text: []byte("zxcv")},
			},
		}, inoutreplace)
		require.NotEqual(t, inout, inoutreplace, "operation shouldn't mutate")
	}
}

func TestExput(t *testing.T) {
	exput := pkg.NewExput([]byte("asdf"))

	{
		exputadd := exput.Append([]byte("asdf"))
		assert.Equal(t, pkg.Exput{
			pkg.Line{Text: []byte("asdf")},
			pkg.Line{Text: []byte("asdf")},
		}, exputadd)
		require.NotEqual(t, exput, exputadd, "operation shouldn't mutate")
	}

	{
		exputaddbefore := exput.AddBefore(0, []byte("qwer"))
		assert.Equal(t, pkg.Exput{
			pkg.Line{
				Text:   []byte("asdf"),
				Before: []pkg.Text{pkg.Text([]byte("qwer"))},
			},
		}, exputaddbefore)
		require.NotEqual(t, exput, exputaddbefore, "operation shouldn't mutate")
	}

	{
		exputaddafter := exput.AddAfter(0, []byte("zxcv"))
		assert.Equal(t, pkg.Exput{
			pkg.Line{
				Text:  []byte("asdf"),
				After: []pkg.Text{pkg.Text([]byte("zxcv"))},
			},
		}, exputaddafter)
		require.NotEqual(t, exput, exputaddafter, "operation shouldn't mutate")
	}

	{
		var bytes [][]byte
		assert.Equal(t, bytes, exput.Remove(0).Bytes())
	}

	{
		exputreplace := exput.Replace(0, []byte("fdsa"))
		assert.Equal(t, pkg.Exput{
			pkg.Line{Text: []byte("fdsa")},
		}, exputreplace)
		require.NotEqual(t, exput, exputreplace, "operation shouldn't mutate")
	}

	{
		exputsplit := exput.Split([]byte{'d'})
		assert.Equal(t, pkg.Exput{
			pkg.Line{Text: []byte("as")},
			pkg.Line{Text: []byte("f")},
		}, exputsplit)
		require.NotEqual(t, exput, exputsplit, "operation shouldn't mutate")
	}

	{
		// Maintains order.
		exputsplit := exput.Append([]byte("qwer")).Split([]byte{'d'})
		assert.Equal(t, pkg.Exput{
			pkg.Line{Text: []byte("as")},
			pkg.Line{Text: []byte("f")},
			pkg.Line{Text: []byte("qwer")},
		}, exputsplit)
		require.NotEqual(t, exput, exputsplit, "operation shouldn't mutate")
	}

	{
		exputsplit := exput.Split([]byte{'x'})
		assert.Equal(t, pkg.Exput{
			pkg.Line{Text: []byte("asdf")},
		}, exputsplit)
		require.Equal(t, exput, exputsplit, "end result is the same")
	}

	assert.Equal(t, pkg.Inoutput{Input: exput}, exput.Inoutput(pkg.Input))
	assert.Equal(t, pkg.Inoutput{Output: exput}, exput.Inoutput(pkg.Output))
	assert.Equal(t, pkg.Inoutput{}, exput.Inoutput(pkg.IOKind("asdf")))

	assert.Equal(t, [][]byte{
		[]byte("qw"),
		[]byte("as"),
		[]byte("zx"),
		[]byte("er"),
		[]byte("df"),
		[]byte("cv"),
	}, (pkg.Exput{
		pkg.Line{
			Text:   []byte("as"),
			Before: []pkg.Text{pkg.Text([]byte("qw"))},
			After:  []pkg.Text{pkg.Text([]byte("zx"))},
		},
		pkg.Line{
			Text: []byte("secret"),
		},
		pkg.Line{
			Text:   []byte("df"),
			Before: []pkg.Text{pkg.Text([]byte("er"))},
			After:  []pkg.Text{pkg.Text([]byte("cv"))},
		},
	}).Remove(1).Bytes())
}

func TestClean(t *testing.T) {
	tcs := map[string]struct {
		in  []byte
		out []byte
	}{
		"colored": {
			in:  []byte("\033[35masdf"),
			out: []byte("asdf"),
		},

		"plain": {
			in:  []byte("asdf"),
			out: []byte("asdf"),
		},

		"invalid start": {
			in:  []byte("\033(35masdf"),
			out: []byte("(35masdf"),
		},

		"unterminated": {
			in:  []byte("\033[35asdf"),
			out: []byte(""),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			text := pkg.Text(tc.in)
			assert.Equal(t, tc.out, text.Clean())
		})
	}
}
