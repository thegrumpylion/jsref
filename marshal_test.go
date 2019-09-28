package jsref

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	s := struct {
		Name string
		Age  int
	}{
		Name: "tester",
		Age:  666,
	}

	out, err := Marshal(s)
	fmt.Println(out, err)
}
