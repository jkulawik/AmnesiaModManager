package configs

import (
	"encoding/xml"
	"fmt"
	"os"
)

// TODO Improve: A lot of this data isn't really needed and could be removed if the unmarshaller still works afterwards

type MainInitXML struct {
	ConfigFiles MainInitXMLConfigFiles `xml:"ConfigFiles"`
	Directories MainInitXMLDirectories `xml:"Directories"`
	Variables   MainInitXMLVariables   `xml:"Variables"`
	StartMap    MainInitXMLStartMap    `xml:"StartMap"`
}

type MainInitXMLConfigFiles struct {
	Resources string `xml:"Resources,attr"`
	Materials string `xml:"Materials,attr"`

	Game    string `xml:"Game,attr"`
	Menu    string `xml:"Menu,attr"`
	PreMenu string `xml:"PreMenu,attr"`
	Demo    string `xml:"Demo,attr"`

	DefaultMainSettings     string `xml:"DefaultMainSettings,attr"`
	DefaultMainSettingsSDL2 string `xml:"DefaultMainSettingsSDL2,attr"`

	DefaultMainSettingsLow    string `xml:"DefaultMainSettingsLow,attr"`
	DefaultMainSettingsMedium string `xml:"DefaultMainSettingsMedium,attr"`
	DefaultMainSettingsHigh   string `xml:"DefaultMainSettingsHigh,attr"`

	DefaultUserSettings string `xml:"DefaultUserSettings,attr"`
	DefaultUserKeys     string `xml:"DefaultUserKeys,attr"`
	DefaultUserKeysSDL2 string `xml:"DefaultUserKeysSDL2,attr"`

	DefaultBaseLanguage string `xml:"DefaultBaseLanguage,attr"`
	DefaultGameLanguage string `xml:"DefaultGameLanguage,attr"`
}

type MainInitXMLDirectories struct {
	MainSaveFolder     string `xml:"MainSaveFolder,attr"`
	BaseLanguageFolder string `xml:"BaseLanguageFolder,attr"`
	GameLanguageFolder string `xml:"GameLanguageFolder,attr"`
	CustomStoryPath    string `xml:"CustomStoryPath,attr"`
}

type MainInitXMLVariables struct {
	GameName      string `xml:"GameName,attr"`
	AllowHardMode string `xml:"AllowHardMode,attr"`
}

type MainInitXMLStartMap struct {
	File   string `xml:"File,attr"`
	Folder string `xml:"Folder,attr"`
	Pos    string `xml:"Pos,attr"`
}

func ReadConversionInit(path string) (*MainInitXML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ReadConversionInit: %w", err)
	}
	// We need to wrap the config in a dummy tag to get it to unmarshal properly
	data = []byte("<dummy>" + string(data) + "</dummy>")

	mi := new(MainInitXML)
	empty := new(MainInitXML)
	err = xml.Unmarshal(data, mi)

	if *mi == *empty {
		return nil, fmt.Errorf("ReadConversionInit: XML parser returned an empty object with error: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("ReadConversionInit: %w", err)
	}
	return mi, nil
}

func (initXml *MainInitXML) GetLangName() string {
	return initXml.ConfigFiles.DefaultGameLanguage
}
