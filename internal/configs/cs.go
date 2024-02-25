package configs

import (
	"encoding/xml"
	"fmt"
	"io/fs"
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
	fileSystem := os.DirFS(".")
	// data, err := os.ReadFile(filepath)
	data, err := fs.ReadFile(fileSystem, filepath)
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
