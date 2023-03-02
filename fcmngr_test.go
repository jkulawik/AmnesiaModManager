package main

import (
	"fmt"
	"testing"
)

var TestWhiteNight = FullConversion{
	name:            "White Night",
	mainInitConfig:  "testdata/wn_config/main_init.cfg",
	logo:            "wn_menu_logo.png",
	uniqueResources: []string{"/wn_models", "/wn_sounds", "/wn_graphics", "/wn_models", "/wn_music"},
}

func TestReadConversionInit(t *testing.T) {
	fc, err := ReadConversionInit("testdata/SomeMod/config/main_init.cfg")
	t.Logf("FC: %s", fc)

	if err != nil {
		t.Error(err)
	}

	if fc.Variables.GameName != "A full conversion" {
		t.Errorf("Wrong FC name: %s", fc.Variables.GameName)
	}
}

func TestGetUniqueResources(t *testing.T) {
	res, err := GetUniqueResources("testdata/SomeMod/config/resources.cfg")

	if err != nil {
		t.Error(err)
	}

	if len(res) != 1 || res[0] != "/SomeMod" {
		t.Error("Parsed resource list differs from the expected one")
		t.Log(res)
	}

}

func TestGetLogoFromMenuConfig(t *testing.T) {
	logo, err := GetLogoFromMenuConfig("testdata/wn_config/menu.cfg")
	t.Logf("Menu config: %s", logo)

	if err != nil {
		t.Error(err)
	}

	if logo != "wn_menu_logo.png" {
		t.Errorf("Wrong FC logo: %s", logo)
	}
}

func TestGetConversionFromInit(t *testing.T) {
	path := "testdata/wn_config/main_init.cfg"
	fc, err := GetConversionFromInit(path, "testdata")
	if err != nil {
		t.Error(err)
	}

	// Structs with arrays aren't comparable, easiest solution is to cast them to string
	fcString := fmt.Sprintf("%v", *fc)
	testString := fmt.Sprintf("%v", TestWhiteNight)

	if fcString != testString {
		t.Errorf("FC differs from expected. Original:\n%s\nGot:\n%s", testString, fcString)
	}
}

func TestGetMainInitConfigs(t *testing.T) {
	mainInits, err := GetMainInitConfigs("testdata") // FIXME adding a trailing slash casuses runtime panic

	if err != nil {
		t.Error(err)
	}

	if len(mainInits) < 2 {
		t.Errorf("Did not find one of the main inits. Got: %s", mainInits)
	} else if mainInits[0] != "testdata/SomeMod/config/main_init.cfg" && mainInits[1] != "testdata/wn_config/main_init.cfg" {
		t.Errorf("Did not find one of the main inits. Got: %s", mainInits)
	}
}

func TestGetFullConversions(t *testing.T) {
	fcList, err := GetFullConversions("testdata")

	if err != nil {
		t.Error(err)
	}

	for _, fc := range fcList {
		t.Log(fc)
	}

	if len(fcList) != 2 {
		t.Errorf("FC list wasn't the expected length. Expected 2, got %d", len(fcList))
	}
}
