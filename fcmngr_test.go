package main

import "testing"

func TestGetMainInitConfigs(t *testing.T) {
	mainInits, err := GetMainInitConfigs("testdata")

	if err != nil {
		t.Error(err)
	}

	if mainInits[0] != "testdata/SomeMod/config/main_init.cfg" {
		t.Errorf("Did not find SomeMod main init. Instead got %s", mainInits[0])
	}
}

// func TestGetConversionFromInit(t *testing.T) {
// 	fullConversion, err := GetConversionFromInit("./testdata/SomeMod/main_init.cfg")
// }
