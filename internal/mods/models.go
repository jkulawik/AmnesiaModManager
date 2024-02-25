package mods

import "strings"

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
	Logo            string
	IsHybrid        bool
	IsSteamWorkshop bool
}

func (cs CustomStory) ListFolders() []string {
	return []string{cs.Dir}
}

func (cs CustomStory) GetStoryText() string {
	// desc := fmt.Sprintf("Folder:\n%s\nDescription:\n%s", cs.Dir, cs.Desc)
	text := "Folder:\n" + cs.Dir
	if cs.IsHybrid {
		text += "\nThis custom story is a hybrid full conversion and can be also launched from the game."
	}
	if cs.IsSteamWorkshop {
		text += "\nThis custom story was downloaded from the Steam Workshop."
	}
	description := strings.Replace(cs.Desc, "\n", " ", -1) // -1 means no limit on how many replaced
	description = strings.Replace(description, "\t", " ", -1)
	description = strings.Replace(description, "  ", " ", -1)
	description = strings.ReplaceAll(description, "[br]", "\n")
	description = strings.TrimSpace(description)
	text += "\n\nDescription:\n" + description
	return text
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
