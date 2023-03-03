package main

import "testing"

func TestCheckIsRootDir(t *testing.T) {
	err := CheckIsRootDir(".")

	if err != nil {
		t.Error(err)
		t.Error("Root directory is an Amnesia installation but one isn't detected")
	}
}
