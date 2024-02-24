package mods

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
}

func (cs CustomStory) ListFolders() []string {
	return []string{cs.Dir}
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
	UniqueResources []string
}

func (fc FullConversion) ListFolders() []string {
	return fc.UniqueResources
}
