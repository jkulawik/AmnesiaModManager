package main

import "testing"

var TestWhiteNight = FullConversion{
	name:            "White Night",
	mainInitConfig:  "",
	resourcesConfig: "",
	logo:            "wn_menu_logo.png",
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

func TestReadMenuConfig(t *testing.T) {
	menu, err := ReadMenuConfig("testdata/wn_config/menu.cfg")
	t.Logf("Menu config: %s", menu)

	if err != nil {
		t.Error(err)
	}

	if menu != nil && menu.MenuLogo != "wn_menu_logo.png" {
		t.Errorf("Wrong FC name: %s", menu.MenuLogo)
	}
}

func TestGetLogoFromConfig(t *testing.T) {
	t.Error("unimplemented")
}

func TestGetConversionFromInit(t *testing.T) {
	t.Error("unimplemented")
}

func TestGetMainInitConfigs(t *testing.T) {
	mainInits, err := GetMainInitConfigs("testdata")

	if err != nil {
		t.Error(err)
	}

	if mainInits[0] != "testdata/SomeMod/config/main_init.cfg" && mainInits[1] != "testdata/wn_config/main_init.cfg" {
		t.Errorf("Did not find one of the main inits. Got: %s", mainInits)
	}
}

func TestGetConversions(t *testing.T) {
	t.Error("unimplemented")
}
