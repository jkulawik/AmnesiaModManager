package mods

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"modmanager/internal/configs"
	"modmanager/internal/logger"
)

func GetMainInitFilepaths(path string) ([]string, error) {
	fileSystem := os.DirFS(path)

	mainInits := make([]string, 0)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error.Println(err, "in", path)
			return fmt.Errorf("GetMainInitFilepaths: %w", err)
		}
		// Ignore base game init
		if !d.IsDir() && d.Name() == "main_init.cfg" && path != "config/main_init.cfg" {
			mainInits = append(mainInits, path)
		}
		return nil
	})

	if len(mainInits) == 0 {
		return nil, errors.New("no full conversion init files found")
	}
	return mainInits, nil
}

func GetConversionFromInit(workdir, path string) (*FullConversion, error) {
	init, err := configs.ReadConversionInit(path)
	if err != nil {
		return nil, fmt.Errorf("GetConversionFromInit (%s): %w", path, err)
	}

	fc := new(FullConversion)
	fc.Name = init.Variables.GameName
	fc.MainInitConfig = path
	res, err := configs.GetUniqueResourceDirs(workdir + "/" + init.ConfigFiles.Resources)
	if err != nil {
		return nil, fmt.Errorf("GetConversionFromInit (%s): %w", path, err)
	}
	fc.UniqueResources = res
	logo, err := configs.GetLogoPathFromMenuConfig(workdir+"/"+init.ConfigFiles.Menu, res)
	if err != nil {
		logger.Warn.Printf("GetConversionFromInit (%s): %s\n", path, err)
		logo = ""
	}
	fc.Logo = logo
	fc.LangFile = init.GetLangName()

	return fc, nil
}

func GetFullConversions(workdir string) ([]*FullConversion, error) {
	initList, err := GetMainInitFilepaths(workdir)

	if err != nil {
		return nil, fmt.Errorf("GetFullConversions: %w", err)
	}
	if len(initList) == 0 {
		return nil, errors.New("GetFullConversions: did not find any main_init files")
	}

	logger.Info.Println("Found main init configs:", initList)

	fcList := make([]*FullConversion, 0, len(initList))
	for _, init := range initList {
		fc, err := GetConversionFromInit(workdir, workdir+"/"+init)
		if err != nil {
			logger.Error.Println("Error while reading full conversion from", init, ":", err)
			continue
		}
		fcList = append(fcList, fc)
	}

	if len(fcList) == 0 {
		return nil, errors.New("GetFullConversions: did not find any valid full conversions")
	}

	return fcList, nil
}
