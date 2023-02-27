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

var _ Mod = (*CustomStory)(nil) // Check if CS implements interface (at compile time)

// ------ HPL2 XML parsing ------ //

type CSXML struct {
	ImgFile              string `xml:"ImgFile,attr"`
	Name                 string `xml:"Name,attr"`
	Author               string `xml:"Author,attr"`
	ExtraLangFilePrefix  string `xml:"ExtraLangFilePrefix,attr"`
	DefaultExtraLangFile string `xml:"DefaultExtraLangFile,attr"`
}

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
