package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
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

	if !csInDir {
		return errors.New("custom story folder not found")
	}
	if !coreInDir || !entInDir {
		return errors.New("work directory is not an Amnesia install")
	}
	if coreInDir && entInDir && csInDir {
		return nil
	}

	return errors.New("unknown issue with the work directory") // func should never reach here but static analysis complains about not returning
}

func ReadCustomStoryConfig(filepath string) (*CSXML, error) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	csxml := new(CSXML)
	err = xml.Unmarshal(data, csxml)

	if err != nil {
		return nil, err
	} else {
		return csxml, nil
	}
}

func GetDescFromLang(filepath string) (string, error) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return "Lang file not found", err
	}

	langxml := new(LangXML)
	err = xml.Unmarshal(data, langxml)
	// t.Log(langxml)

	if err != nil {
		return "Error while parsing lang file XML.", err
	}

	// Search categories
	var categoryCustomStoryMain LangXMLCategory
	isMainCategoryInLang := false

	for _, cat := range langxml.Categories {
		if cat.Name == "CustomStoryMain" {
			categoryCustomStoryMain = cat
			isMainCategoryInLang = true
			break
		}
	}
	if !isMainCategoryInLang {
		return "This custom story has no description (missing CustomStoryMain category).", nil
	}
	//t.Log(categoryCustomStoryMain)

	// Search category for entry
	var entryCustomStoryDesc LangXMLEntry
	isDescInCategory := false

	for _, entry := range categoryCustomStoryMain.Entries {
		if entry.Name == "Description" {
			entryCustomStoryDesc = entry
			isDescInCategory = true
			break
		}
	}
	if !isDescInCategory {
		return "This custom story has no description (missing description entry).", nil
	}

	return entryCustomStoryDesc.Content, nil
}

func makeInvalidStory(dir string) *CustomStory {
	EmptyDirStory := new(CustomStory)

	EmptyDirStory.name = "Invalid custom story"
	EmptyDirStory.author = "N/A"
	EmptyDirStory.dir = dir
	EmptyDirStory.desc = "A folder was found which is not a valid custom story."

	return EmptyDirStory
}

func GetStoryFromDir(dir string) (*CustomStory, error) {

	csxml, err := ReadCustomStoryConfig(dir + "custom_story_settings.cfg")
	cs := new(CustomStory)

	if err != nil {
		if strings.Contains(err.Error(), "custom_story_settings.cfg: no such file or directory") {
			cs = makeInvalidStory(dir)
			return cs, nil
		} else {
			return nil, err
		}
	}

	cs.dir = dir
	cs.author = csxml.Author
	cs.name = csxml.Name
	cs.imgFile = csxml.ImgFile

	if _, err := os.Stat(cs.dir + cs.imgFile); err != nil {
		cs.imgFile = ""
	}

	if csxml.ExtraLangFilePrefix != "" {
		if csxml.DefaultExtraLangFile == "" {
			cs.langFile = csxml.ExtraLangFilePrefix + "english.lang"
		} else {
			cs.langFile = csxml.ExtraLangFilePrefix + csxml.DefaultExtraLangFile
		}
	} else if csxml.DefaultExtraLangFile != "" {
		cs.langFile = "extra_" + csxml.DefaultExtraLangFile
	} else {
		cs.langFile = "extra_english.lang"
	}

	cs.desc, err = GetDescFromLang(dir + cs.langFile)

	if err != nil {
		fmt.Println(cs, err)
		return cs, err
	}

	return cs, nil
}

func GetCustomStories(dir string) ([]*CustomStory, error) {
	filelist, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	csList := make([]*CustomStory, 0, len(filelist))

	for _, direntry := range filelist {
		if direntry.IsDir() {
			path := dir + "/" + direntry.Name() + "/"
			cs, err := GetStoryFromDir(path)

			if cs != nil {
				csList = append(csList, cs)
			}
			if err != nil {
				fmt.Println(err)
				//Can't return nil dueo to an error because finding one doesn't mean the entire list is invalid
			}
		}
	}

	if len(csList) == 0 {
		return nil, errors.New("did not find any folders in the custom story directory")
	}

	return csList, nil
}
