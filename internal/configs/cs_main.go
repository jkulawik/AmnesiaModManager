package configs

import (
	"encoding/xml"
	"fmt"
	"os"
)

type CSXML struct {
	ImgFile              string `xml:"ImgFile,attr"`
	Name                 string `xml:"Name,attr"`
	Author               string `xml:"Author,attr"`
	ExtraLangFilePrefix  string `xml:"ExtraLangFilePrefix,attr"`
	DefaultExtraLangFile string `xml:"DefaultExtraLangFile,attr"`
	InitCfgFile          string `xml:"InitCfgFile,attr"`
}

func ReadCustomStoryConfig(filepath string) (*CSXML, error) {
	// fileSystem := os.DirFS(".")
	// data, err := fs.ReadFile(fileSystem, filepath)
	data, err := os.ReadFile(filepath) // this is better because it supports going up a directory, which is needed for Steam workshop
	if err != nil {
		return nil, fmt.Errorf("ReadCustomStoryConfig: %w", err)
	}

	csxml := new(CSXML)
	// empty := new(configs.CSXML)
	err = xml.Unmarshal(data, csxml)

	// if *csxml == *empty {
	// }

	if err != nil {
		return nil, fmt.Errorf("ReadCustomStoryConfig: XML parser encountered an error: %w", err)
	}
	return csxml, nil
}

func (csxml *CSXML) GetLangName() string {
	if csxml.ExtraLangFilePrefix != "" {
		if csxml.DefaultExtraLangFile == "" {
			return csxml.ExtraLangFilePrefix + "english.lang"
		} else {
			return csxml.ExtraLangFilePrefix + csxml.DefaultExtraLangFile
		}
	} else if csxml.DefaultExtraLangFile != "" {
		return "extra_" + csxml.DefaultExtraLangFile
	} else {
		return "extra_english.lang"
	}
}
