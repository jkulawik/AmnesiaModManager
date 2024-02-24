package configs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"modmanager/internal/logger"
	"os"
)

type MenuXML struct {
	Main MenuXMLMain `xml:"Main"`
}

type MenuXMLMain struct {
	MainFadeInTime      string `xml:"MainFadeInTime,attr"`
	MainFadeOutTimeFast string `xml:"MainFadeOutTimeFast,attr"`
	MainFadeOutTimeSlow string `xml:"MainFadeOutTimeSlow,attr"`

	TopMenuFadeInTime             string `xml:"TopMenuFadeInTime,attr"`
	TopMenuFadeOutTime            string `xml:"TopMenuFadeOutTime,attr"`
	TopMenuFontRelativeSize       string `xml:"TopMenuFontRelativeSize,attr"`
	TopMenuStartRelativePos       string `xml:"TopMenuStartRelativePos,attr"`
	TopMenuStartRelativePosInGame string `xml:"TopMenuStartRelativePosInGame,attr"`
	TopMenuFont                   string `xml:"TopMenuFont,attr"`

	MainMenuLogoStartRelativePos string `xml:"MainMenuLogoStartRelativePos,attr"`
	MainMenuLogoRelativeSize     string `xml:"MainMenuLogoRelativeSize,attr"`

	BGScene            string `xml:"BGScene,attr"`
	BGCamera_FOV       string `xml:"BGCamera_FOV,attr"`
	BGCamera_ZoomedFOV string `xml:"BGCamera_ZoomedFOV,attr"`

	ZoomSound string `xml:"ZoomSound,attr"`
	Music     string `xml:"Music,attr"`

	MenuLogo string `xml:"MenuLogo,attr"`
}

func GetLogoPathFromMenuConfig(filepath string, resources []string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("error while reading menu config: %w", err)
	}

	// We need to wrap the config in a dummy tag to get it to unmarshal properly
	data = []byte("<dummy>" + string(data) + "</dummy>")

	menu := new(MenuXML)
	empty := new(MenuXML)
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
	}
	return logoCandidates[0], nil
}
