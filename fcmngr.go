package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
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
			log.Fatal(err)
			fmt.Println(err, "in", path)
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
