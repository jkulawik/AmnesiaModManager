package main

import (
	"testing"
)

func TestCheckIsRootDir(t *testing.T) {
	err := CheckIsRootDir(".")

	if err != nil {
		t.Error(err)
		t.Error("Root directory is an Amnesia installation but one isn't detected")
	}
}

func TestFormatStringList(t *testing.T) {
	strList := []string{"A", "B"}

	formatted := formatStringList(strList)

	if formatted != "A\nB\n" {
		t.Errorf("String list format returned incorrect string:\n%s", formatted)
	}
}

func TestGetStringSpacer(t *testing.T) {
	s := getStringSpacer(2)

	if s != "  " {
		t.Errorf("GetStringSpacer returned incorrect string: %q (length %d)", s, len(s))
	}
}
