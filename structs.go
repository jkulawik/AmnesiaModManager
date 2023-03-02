package main

import "encoding/xml"

type Mod interface {
	// Need basically only to allow deletion of both CS and FC mods
	listFolders() []string
}

type CustomStory struct {
	name     string
	author   string
	langFile string
	dir      string
	desc     string
	imgFile  string
}

func (cs CustomStory) listFolders() []string {
	return []string{cs.dir}
}

type FullConversion struct {
	name            string
	mainInitConfig  string
	resourcesConfig string
	logo            string
}

var _ Mod = (*CustomStory)(nil) // Check if CS implements interface (at compile time)

// ------ Custom story config XML ------ //

type CSXML struct {
	ImgFile              string `xml:"ImgFile,attr"`
	Name                 string `xml:"Name,attr"`
	Author               string `xml:"Author,attr"`
	ExtraLangFilePrefix  string `xml:"ExtraLangFilePrefix,attr"`
	DefaultExtraLangFile string `xml:"DefaultExtraLangFile,attr"`
}

// ------ Lang file XML ------ //

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

// ------ FC main init XML ------ //

// Improve: A lot of this data isn't really needed and could be removed if the unmarshaller still works afterwards

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

// ------ Resources.cfg XML ------ //

type ResourcesXML struct {
	XMLName   xml.Name                `xml:"Resources"`
	Directory []ResourcesXMLDirectory `xml:"Directory"`
}

type ResourcesXMLDirectory struct {
	Path       string `xml:"Path,attr"`
	AddSubDirs string `xml:"AddSubDirs,attr"`
}

// ------ Menu.cfg XML ------ //

type MenuXML struct {
	Main MenuXMLMain `xml:"Main"`
}

type MenuXMLMain struct {
	MainFadeInTime      string `xml:"MainFadeInTime,attr"`
	MainFadeOutTimeFast string `xml:"MainFadeOutTimeFast,attr"`
	MainFadeOutTimeSlow string `xml:"MainFadeOutTimeSlow,attr"`

	TopMenuFadeInTime             string `xml:"TopMenuFadeInTime,attr"`
	TopMenuFadeOutTime            string `xml:"TopMenuFadeOutTime,attr"`
	TopMenuFontRelativeSize       string `xml:"TopMenuFontRelativeSize,attr"`
	TopMenuStartRelativePos       string `xml:"TopMenuStartRelativePos,attr"`
	TopMenuStartRelativePosInGame string `xml:"TopMenuStartRelativePosInGame,attr"`
	TopMenuFont                   string `xml:"TopMenuFont,attr"`

	MainMenuLogoStartRelativePos string `xml:"MainMenuLogoStartRelativePos,attr"`
	MainMenuLogoRelativeSize     string `xml:"MainMenuLogoRelativeSize,attr"`

	BGScene            string `xml:"BGScene,attr"`
	BGCamera_FOV       string `xml:"BGCamera_FOV,attr"`
	BGCamera_ZoomedFOV string `xml:"BGCamera_ZoomedFOV,attr"`

	ZoomSound string `xml:"ZoomSound,attr"`
	Music     string `xml:"Music,attr"`

	MenuLogo string `xml:"MenuLogo,attr"`
}
