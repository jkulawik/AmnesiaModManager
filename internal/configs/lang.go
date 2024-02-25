package configs

import (
	"encoding/xml"
	"fmt"
	"os"
)

type LangXML struct {
	XMLName    xml.Name          `xml:"LANGUAGE"`
	Categories []LangXMLCategory `xml:"CATEGORY"`
}

type LangXMLCategory struct {
	Name    string         `xml:"Name,attr"`
	Entries []LangXMLEntry `xml:"Entry"`
}

type LangXMLEntry struct {
	Name    string `xml:"Name,attr"`
	Content string `xml:",chardata"`
}

func GetDescFromLang(filepath string) (string, error) {

	data, err := os.ReadFile(filepath)
	if err != nil {
		return "Lang file not found", fmt.Errorf("GetDescFromLang: %w", err)
	}

	langxml := new(LangXML)
	err = xml.Unmarshal(data, langxml)
	// t.Log(langxml)

	if err != nil {
		return "Error while parsing lang file XML.", fmt.Errorf("GetDescFromLang: %w", err)
	}

	// Search categories
	var categoryCustomStoryMain *LangXMLCategory
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
	var entryCustomStoryDesc *LangXMLEntry
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
