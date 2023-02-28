package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
)

var originalResources = []string{
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
