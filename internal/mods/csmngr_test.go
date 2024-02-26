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
	"Description of Custom Story.",
	"customstory.png",
	"",
	"",
	false,
	false,
}

var testStoryEscape = CustomStory{
	"Another Madhouse mod",
	"Sabatu",
	"extra_english.lang",
	"testdata/custom_stories/_ESCAPE",
	"Another epic plot about people getting Amnesia",
	"yellow.jpg",
	"",
	"",
	false,
	false,
}

var testStoryBad = CustomStory{
	"Mod with no image",
	"Cmon",
	"extra_english.lang",
	"testdata/custom_stories/BadMod",
	"I don't know how you can miss the damn image.",
	"",
	"",
	"",
	false,
	false,
}

func TestGetStoryFromDir(t *testing.T) {
	cs, err := GetStoryFromDir("testdata/custom_stories/MyMod")

	if err != nil && !strings.Contains(err.Error(), "invalid sequence \"--\" not allowed in comments") {
		t.Error(err)
	}

	if cs == nil {
		t.Errorf("Received nil mod pointer")
	} else if *cs != testStoryMyMod {
		t.Errorf("Custom story did not match. Mock:\n%v\nGot:\n%v", testStoryMyMod, cs)
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
		t.Errorf("Custom story did not match. Mock:\n%v\nGot:\n%v", testStoryEscape, cs)
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
		t.Errorf("Custom story did not match. Mock:\n%v\nGot:\n%v", testStoryBad, cs)
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
