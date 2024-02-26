package configs

import (
	"encoding/xml"
	"fmt"
	"os"
)

const replacement_byte = byte(0x3D) // ASCII "="

type LangXML struct {
	// turns out the marshaller doesn't need the root tag to get its contents;
	// this allows us to handle lang files with non-standard root tag names (extremely rare but oh well)
	// XMLName    xml.Name          `xml:"LANGUAGE"`
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

	// Swap invalid -- sequences inside comments to ==
	isInsideComment := false
	for i := 0; i < len(data)-3; i++ {
		if string(data[i:i+4]) == "<!--" {
			if isInsideComment {
				// we found a <!-- inside a comment
				data[i+2] = replacement_byte
				data[i+3] = replacement_byte
				// logger.Info.Printf("Found <!-- inside comment, replaced with [%s]\n", string(data[i:i+4]))
				continue
			} else {
				isInsideComment = true
				continue // this is not enough to stop replacing -- in a valid <!-- token, because !-- will still get picked up
			}
		}
		if string(data[i:i+3]) == "-->" {
			isInsideComment = false
			continue
		}
		// Make sure we can access i-2 bytes
		if i < 2 {
			continue
		}
		if isInsideComment && string(data[i-2:i]) != "<!" && string(data[i:i+2]) == "--" {
			// logger.Info.Printf("found wrong sequence -- at position %d: [%s]\n", i, data[i-2:i+3])
			data[i] = replacement_byte
			data[i+1] = replacement_byte
		}
	}

	langxml := new(LangXML)
	err = xml.Unmarshal(data, langxml)

	// Note: ignoring the -- in comments error does not help because the data won't be unmarshalled anyway.
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
