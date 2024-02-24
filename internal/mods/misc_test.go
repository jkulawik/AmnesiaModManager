package mods

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

func TestIsModNil(t *testing.T) {
	var selectedMod Mod
	// selectedMod = nil
	x := IsModNil(selectedMod)
	if !x {
		t.Error("Mod is nil but IsModNil returned false")
	}
	var selectedStory *CustomStory
	// selectedStory = nil
	selectedMod = selectedStory
	x = IsModNil(selectedMod)
	if !x {
		t.Error("Mod is type CS nil but IsModNil returned false")
	}
	var selectedConversion *FullConversion
	// selectedConversion = nil
	selectedMod = selectedConversion
	x = IsModNil(selectedMod)
	if !x {
		t.Error("Mod is type FC nil but IsModNil returned false")
	}
}

// TODO test deleteFolders
