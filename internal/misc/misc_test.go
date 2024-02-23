package misc

import (
	"modmanager/internal/mods"
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
	var selectedMod mods.Mod
	// selectedMod = nil
	x := isModNil(selectedMod)
	if !x {
		t.Error("Mod is nil but IsModNil returned false")
	}
	var selectedStory *mods.CustomStory
	// selectedStory = nil
	selectedMod = selectedStory
	x = isModNil(selectedMod)
	if !x {
		t.Error("Mod is type CS nil but IsModNil returned false")
	}
	var selectedConversion *mods.FullConversion
	// selectedConversion = nil
	selectedMod = selectedConversion
	x = isModNil(selectedMod)
	if !x {
		t.Error("Mod is type FC nil but IsModNil returned false")
	}
}

// TODO test deleteFolders
