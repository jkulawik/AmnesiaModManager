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

func TestGetConversionFromInit(t *testing.T) {
}
