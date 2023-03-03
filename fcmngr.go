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

func GetMainInitConfigs() ([]string, error) {
	fileSystem := os.DirFS(".")

	mainInits := make([]string, 0)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			ErrorLogger.Println(err, "in", path)
			return err
		}
		if !d.IsDir() && d.Name() == mainInitStr {
			mainInits = append(mainInits, path)
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

func GetLogoFromMenuConfig(filepath string, resources []string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			WarningLogger.Println("GetLogoFromMenuConfig: file", filepath, "doesn't exist")
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
	if err != nil {
		return "", err
	}

	if *menu == *empty {
		return "", errors.New("GetLogoFromMenuConfig: " + filepath + ": XML parser returned an empty object")
	}

	// Find the logo path
	fileSystem := os.DirFS(".")
	logoCandidates := make([]string, 0)

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			ErrorLogger.Println(err)
			return err
		}
		if !d.IsDir() && d.Name() == menu.Main.MenuLogo {
			alreadyFound := false
			for _, candidate := range logoCandidates {
				if path == candidate {
					alreadyFound = true
					break
				}
			}
			if !alreadyFound {
				logoCandidates = append(logoCandidates, path)
			}
		}
		return nil
	}

	// Search custom resource dirs first
	// (same name as vanilla logo could have been used, so we can't just search from root once)
	for _, res := range resources {
		fs.WalkDir(fileSystem, res, walkFunc)
	}

	// Search base game folders as a last ditch resort; this will search some dirs again so walkFunc checks for doubles
	fs.WalkDir(fileSystem, ".", walkFunc)

	InfoLogger.Println("Logo candidates:", logoCandidates)
	if len(logoCandidates) == 0 {
		return "", errors.New("mod logo could not be found")
	} else {
		return logoCandidates[0], nil
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
	res, err := GetUniqueResources(init.ConfigFiles.Resources)
	if err != nil {
		return nil, err
	}
	fc.uniqueResources = res
	menuPath := init.ConfigFiles.Menu
	logo, err := GetLogoFromMenuConfig(menuPath, res)
	if err != nil {
		WarningLogger.Println("Error while searching for logo:", err)
	}
	fc.logo = logo

	return fc, nil
}

func GetFullConversions() ([]*FullConversion, error) {

	initList, err := GetMainInitConfigs()

	if err != nil {
		return nil, err
	}

	// Find and remove base game init

	for i, init := range initList {
		if init == "config/main_init.cfg" {
			initList = append(initList[:i], initList[i+1:]...)
			break
		}
	}
	InfoLogger.Println("Found main init configs:", initList)

	fcList := make([]*FullConversion, 0, len(initList))

	for _, init := range initList {
		fc, err := GetConversionFromInit(init)

		if err != nil {
			ErrorLogger.Println("Error while reading full conversion from", init, "-", err)
		}

		fcList = append(fcList, fc)
	}

	if len(fcList) == 0 {
		return nil, errors.New("did not find any full conversions")
	}

	return fcList, nil
}
