package main

import (
	"image"
	"modmanager/internal/logger"
	"os"
	"strings"

	// the _ means to import a package purely for its initialization side effects;
	// in this case png has to be registered, otherwise it causes errors about pngs being tga for some god forsaken reason
	_ "image/jpeg"
	_ "image/png"

	"fyne.io/fyne/v2/canvas"
	"github.com/ftrvxmtrx/tga"
)

var imageCache = make(map[string]*canvas.Image)

func loadTGA(path string) image.Image {
	// TODO This should return errors
	imgRaw, err := os.Open(path)
	if err != nil {
		logger.Error.Println(err)
	}
	img, err := tga.Decode(imgRaw)
	if err != nil {
		logger.Error.Println(err)
	}
	return img
}

func getImageFromFile(path string) *canvas.Image {
	// TODO this should handle and return errors
	img, ok := imageCache[path]
	if ok {
		return img
	}
	if strings.Contains(path, ".tga") {
		tga := loadTGA(path)
		img = canvas.NewImageFromImage(tga)
	} else {
		img = canvas.NewImageFromFile(path)
	}
	imageCache[path] = img
	return img
}
