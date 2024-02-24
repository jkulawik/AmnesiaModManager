package mods

import (
	"strings"
	"testing"
)

var testStoryMyMod = CustomStory{
	"Tutorial",
	"Mudbill",
	"new_story.lang",
	"testdata/custom_stories/MyMod",
	"Error while parsing lang file XML.",
	"customstory.png",
}

var testStoryEscape = CustomStory{
	"Another Madhouse mod",
	"Sabatu",
	"extra_english.lang",
	"testdata/custom_stories/_ESCAPE",
	"Another epic plot about people getting Amnesia",
	"yellow.jpg",
}

var testStoryBad = CustomStory{
	"Mod with no image",
	"Cmon",
	"extra_english.lang",
	"testdata/custom_stories/BadMod",
	"I don't know how you can miss the damn image.",
	"",
}

func TestReadCustomStoryConfig(t *testing.T) {
	csxml, err := ReadCustomStoryConfig("testdata/custom_stories/MyMod/custom_story_settings.cfg")

	if err != nil {
		t.Error(err)
		return
	}
	t.Log(*csxml)

	if csxml.Author != "Mudbill" {
		t.Errorf("wrong Author parameter: %s instead of Mudbill", csxml.Author)
	}
	if csxml.Name != "Tutorial" {
		t.Errorf("wrong Author parameter: %s instead of Tutorial", csxml.Name)
	}
	if csxml.ImgFile != "customstory.png" {
		t.Errorf("wrong Author parameter: %s instead of customstory.png", csxml.ImgFile)
	}

}

func TestGetDescFromLang(t *testing.T) {

	desc, err := GetDescFromLang("testdata/custom_stories/_ESCAPE/extra_english.lang")

	if err != nil {
		t.Error(err)
	}

	if desc != "Another epic plot about people getting Amnesia" {
		t.Errorf("wrong description: %s", desc)
	}
}

func TestGetStoryFromDir(t *testing.T) {
	cs, err := GetStoryFromDir("testdata/custom_stories/MyMod")

	if err != nil && !strings.Contains(err.Error(), "invalid sequence \"--\" not allowed in comments") {
		t.Error(err)
	}

	if cs == nil {
		t.Errorf("Received nil mod pointer")
	} else if *cs != testStoryMyMod {
		t.Errorf("Custom story did not match. Mock:\n%s\nRead:\n%s", cs, testStoryMyMod)
	}
}

func TestGetStoryFromDir2(t *testing.T) {
	cs, err := GetStoryFromDir("testdata/custom_stories/_ESCAPE")

	if err != nil {
		t.Error(err)
	}

	if cs == nil {
		t.Errorf("Received nil mod pointer")
	} else if *cs != testStoryEscape {
		t.Errorf("Custom story did not match. Mock:\n%s\nRead:\n%s", cs, testStoryEscape)
	}
}

func TestGetStoryNoImg(t *testing.T) {
	cs, err := GetStoryFromDir("testdata/custom_stories/BadMod")

	if err != nil {
		t.Error(err)
	}

	if cs == nil {
		t.Errorf("Received nil mod pointer")
	} else if *cs != testStoryBad {
		t.Errorf("Custom story did not match. Mock:\n%s\nRead:\n%s", cs, testStoryBad)
	}
}

func TestGetCustomStories(t *testing.T) {

	storyList, err := GetCustomStories("testdata/custom_stories")

	if err != nil {
		t.Error(err)
	}

	MyModFound := false
	EscapeModFound := false
	for _, cs := range storyList {
		t.Log(*cs)
		if *cs != testStoryEscape {
			EscapeModFound = true
		}
		if *cs != testStoryMyMod {
			MyModFound = true
		}
	}

	if !MyModFound {
		t.Error("did not find MyMod")
	}
	if !EscapeModFound {
		t.Error("did not find _ESCAPE")
	}
}
