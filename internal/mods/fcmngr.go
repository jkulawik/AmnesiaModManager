package mods

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

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

	// logger.Info.Println(init.Variables.GameName, "workdir:", workdir)
	// Note: hybrid mods will show up in this search but their workdir is the root dir, so their resources.cfg path won't be valid

	// Clean up Windows filth
	clean_res_path := strings.ReplaceAll(init.ConfigFiles.Resources, "\\", "/")
	res, err := configs.GetUniqueResourceDirs(workdir + "/" + clean_res_path)
	if err != nil {
		return nil, fmt.Errorf("GetConversionFromInit (%s): %w", path, err)
	}
	fc.UniqueResources = res

	// Clean up Windows filth
	clean_logo_path := strings.ReplaceAll(init.ConfigFiles.Menu, "\\", "/")
	logo, err := configs.GetLogoPathFromMenuConfig(workdir+"/"+clean_logo_path, res)
	if err != nil {
		logger.Warn.Printf("GetConversionFromInit (%s): %s\n", path, err)
		logo = ""
	}
	fc.Logo = logo
	// Clean up Windows filth
	fc.LangFile = strings.ReplaceAll(init.GetLangName(), "\\", "/")

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
