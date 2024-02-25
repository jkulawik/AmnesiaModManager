package configs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ResourcesXML struct {
	XMLName   xml.Name                `xml:"Resources"`
	Directory []ResourcesXMLDirectory `xml:"Directory"`
}

type ResourcesXMLDirectory struct {
	Path       string `xml:"Path,attr"`
	AddSubDirs string `xml:"AddSubDirs,attr"`
}

var baseResources = map[string]bool{
	"/_temp":          true,
	"/fonts":          true,
	"/maps":           true,
	"/textures":       true,
	"/models":         true,
	"/gui":            true,
	"/static_objects": true,
	"/sounds":         true,
	"/main_menu":      true,
	"/shaders":        true,
	"/lights":         true,
	"/billboards":     true,
	"/entities":       true,
	"/graphics":       true,
	"/viewer":         true,
	"/particles":      true,
	"/music":          true,
	"/flashbacks":     true,
	"/misc":           true,
	"/commentary":     true,
}

func GetUniqueResourceDirs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GetUniqueResourceDirs: %w", err)
	}

	res := new(ResourcesXML)
	err = xml.Unmarshal(data, res)

	if err != nil {
		return nil, fmt.Errorf("GetUniqueResourceDirs: %w", err)
	}
	if len(res.Directory) == 0 {
		return nil, errors.New("GetUniqueResourceDirs: no directories found in resources file")
	}

	resourceDirs := make([]string, 0)
	for _, entry := range res.Directory {
		// Clean up Windows filth
		clean_path := strings.ReplaceAll(entry.Path, "\\", "/")
		// Don't include folders of the base game
		_, isBase := baseResources[clean_path]
		if !isBase {
			resourceDirs = append(resourceDirs, clean_path)
		}
	}

	for i, entry := range resourceDirs {

		// Remove root slash
		if string(entry[0]) == "/" {
			resourceDirs[i] = entry[1:]
		}
	}

	return resourceDirs, nil
}
