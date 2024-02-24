package mods

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"modmanager/internal/configs"
	"modmanager/internal/logger"
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

func GetMainInitConfigs(path string) ([]string, error) {
	fileSystem := os.DirFS(path)

	mainInits := make([]string, 0)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error.Println(err, "in", path)
			return fmt.Errorf("GetMainInitConfigs: %w", err)
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

func ReadConversionInit(path string) (*configs.MainInitXML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ReadConversionInit: %w", err)
	}
	// We need to wrap the config in a dummy tag to get it to unmarshal properly
	data = []byte("<dummy>" + string(data) + "</dummy>")

	mi := new(configs.MainInitXML)
	empty := new(configs.MainInitXML)
	err = xml.Unmarshal(data, mi)

	if *mi == *empty {
		return nil, fmt.Errorf("ReadConversionInit: XML parser returned an empty object with error: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("ReadConversionInit: %w", err)
	} else {
		return mi, nil
	}
}

func GetUniqueResources(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GetUniqueResources: %w", err)
	}

	res := new(configs.ResourcesXML)
	err = xml.Unmarshal(data, res)

	if len(res.Directory) == 0 {
		return nil, fmt.Errorf("GetUniqueResources: XML parser returned an empty object with error: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("GetUniqueResources: %w", err)
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

		for i, entry := range resFolders {
			if string(entry[0]) == "/" {
				resFolders[i] = entry[1:]
			}
		}

		return resFolders, nil
	}
}

func GetLogoFromMenuConfig(filepath string, resources []string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error while reading menu config: %w", err)
	}

	// We need to wrap the config in a dummy tag to get it to unmarshal properly
	data = []byte("<dummy>" + string(data) + "</dummy>")

	menu := new(configs.MenuXML)
	empty := new(configs.MenuXML)
	err = xml.Unmarshal(data, menu)

	var searchName string
	if menu.Main.MenuLogo == "" {
		logger.Info.Println(filepath, "mod doesn't specify a logo. Trying default name")
		searchName = "menu_logo.tga"
	} else if err != nil {
		logger.Error.Println(err)
		searchName = "menu_logo.tga"
	} else if *menu == *empty {
		logger.Warn.Println("GetLogoFromMenuConfig: " + filepath + ": XML parser returned an empty object. Trying default name")
		searchName = "menu_logo.tga"
	} else {
		searchName = menu.Main.MenuLogo
	}

	// Find the logo path
	fileSystem := os.DirFS(".")
	logoCandidates := make([]string, 0)

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error.Println(err)
			return fmt.Errorf("GetLogoFromMenuConfig: %w occured while crawling the filesystem", err)
		}
		if !d.IsDir() && d.Name() == searchName {
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
		searchRoot := res
		if string(res[0]) == "/" {
			searchRoot = res[1:]
		}
		fs.WalkDir(fileSystem, searchRoot, walkFunc)
	}

	// Search base game folders as a last ditch resort; this will search some dirs again so walkFunc checks for doubles
	fs.WalkDir(fileSystem, ".", walkFunc)

	logger.Info.Println("Logo candidates:", logoCandidates)
	if len(logoCandidates) == 0 {
		return "", errors.New("mod logo could not be found")
	} else {
		return logoCandidates[0], nil
	}
}

func GetConversionFromInit(path string) (*FullConversion, error) {
	init, err := ReadConversionInit(path)
	if err != nil {
		return nil, fmt.Errorf("GetConversionFromInit: %w", err)
	}

	fc := new(FullConversion)
	fc.Name = init.Variables.GameName
	fc.MainInitConfig = path
	res, err := GetUniqueResources(init.ConfigFiles.Resources)
	if err != nil {
		return nil, fmt.Errorf("GetConversionFromInit: %w", err)
	}
	fc.UniqueResources = res
	menuPath := init.ConfigFiles.Menu
	logo, err := GetLogoFromMenuConfig(menuPath, res)
	if err != nil {
		logger.Warn.Println("Unexpected error while searching for logo, no logo will be used. Error:", err)
		logo = ""
	}
	fc.Logo = logo

	return fc, nil
}

func GetFullConversions(path string) ([]*FullConversion, error) {

	initList, err := GetMainInitConfigs(path)

	if err != nil {
		return nil, fmt.Errorf("GetFullConversions: %w", err)
	}

	// Find and remove base game init

	for i, init := range initList {
		if init == "config/main_init.cfg" {
			initList = append(initList[:i], initList[i+1:]...)
			break
		}
	}
	logger.Info.Println("Found main init configs:", initList)

	fcList := make([]*FullConversion, 0, len(initList))

	for _, init := range initList {
		fc, err := GetConversionFromInit(init)

		if err != nil {
			logger.Error.Println("Error while reading full conversion from", init, "-", err)
		}

		if fc != nil {
			fcList = append(fcList, fc)
		}
	}

	if len(fcList) == 0 {
		return nil, errors.New("did not find any valid full conversions")
	}

	return fcList, nil
}
