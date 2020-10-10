package jsref

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnamrshalArray(t *testing.T) {
	val := js.ValueOf([]interface{}{"one", "two", "three"})
	out := []string{}
	err := Unmarshal(&out, val)

	require := require.New(t)

	require.Nil(err)
	require.Equal("one", out[0])
	require.Equal("two", out[1])
	require.Equal("three", out[2])
}

func TestUnamrshalMap(t *testing.T) {
	obj := js.ValueOf(map[string]interface{}{
		"ena":  "one",
		"dio":  "two",
		"tria": "three",
	})
	keys := ObjectKeys(obj)
	m := map[string]string{}

	err := Unmarshal(&m, obj)

	require := require.New(t)

	require.Nil(err)
	require.ElementsMatch([]string{"ena", "dio", "tria"}, keys)
	require.Equal("one", m["ena"])
	require.Equal("two", m["dio"])
	require.Equal("three", m["tria"])
}

func TestUnamrshalStruct(t *testing.T) {
	obj := js.ValueOf(map[string]interface{}{
		"Ena":  "one",
		"dio":  true,
		"Tria": 3,
	})
	out := &struct {
		Ena  string
		Dio  bool `jsref:"dio"`
		Tria uint
	}{}
	err := Unmarshal(out, obj)

	require := require.New(t)

	require.Nil(err)
	require.Equal("one", out.Ena)
	require.Equal(true, out.Dio)
	require.Equal(uint(3), out.Tria)

}
