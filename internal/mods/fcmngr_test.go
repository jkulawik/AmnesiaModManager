package mods

import (
	"fmt"
	"testing"
)

var TestWhiteNight = FullConversion{
	Name:            "White Night",
	MainInitConfig:  "wn_config/main_init.cfg",
	Logo:            "wn_graphics/graphics/main_menu/wn_menu_logo.png",
	UniqueResources: []string{"wn_models", "wn_sounds", "wn_graphics", "wn_models", "wn_music"},
}

var someModRes = []string{"SomeMod", "SomeMod/misc"}

func TestReadConversionInit(t *testing.T) {
	fc, err := ReadConversionInit("testdata/SomeMod/config/main_init.cfg")
	t.Logf("FC: %s", fc)

	if err != nil {
		t.Fatal(err)
	}

	if fc.Variables.GameName != "A full conversion" {
		t.Errorf("Wrong FC name: %s", fc.Variables.GameName)
	}
}

func TestGetUniqueResources(t *testing.T) {
	res, err := GetUniqueResources("testdata/SomeMod/config/resources.cfg")

	if err != nil {
		t.Fatal(err)
	}

	resString := fmt.Sprintf("%v", res)
	testString := fmt.Sprintf("%v", someModRes)

	if resString != testString {
		t.Errorf("Parsed resource list differs expected. Original:\n%s\nGot:\n%s", testString, resString)
	}

}

func TestGetLogoFromMenuConfigPng(t *testing.T) {
	logo, err := GetLogoFromMenuConfig("testdata/wn_config/menu.cfg", TestWhiteNight.UniqueResources)

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logo path: %s", logo)

	if logo != "testdata/wn_graphics/graphics/main_menu/wn_menu_logo.png" {
		t.Errorf("Wrong FC logo: %s", logo)
	}
}

func TestGetLogoFromMenuConfigTga(t *testing.T) {
	logo, err := GetLogoFromMenuConfig("testdata/SomeMod/config/menu.cfg", []string{"/SomeMod"})

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logo path: %s", logo)

	if logo != "testdata/SomeMod/menu_logo.tga" {
		t.Errorf("Wrong FC logo: %s", logo)
	}
}

func TestGetConversionFromInit(t *testing.T) {
	fc, err := GetConversionFromInit("testdata/wn_config/main_init.cfg")
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

func TestGetMainInitConfigs(t *testing.T) {
	mainInits, err := GetMainInitConfigs("testdata")

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
