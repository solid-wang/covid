package validation

import (
	"fmt"
	"testing"
)

func TestIntSliceEqual(t *testing.T) {
	a := []int{1, 2}
	b := []int{1, 2}
	c := []int{2, 1}
	d := []int{2, 3}
	e := []int{}
	f := []int{}
	fmt.Printf("1 %t\n", IntSliceEqual(a, b))
	fmt.Printf("2 %t\n", IntSliceEqual(b, c))
	fmt.Printf("3 %t\n", IntSliceEqual(b, d))
	fmt.Printf("4 %t\n", IntSliceEqual(b, e))
	fmt.Printf("5 %t\n", IntSliceEqual(e, f))

}
