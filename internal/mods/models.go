package mods

import "fmt"

// TODO make all methods use pointer receivers

type Mod interface {
	// This interface is needed basically only to allow deletion of both CS and FC mods
	ListFolders() []string
}

var _ Mod = (*CustomStory)(nil)    // Check if CS implements interface (at compile time)
var _ Mod = (*FullConversion)(nil) // Check if FC implements interface (at compile time)

type CustomStory struct {
	Name     string
	Author   string
	LangFile string
	Dir      string
	Desc     string
	ImgFile  string

	InitCfgFile     string
	IsHybrid        bool
	IsSteamWorkshop bool
}

func (cs CustomStory) ListFolders() []string {
	return []string{cs.Dir}
}

func (cs CustomStory) GetStoryText() string {
	desc := fmt.Sprintf("Folder:\n%s\nDescription:\n%s", cs.Dir, cs.Desc)
	if cs.IsHybrid {
		desc += "\n\nThis custom story is a full conversion and can also be launched in-game."
	}
	return desc
}

/*
	About the folder string format:
	After a lengthy battle with the way that filepaths are processed by various packages,
	the FC format was settled with no slashes at the start and end
	(fs.WalkDir treats those as invalid).
*/

type FullConversion struct {
	Name            string
	MainInitConfig  string
	Logo            string
	LangFile        string
	UniqueResources []string
}

func (fc FullConversion) ListFolders() []string {
	return fc.UniqueResources
}
