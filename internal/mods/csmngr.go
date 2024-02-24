package mods

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"modmanager/internal/configs"
	"modmanager/internal/logger"
)

func MakeStoryText(cs *CustomStory) string {
	return fmt.Sprintf("Folder:\n%s\nDescription:\n%s", cs.Dir, cs.Desc)
}

func ReadCustomStoryConfig(filepath string) (*configs.CSXML, error) {
	fileSystem := os.DirFS(".")
	// data, err := os.ReadFile(filepath)
	data, err := fs.ReadFile(fileSystem, filepath)
	if err != nil {
		return nil, fmt.Errorf("ReadCustomStoryConfig: %w", err)
	}

	csxml := new(configs.CSXML)
	// empty := new(configs.CSXML)
	err = xml.Unmarshal(data, csxml)

	// if *csxml == *empty {
	// }

	if err != nil {
		return nil, fmt.Errorf("ReadCustomStoryConfig: XML parser encountered an error: %w", err)
	}
	return csxml, nil
}

func GetDescFromLang(filepath string) (string, error) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return "Lang file not found", fmt.Errorf("GetDescFromLang: %w", err)
	}

	langxml := new(configs.LangXML)
	err = xml.Unmarshal(data, langxml)
	// t.Log(langxml)

	if err != nil {
		return "Error while parsing lang file XML.", fmt.Errorf("GetDescFromLang: %w", err)
	}

	// Search categories
	var categoryCustomStoryMain *configs.LangXMLCategory
	for _, cat := range langxml.Categories {
		if cat.Name == "CustomStoryMain" {
			categoryCustomStoryMain = &cat
			break
		}
	}
	if categoryCustomStoryMain == nil {
		return "This custom story has no description (missing CustomStoryMain category).", nil
	}
	//t.Log(categoryCustomStoryMain)

	// Search category for entry
	var entryCustomStoryDesc *configs.LangXMLEntry
	for _, entry := range categoryCustomStoryMain.Entries {
		if entry.Name == "Description" {
			entryCustomStoryDesc = &entry
			break
		}
	}
	if entryCustomStoryDesc == nil {
		return "This custom story has no description (missing description entry).", nil
	}

	return entryCustomStoryDesc.Content, nil
}

func makeInvalidStory(dir string) *CustomStory {
	EmptyDirStory := new(CustomStory)

	EmptyDirStory.Name = "Invalid custom story"
	EmptyDirStory.Author = "N/A"
	EmptyDirStory.Dir = dir
	EmptyDirStory.Desc = "A folder was found which is not a valid custom story."

	return EmptyDirStory
}

func GetStoryFromDir(dir string) (*CustomStory, error) {
	csxml, err := ReadCustomStoryConfig(dir + "/custom_story_settings.cfg")
	cs := new(CustomStory)

	if err != nil {
		if strings.Contains(err.Error(), "custom_story_settings.cfg: no such file or directory") {
			logger.Warn.Println(err)
			cs = makeInvalidStory(dir)
			return cs, nil
		}
		return nil, fmt.Errorf("GetStoryFromDir: %w", err)
	}

	cs.Dir = dir
	cs.Author = csxml.Author
	cs.Name = csxml.Name
	cs.ImgFile = csxml.ImgFile

	// Check if img file exists
	if _, err := os.Stat(cs.Dir + "/" + cs.ImgFile); err != nil {
		logger.Warn.Println(err)
		cs.ImgFile = ""
	}

	if csxml.ExtraLangFilePrefix != "" {
		if csxml.DefaultExtraLangFile == "" {
			cs.LangFile = csxml.ExtraLangFilePrefix + "english.lang"
		} else {
			cs.LangFile = csxml.ExtraLangFilePrefix + csxml.DefaultExtraLangFile
		}
	} else if csxml.DefaultExtraLangFile != "" {
		cs.LangFile = "extra_" + csxml.DefaultExtraLangFile
	} else {
		cs.LangFile = "extra_english.lang"
	}

	cs.Desc, err = GetDescFromLang(dir + "/" + cs.LangFile)

	if err != nil {
		if err.Error() == "XML syntax error on line 3: invalid sequence \"--\" not allowed in comments" { // TODO use string contains here
			logger.Warn.Println(cs.Dir, err)
			return cs, nil
		}
		return cs, fmt.Errorf("GetStoryFromDir: %w", err)
	}

	return cs, nil
}

func GetCustomStories(dir string) ([]*CustomStory, error) {
	filelist, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("GetCustomStories: %w", err)
	}

	csList := make([]*CustomStory, 0, len(filelist))
	for _, direntry := range filelist {
		if !direntry.IsDir() {
			continue
		}
		cs, err := GetStoryFromDir(dir + "/" + direntry.Name())

		if err != nil {
			logger.Error.Println("GetCustomStories: ", err)
			// Can't return nil due to an error because finding one doesn't mean the entire list is invalid
		}
		if cs != nil {
			csList = append(csList, cs)
		}
	}

	if len(csList) == 0 {
		return nil, errors.New("did not find any folders in the custom story directory")
	}

	return csList, nil
}
