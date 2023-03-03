package main

import (
	"strings"
	"testing"
)

var TestStoryMyMod = CustomStory{
	"Tutorial",
	"Mudbill",
	"new_story.lang",
	"custom_stories/MyMod/",
	"Error while parsing lang file XML.",
	"customstory.png",
}

var TestStoryEscape = CustomStory{
	"Another Madhouse mod",
	"Sabatu",
	"extra_english.lang",
	"custom_stories/_ESCAPE/",
	"Another epic plot about people getting Amnesia",
	"yellow.jpg",
}

var TestStoryBad = CustomStory{
	"Mod with no image",
	"Cmon",
	"extra_english.lang",
	"custom_stories/BadMod/",
	"I don't know how you can miss the damn image.",
	"",
}

func TestReadCustomStoryConfig(t *testing.T) {

	csxml, err := ReadCustomStoryConfig("custom_stories/MyMod/custom_story_settings.cfg")
	t.Log(*csxml)

	if err != nil {
		t.Error(err)
	}

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

	desc, err := GetDescFromLang("custom_stories/_ESCAPE/extra_english.lang")

	if err != nil {
		t.Error(err)
	}

	if desc != "Another epic plot about people getting Amnesia" {
		t.Errorf("wrong description: %s", desc)
	}
}

func TestGetStoryFromDir(t *testing.T) {
	cs, err := GetStoryFromDir("custom_stories/MyMod/")

	if err != nil && !strings.Contains(err.Error(), "invalid sequence \"--\" not allowed in comments") {
		t.Error(err)
	}

	if *cs != TestStoryMyMod {
		t.Errorf("Custom story did not match. Mock:\n%s\nRead:\n%s", cs, TestStoryMyMod)
	}
}

func TestGetStoryFromDir2(t *testing.T) {
	cs, err := GetStoryFromDir("custom_stories/_ESCAPE/")

	if err != nil {
		t.Error(err)
	}

	if *cs != TestStoryEscape {
		t.Errorf("Custom story did not match. Mock:\n%s\nRead:\n%s", cs, TestStoryEscape)
	}
}

func TestGetStoryNoImg(t *testing.T) {
	cs, err := GetStoryFromDir("custom_stories/BadMod/")

	if err != nil {
		t.Error(err)
	}

	if *cs != TestStoryBad {
		t.Errorf("Custom story did not match. Mock:\n%s\nRead:\n%s", cs, TestStoryBad)
	}
}

func TestGetCustomStories(t *testing.T) {

	storyList, err := GetCustomStories("custom_stories")

	if err != nil {
		t.Error(err)
	}

	MyModFound := false
	EscapeModFound := false
	for _, cs := range storyList {
		t.Log(*cs)
		if *cs != TestStoryEscape {
			EscapeModFound = true
		}
		if *cs != TestStoryMyMod {
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
