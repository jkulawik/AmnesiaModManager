package mods

type Mod interface {
	// This interface is needed basically only to allow deletion of both CS and FC mods
	listFolders() []string
}

var _ Mod = (*CustomStory)(nil)    // Check if CS implements interface (at compile time)
var _ Mod = (*FullConversion)(nil) // Check if FC implements interface (at compile time)

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

/*
	About the folder string format:
	After a lengthy battle with the way that filepaths are processed by various packages,
	the FC format was settled with no slashes at the start and end
	(fs.WalkDir treats those as invalid).
*/

type FullConversion struct {
	name            string
	mainInitConfig  string
	logo            string
	uniqueResources []string
}

func (fc FullConversion) listFolders() []string {
	return fc.uniqueResources
}
