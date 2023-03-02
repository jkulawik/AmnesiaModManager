package main

import (
	"encoding/xml"
	"errors"
	"io/fs"
	"os"
	"strings"
)

var baseResources = []string{
	"/_temp",
	"/fonts",
	"/maps",
	"/textures",
	"/models",
	"/gui",
	"/static_objects",
	"/sounds",
	"/main_menu",
	"/shaders",
	"/lights",
	"/billboards",
	"/entities",
	"/graphics",
	"/viewer",
	"/particles",
	"/models",
	"/music",
	"/flashbacks",
	"/textures",
	"/misc",
	"/commentary",
}

const mainInitStr = "main_init.cfg"

func GetMainInitConfigs(workdir string) ([]string, error) {
	fileSystem := os.DirFS(".")

	mainInits := make([]string, 0)

	fs.WalkDir(fileSystem, workdir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == mainInitStr {
			mainInits = append(mainInits, path)
		}

		if err != nil {
			ErrorLogger.Println(err, "in", path)
		}
		return nil
	})

	if len(mainInits) == 0 {
		return nil, errors.New("no full conversion init files found")
	} else {
		return mainInits, nil
	}
}

func ReadConversionInit(path string) (*MainInitXML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// We need to wrap the config in a dummy tag to get it to unmarshal properly
	data = []byte("<dummy>" + string(data) + "</dummy>")

	mi := new(MainInitXML)
	empty := new(MainInitXML)
	err = xml.Unmarshal(data, mi)

	if *mi == *empty {
		return nil, errors.New(path + ": XML parser returned an empty object")
	}

	if err != nil {
		return nil, err
	} else {
		return mi, nil
	}
}

func GetUniqueResources(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	res := new(ResourcesXML)
	err = xml.Unmarshal(data, res)

	if len(res.Directory) == 0 {
		return nil, errors.New(path + ": XML parser returned an empty object")
	}

	if err != nil {
		return nil, err
	} else {
		resFolders := make([]string, 0)
		for _, entry := range res.Directory {
			// Don't include folders of the base game
			isInBaseResources := false

			for _, baseResource := range baseResources {
				if entry.Path == baseResource {
					isInBaseResources = true
					break
				}
			}

			if !isInBaseResources {
				resFolders = append(resFolders, entry.Path)
			}

		}

		return resFolders, nil
	}
}

func GetLogoFromMenuConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			WarningLogger.Println(path, ": specified menu.cfg file doesn't exist")
			return "", nil
		} else {
			return "", err
		}
	}

	// We need to wrap the config in a dummy tag to get it to unmarshal properly
	data = []byte("<dummy>" + string(data) + "</dummy>")

	menu := new(MenuXML)
	empty := new(MenuXML)
	err = xml.Unmarshal(data, menu)

	if *menu == *empty {
		return "", errors.New(path + ": XML parser returned an empty object")
	}

	if err != nil {
		return "", err
	} else {
		return menu.Main.MenuLogo, nil
	}
}

func GetConversionFromInit(path string) (*FullConversion, error) {
	init, err := ReadConversionInit(path)
	if err != nil {
		return nil, err
	}

	fc := new(FullConversion)
	fc.name = init.Variables.GameName
	fc.mainInitConfig = path
	res, err := GetUniqueResources("testdata/" + init.ConfigFiles.Resources)
	if err != nil {
		return nil, err
	}
	fc.uniqueResources = res
	logo, err := GetLogoFromMenuConfig("testdata/" + init.ConfigFiles.Menu)
	if err != nil {
		return nil, err
	}
	fc.logo = logo

	return fc, nil
}

func GetFullConversions(workdir string) ([]*FullConversion, error) {

	initList, err := GetMainInitConfigs(workdir)

	if err != nil {
		return nil, err
	}
	t.Log(initList)

	fcList := make([]*FullConversion, 0, len(initList))

	for _, init := range initList {
		fc, err := GetConversionFromInit(init)

		if err != nil {
			return nil, err
		}

		fcList = append(fcList, fc)
	}

	if len(fcList) == 0 {
		return nil, errors.New("did not find any full conversions")
	}

	return fcList, nil
}
