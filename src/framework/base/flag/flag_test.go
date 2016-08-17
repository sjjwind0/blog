package flag

import (
	"testing"
)

func Test_Parse0(t *testing.T) {
	command := "test"
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse1(t *testing.T) {
	command := "test a"
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse2(t *testing.T) {
	command := "test     a"
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse3(t *testing.T) {
	command := "   test     a   "
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse4(t *testing.T) {
	command := "   test     a "
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse5(t *testing.T) {
	command := "   test     a   b  "
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse6(t *testing.T) {
	command := "   test "
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}

func Test_Parse7(t *testing.T) {
	command := "   "
	c, err := Parse(command)
	if err != nil {
		t.Error(err.Error())
		return
	}
	c.show()
}
