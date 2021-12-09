package jsref

import (
	"fmt"
	"testing"
	"time"

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
	testTime := time.Now()
	testTimeString := testTime.Format(time.RFC3339)
	s := struct {
		Name       string
		Age        int
		List       []string
		Time       time.Time
		unexported string
	}{
		Name:       "tester",
		Age:        666,
		List:       lst,
		Time:       testTime,
		unexported: "skip me",
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
	timeJS := out.Get("Time")
	require.Equal(testTimeString, timeJS.String())
	unexported := out.Get("unexported")
	require.True(unexported.IsUndefined())

}

func TestExample(t *testing.T) {
	type addr struct {
		Street string
		Num    int
	}
	type user struct {
		Name    string
		Email   []string
		IsAdmin bool
		Limit   int
		Addr    *addr
	}
	usr := &user{
		Name:    "admin",
		Email:   []string{"admin@domain", "root@domain"},
		IsAdmin: true,
		Limit:   9001,
		Addr: &addr{
			Street: "somewhere",
			Num:    99,
		},
	}

	usrjs, err := Marshal(usr)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("JS Values")
	fmt.Println("---------")
	fmt.Println("Name", usrjs.Get("Name").String())
	fmt.Println("IsAdmin", usrjs.Get("IsAdmin").Bool())
	fmt.Println("Limit", usrjs.Get("Limit").Int())
	for i := 0; i < usrjs.Get("Email").Length(); i++ {
		fmt.Println("Email", i, usrjs.Get("Email").Index(i))
	}
	fmt.Println("Addr.Street", usrjs.Get("Addr").Get("Street").String())
	fmt.Println("Addr.Num", usrjs.Get("Addr").Get("Num").Int())
	fmt.Println("---------")

	usrjs.Set("Name", "user")
	usrjs.Set("IsAdmin", false)
	usrjs.Set("Limit", 100)
	usrjs.Set("Email", []interface{}{"user@domain"})
	usrjs.Get("Addr").Set("Street", "somewhere else")
	usrjs.Get("Addr").Set("Num", 2)

	usr = &user{}

	err = Unmarshal(usr, usrjs)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Go Values")
	fmt.Println("---------")
	fmt.Println("Name", usr.Name)
	fmt.Println("IsAdmin", usr.IsAdmin)
	fmt.Println("Limit", usr.Limit)
	for i, e := range usr.Email {
		fmt.Println("Email", i, e)
	}
	fmt.Println("Addr.Street", usr.Addr.Street)
	fmt.Println("Addr.Num", usr.Addr.Num)
	fmt.Println("---------")
}
