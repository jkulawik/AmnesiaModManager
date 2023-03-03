package main

import (
	"errors"
	"os"
)

func CheckIsRootDir(dir string) error {
	if dir == "" {
		dir = "."
	}

	filelist, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	coreInDir := false
	entInDir := false
	csInDir := false

	for _, direntry := range filelist {
		coreInDir = (coreInDir || direntry.Name() == "core")
		entInDir = (entInDir || direntry.Name() == "entities")
		csInDir = (csInDir || direntry.Name() == "custom_stories")
	}

	// Not checking for Amnesia binaries because their names differ between releases

	if !coreInDir || !entInDir {
		return errors.New("work directory is not an Amnesia install")
	}
	if !csInDir {
		return errors.New("custom story folder not found")
	}
	if coreInDir && entInDir && csInDir {
		return nil
	}

	return errors.New("unknown issue with the work directory") // func should never reach here but static analysis complains about not returning
}

func formatStringList(list []string) string {
	folderList := ""
	for _, f := range list {
		folderList += f + "\n"
	}
	return folderList
}

// For FC display card - stops it from shrinking
func getStringSpacer(width int) string {
	spacer := ""
	for i := 0; i < width; i++ {
		spacer += " "
	}
	return spacer
}
