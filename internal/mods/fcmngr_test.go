package mods

import (
	"fmt"
	"testing"
)

var TestWhiteNight = FullConversion{
	name:            "White Night",
	mainInitConfig:  "wn_config/main_init.cfg",
	logo:            "wn_graphics/graphics/main_menu/wn_menu_logo.png",
	uniqueResources: []string{"wn_models", "wn_sounds", "wn_graphics", "wn_models", "wn_music"},
}

var someModRes = []string{"SomeMod", "SomeMod/misc"}

func TestReadConversionInit(t *testing.T) {
	fc, err := ReadConversionInit("SomeMod/config/main_init.cfg")
	t.Logf("FC: %s", fc)

	if err != nil {
		t.Error(err)
	}

	if fc.Variables.GameName != "A full conversion" {
		t.Errorf("Wrong FC name: %s", fc.Variables.GameName)
	}
}

func TestGetUniqueResources(t *testing.T) {
	res, err := GetUniqueResources("SomeMod/config/resources.cfg")

	if err != nil {
		t.Error(err)
	}

	resString := fmt.Sprintf("%v", res)
	testString := fmt.Sprintf("%v", someModRes)

	if resString != testString {
		t.Errorf("Parsed resource list differs expected. Original:\n%s\nGot:\n%s", testString, resString)
	}

}

func TestGetLogoFromMenuConfig(t *testing.T) {
	logo, err := GetLogoFromMenuConfig("wn_config/menu.cfg", TestWhiteNight.uniqueResources)
	t.Logf("Logo path: %s", logo)

	if err != nil {
		t.Error(err)
	}

	if logo != "wn_graphics/graphics/main_menu/wn_menu_logo.png" {
		t.Errorf("Wrong FC logo: %s", logo)
	}
}

func TestGetLogoFromMenuConfig2(t *testing.T) {
	logo, err := GetLogoFromMenuConfig("SomeMod/config/menu.cfg", []string{"/SomeMod"})
	t.Logf("Logo path: %s", logo)

	if err != nil {
		t.Error(err)
	}

	if logo != "SomeMod/menu_logo.tga" {
		t.Errorf("Wrong FC logo: %s", logo)
	}
}

func TestGetConversionFromInit(t *testing.T) {
	path := "wn_config/main_init.cfg"
	fc, err := GetConversionFromInit(path)
	if err != nil {
		t.Errorf("GetConversionFromInit returned an error: %s", err)
	}

	if fc == nil {
		t.Fatal("FATAL: GetConversionFromInit returned nil")
	}

	// Structs with arrays aren't comparable, easiest solution is to cast them to string
	fcString := fmt.Sprintf("%v", *fc)
	testString := fmt.Sprintf("%v", TestWhiteNight)

	if fcString != testString {
		t.Errorf("FC differs from expected. Original:\n%s\nGot:\n%s", testString, fcString)
	}
}

func TestGetMainInitConfigs(t *testing.T) {
	mainInits, err := GetMainInitConfigs()

	if err != nil {
		t.Error(err)
	}

	if len(mainInits) < 2 {
		t.Errorf("Did not find one of the main inits. Got: %s", mainInits)
	} else if mainInits[0] != "SomeMod/config/main_init.cfg" && mainInits[1] != "wn_config/main_init.cfg" {
		t.Errorf("Did not find one of the main inits. Got: %s", mainInits)
	}
}

func TestGetFullConversions(t *testing.T) {
	fcList, err := GetFullConversions()

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
