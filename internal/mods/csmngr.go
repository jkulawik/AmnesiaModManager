package mods

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"modmanager/internal/configs"
	"modmanager/internal/logger"
)

func makeInvalidStory(dir string) *CustomStory {
	EmptyDirStory := new(CustomStory)

	EmptyDirStory.Name = "Invalid custom story"
	EmptyDirStory.Author = "N/A"
	EmptyDirStory.Dir = dir
	EmptyDirStory.Desc = "A folder was found which is not a valid custom story."

	return EmptyDirStory
}

func GetStoryFromDir(dir string) (*CustomStory, error) {
	csxml, err := configs.ReadCustomStoryConfig(dir + "/custom_story_settings.cfg")
	cs := new(CustomStory)

	if err != nil {
		if strings.Contains(err.Error(), "custom_story_settings.cfg: no such file or directory") {
			logger.Warn.Println(err)
			cs = makeInvalidStory(dir)
			return cs, nil
		}
		return nil, fmt.Errorf("GetStoryFromDir (%s): %w", dir, err)
	}

	cs.Dir = dir
	cs.Author = csxml.Author
	cs.Name = csxml.Name
	cs.ImgFile = strings.ReplaceAll(csxml.ImgFile, "\\", "/")
	cs.IsHybrid = csxml.InitCfgFile != ""
	if cs.IsHybrid {
		cs.InitCfgFile = strings.ReplaceAll(csxml.InitCfgFile, "\\", "/")
	}
	cs.IsSteamWorkshop = strings.Contains(dir, "workshop/content/57300")

	// Check if img file exists at specified path; search mod for it otherwise
	if _, err := os.Stat(cs.Dir + "/" + cs.ImgFile); err != nil {
		base_name := filepath.Base(cs.ImgFile)
		cs.ImgFile = ""
		logger.Info.Printf("GetStoryFromDir: searching for %s manually because an error occured: %v\n", base_name, err)
		walkFunc := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				logger.Warn.Println("GetStoryFromDir: fs.WalkDir walkFunc:", err)
				// Note: this return is consumed by fs.WalkDir so it won't show in logs or in app
				return fmt.Errorf("GetStoryFromDir: %w occured while crawling the filesystem", err)
			}
			if d.Name() == base_name {
				cs.ImgFile = path
				return fs.SkipAll
			}
			return nil
		}
		fs.WalkDir(os.DirFS(cs.Dir), ".", walkFunc)
	}

	if cs.IsHybrid {
		fc, err := GetConversionFromInit(dir, cs.Dir+"/"+cs.InitCfgFile)
		if err != nil {
			logger.Warn.Println("GetStoryFromDir:", err)
		} else {
			cs.LangFile = fc.LangFile
			cs.Logo = fc.Logo
		}
	} else {
		cs.LangFile = csxml.GetLangName()
	}
	cs.Desc, err = configs.GetDescFromLang(cs.Dir + "/" + cs.LangFile)

	if err != nil {
		if strings.Contains(err.Error(), "invalid sequence \"--\" not allowed in comments") {
			logger.Warn.Println(cs.Dir, err)
			return cs, nil
		}
		return cs, fmt.Errorf("GetStoryFromDir (%s): %w", dir, err)
	}

	return cs, nil
}

func GetCustomStories(dir string) ([]*CustomStory, error) {
	filelist, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("GetCustomStories: %w", err)
	}

	jobs := make(chan fs.DirEntry, len(filelist))     // bufferred so we can send all jobs without waiting for them to be consumed
	results := make(chan *CustomStory, len(filelist)) // bufferred should make things faster because workers don't have to wait to send results

	// Start 5 worker routines
	for i := 0; i < 1; i++ {
		go func() {
			// Consume jobs and send the results
			for dirEntry := range jobs {
				if !dirEntry.IsDir() {
					// We send nils so that the number of channel sends remains deterministic
					results <- nil
					continue
				}
				cs, err := GetStoryFromDir(dir + "/" + dirEntry.Name())

				if err != nil {
					logger.Error.Println("GetCustomStories: ", err)
					results <- nil
					continue
				}
				results <- cs
			}
		}()
	}

	// Start sending jobs into the pipeline
	// logger.Info.Println("Fill the CS job pipeline")
	for _, dirEntry := range filelist {
		logger.Info.Println("Sending job", dirEntry.Name())
		jobs <- dirEntry
	}
	close(jobs)

	// Start consuming results
	// logger.Info.Println("Start consuming the CS results pipeline")
	csList := make([]*CustomStory, 0, len(filelist))
	for i := 0; i < len(filelist); i++ {
		result := <-results
		logger.Info.Println("Consuming result", result)
		if result == nil {
			continue
		}
		csList = append(csList, result)
	}
	close(results) // kinda optional

	if len(csList) == 0 {
		logger.Info.Println("did not find any folders in the " + dir + " directory")
	}

	return csList, nil
}
