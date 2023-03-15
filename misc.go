package main

import (
	"errors"
	"os"
)

var execMap = map[string]string{
	"windows": ".\\Amnesia.exe",
	"linux":   "./Amnesia.bin.x86_64",
}

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

	// TODO Check for Amnesia binaries?

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

func isModNil(mod Mod) bool {
	// assigning structs which implement interfaces which are nil is not the same as assigning nil;
	// this means that interface == nil will return false in such cases
	csNil := (*CustomStory)(nil)
	fcNil := (*FullConversion)(nil)
	return mod == nil || mod == csNil || mod == fcNil
}

func deleteModDir(path string) error {
	InfoLogger.Println("trying to delete:", path)

	// Check if img file exists
	if _, err := os.Stat(path); err != nil {
		ErrorLogger.Println(err)
		return err
	}

	// There should be no trailing slashes anywhere and we need to add one for the deletion to succeed
	lastChar := path[len(path)-1:]
	if lastChar != "/" {
		path += "/"
	}

	err := os.RemoveAll(path)
	return err
}
