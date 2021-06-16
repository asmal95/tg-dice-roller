package dice

import (
	"testing"
)

func TestRoll2d20p3(t *testing.T) {
	res, ex, err := Roll("2d20+3")
	if res == 0 {
		t.Errorf("The result shouldn't be 0")
	}
	if ex == "" {
		t.Errorf("The expiration shouldn't be empty")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v = %v\n", ex, res)
}

func TestRolld20p3(t *testing.T) {
	res, ex, err := Roll("d20+3")
	if res == 0 {
		t.Errorf("The result shouldn't be 0")
	}
	if ex == "" {
		t.Errorf("The expiration shouldn't be empty")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v = %v\n", ex, res)
}

func TestRoll2d20m3(t *testing.T) {
	res, ex, err := Roll("2d20-3")
	if res == 0 {
		t.Errorf("The result shouldn't be 0")
	}
	if ex == "" {
		t.Errorf("The expiration shouldn't be empty")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("%v = %v\n", ex, res)
}

func TestRolld20(t *testing.T) {
	res, ex, err := Roll("d20")
	if res == 0 {
		t.Errorf("The result shouldn't be 0")
	}
	if ex != "" {
		t.Errorf("The expiration should be empty")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("res = %v\n", res)
}

func TestRolld0(t *testing.T) {
	_, _, err := Roll("d0")
	if err == nil {
		t.Errorf("the error shouldn't be nil")
	}
	t.Log(err)
}

func TestRolld10001(t *testing.T) {
	_, _, err := Roll("d10001")
	if err == nil {
		t.Errorf("the error shouldn't be nill")
	}
	t.Log(err)
}

func TestRoll0d20(t *testing.T) {
	_, _, err := Roll("0d20")
	if err == nil {
		t.Errorf("the error shouldn't be nill")
	}
	t.Log(err)
}

func TestRoll31d20(t *testing.T) {
	_, _, err := Roll("31d20")
	if err == nil {
		t.Errorf("the error shouldn't be nill")
	}
	t.Log(err)
}

func TestRollIncorrect(t *testing.T) {
	_, _, err := Roll("asdg")
	if err != nil {
		t.Log(err)
	} else {
		t.Errorf("result should have error")
	}
}
