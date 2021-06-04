package dice

import (
	"fmt"
	"testing"
)

func TestRoll2d20p3(t *testing.T) {
	res, ex, _ := Roll("2d20+3")
	fmt.Printf("%v = %v\n", ex, res)
	//todo assertions
}

func TestRolld20p3(t *testing.T) {
	res, ex, _ := Roll("d20+3")
	fmt.Printf("%v = %v\n", ex, res)
}

func TestRolld20(t *testing.T) {
	res, ex, _ := Roll("2d20-3")
	fmt.Printf("%v = %v\n", ex, res)
}

func TestRolld0(t *testing.T) {
	res, ex, _ := Roll("d0")
	fmt.Printf("%v = %v\n", ex, res)
}

func TestRollIncorrect(t *testing.T) {
	_, _, err := Roll("asdg")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		t.Errorf("result should have error")
	}
}
