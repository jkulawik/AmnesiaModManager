package mods

import (
	"fmt"
	"testing"
)

var TestWhiteNight = FullConversion{
	Name:            "White Night",
	MainInitConfig:  "testdata/wn_config/main_init.cfg",
	Logo:            "testdata/wn_graphics/graphics/main_menu/wn_menu_logo.png",
	LangFile:        "english.lang",
	UniqueResources: []string{"wn_models", "wn_sounds", "wn_graphics", "wn_models", "wn_music"},
}

func TestGetConversionFromInit(t *testing.T) {
	fc, err := GetConversionFromInit("testdata", "testdata/wn_config/main_init.cfg")
	if err != nil {
		t.Fatal(err)
	}

	if fc == nil {
		t.Fatal("GetConversionFromInit returned nil")
	}

	// Structs with arrays aren't comparable, easiest solution is to cast them to string
	fcString := fmt.Sprintf("%v", *fc)
	testString := fmt.Sprintf("%v", TestWhiteNight)

	if fcString != testString {
		t.Errorf("FC differs from expected. Original:\n%s\nGot:\n%s", testString, fcString)
	}
}

func TestGetMainInitFilepaths(t *testing.T) {
	mainInits, err := GetMainInitFilepaths("testdata")

	if err != nil {
		t.Fatal(err)
	}

	if len(mainInits) != 2 {
		t.Errorf("Number of main_inits differs from 2. Got: %s", mainInits)
	} else if mainInits[0] != "SomeMod/config/main_init.cfg" && mainInits[1] != "wn_config/main_init.cfg" {
		t.Errorf("Did not find one of the main inits. Got: %s", mainInits)
	}
}

func TestGetFullConversions(t *testing.T) {
	fcList, err := GetFullConversions("testdata")

	if err != nil {
		t.Fatal(err)
	}

	for _, fc := range fcList {
		t.Log(fc)
	}

	if len(fcList) != 2 {
		t.Errorf("FC list wasn't the expected length. Expected 2, got %d", len(fcList))
	}
}
