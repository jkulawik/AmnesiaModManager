package configs

import "encoding/xml"

type ResourcesXML struct {
	XMLName   xml.Name                `xml:"Resources"`
	Directory []ResourcesXMLDirectory `xml:"Directory"`
}

type ResourcesXMLDirectory struct {
	Path       string `xml:"Path,attr"`
	AddSubDirs string `xml:"AddSubDirs,attr"`
}
