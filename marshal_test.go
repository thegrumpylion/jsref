package jsref

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalScalar(t *testing.T) {
	require := require.New(t)
	{
		out, err := Marshal(true)
		require.Nil(err)
		require.Equal(true, out.Bool())
	}
	{
		out, err := Marshal(42)
		require.Nil(err)
		require.Equal(42, out.Int())
	}
	{
		out, err := Marshal("whatever")
		require.Nil(err)
		require.Equal("whatever", out.String())
	}
	{
		in := new(bool)
		*in = true
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal(true, out.Bool())
	}
	{
		in := new(int)
		*in = 42
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal(42, out.Int())
	}
	{
		in := new(string)
		*in = "whatever"
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal("whatever", out.String())
	}
}

func TestMarshalArray(t *testing.T) {
	require := require.New(t)
	{
		in := []string{"one", "two", "three"}
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal(true, IsArray(out))
		for i := 0; i < out.Length(); i++ {
			require.Equal(in[i], out.Index(i).String())
		}
	}
	{
		inVal := []string{"one", "two", "three"}
		in := []*string{}
		for _, val := range inVal {
			v := new(string)
			*v = val
			in = append(in, v)
		}
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal(true, IsArray(out))
		for i := 0; i < out.Length(); i++ {
			require.Equal(*in[i], out.Index(i).String())
		}
	}
	{
		in := []int{3, 44, 555}
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal(true, IsArray(out))
		for i := 0; i < out.Length(); i++ {
			require.Equal(in[i], out.Index(i).Int())
		}
	}
	{
		inVal := []int{3, 44, 555}
		in := []*int{}
		for _, val := range inVal {
			v := new(int)
			*v = val
			in = append(in, v)
		}
		out, err := Marshal(in)
		require.Nil(err)
		require.Equal(true, IsArray(out))
		for i := 0; i < out.Length(); i++ {
			require.Equal(*in[i], out.Index(i).Int())
		}
	}
}

func TestMarshalMap(t *testing.T) {
	require := require.New(t)

	in := map[string]string{
		"ena":  "one",
		"dio":  "two",
		"tria": "three",
	}
	out, err := Marshal(in)

	require.Nil(err)
	require.Equal(true, IsObject(out))

	for _, k := range ObjectKeys(out) {
		require.Equal(in[k], out.Get(k).String())
	}
}

func TestMarshalStruct(t *testing.T) {
	require := require.New(t)

	lst := []string{"sdf", "eeee"}
	s := struct {
		Name string
		Age  int
		List []string
	}{
		Name: "tester",
		Age:  666,
		List: lst,
	}
	out, err := Marshal(s)

	require.Nil(err)
	require.Equal(true, IsObject(out))
	name := out.Get("Name")
	require.Equal("tester", name.String())
	age := out.Get("Age")
	require.Equal(666, age.Int())
	list := out.Get("List")
	require.Equal(true, IsArray(list))
	for i := 0; i < list.Length(); i++ {
		require.Equal(lst[i], list.Index(i).String())
	}
}
