package configs

import (
	"fmt"
	"testing"
)

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

func TestGetLogoFromMenuConfigTga(t *testing.T) {
	logo, err := GetLogoPathFromMenuConfig("testdata/SomeMod/config/menu.cfg", []string{"/SomeMod"})

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Logo path: %s", logo)

	if logo != "testdata/SomeMod/menu_logo.tga" {
		t.Errorf("Wrong FC logo: %s", logo)
	}
}

func TestGetUniqueResourceDirs(t *testing.T) {
	res, err := GetUniqueResourceDirs("testdata/SomeMod/config/resources.cfg")

	if err != nil {
		t.Fatal(err)
	}

	resString := fmt.Sprintf("%v", res)
	testString := fmt.Sprintf("%v", []string{"SomeMod", "SomeMod/misc"})

	if resString != testString {
		t.Errorf("Parsed resource list differs expected. Original:\n%s\nGot:\n%s", testString, resString)
	}

}
