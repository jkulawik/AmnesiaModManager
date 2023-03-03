package main

import (
	"embed"
	"errors"
	"image"
	"log"
	"os"
	"strings"

	"fmt"
	"image/jpeg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ftrvxmtrx/tga"
)

const (
	appInfo        = "Amnesia Mod Manager v1.2\nCopyright 2023 - github.com/jkulawik/ a.k.a. Darkfire"
	helpDeleteInfo = "Saves tied to mods currently do not get deleted.\n" +
		"Custom stories can be deleted entirely.\n" +
		"Full conversions might leave leftovers because many of them\n" +
		"do not properly list all of their folders and files in their config."
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger

	csPath          string = "custom_stories"
	customStories   []*CustomStory
	fullConversions []*FullConversion
	selectedMod     Mod

	windowContent *fyne.Container

	//go:embed default.jpg
	defaultImgFS embed.FS
	//go:embed icon.png
	iconBytes []byte
)

func initLoggers() {
	WarningLogger = log.New(os.Stderr, "WARNING: ", log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
	InfoLogger = log.New(os.Stderr, "INFO: ", log.Lshortfile)
}

func main() {
	// os.Chdir("testdata") // Debug
	initLoggers()
	a := app.New()
	a.SetIcon(fyne.NewStaticResource("amm_icon", iconBytes))
	w := a.NewWindow("Amnesia Mod Manager")

	err := CheckIsRootDir(".")
	displayIfError(err, w)

	customStories, err = GetCustomStories(csPath)
	displayIfError(err, w)
	fullConversions, err = GetFullConversions()
	displayIfError(err, w)

	windowContent = container.NewMax()
	toolbar := makeToolbar(w, a)
	windowContent.Objects = []fyne.CanvasObject{makeModTypeTabs()}
	windowContent.Refresh()

	mainView := container.NewBorder(
		container.NewVBox(toolbar, widget.NewSeparator()),
		nil, nil, nil, windowContent)

	w.SetContent(mainView)

	w.Resize(fyne.NewSize(900, 480))
	w.ShowAndRun()
}

func makeToolbar(window fyne.Window, app fyne.App) fyne.CanvasObject {
	t := widget.NewToolbar(
		widget.NewToolbarAction(theme.InfoIcon(), func() { dialog.ShowInformation("About", appInfo, window) }),
		widget.NewToolbarAction(theme.SettingsIcon(), func() { showSettings(app) }),
		//widget.NewToolbarAction(theme.ConfirmIcon(), func() { fmt.Println("Mark") }),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { refreshMods(window) }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.QuestionIcon(), func() { dialog.ShowInformation("Help: Deleting mods", helpDeleteInfo, window) }),
		widget.NewToolbarAction(theme.DeleteIcon(), func() { deleteSelectedMod(window) }),
	)
	return t
}

func makeModTypeTabs() fyne.CanvasObject {

	csTabContent := container.NewMax()
	csTabContent.Objects = []fyne.CanvasObject{makeCustomStoryListTab()}
	csTabContent.Refresh()

	fcTabContent := container.NewMax()
	fcTabContent.Objects = []fyne.CanvasObject{makeFullConversionListTab()}
	fcTabContent.Refresh()

	return container.NewAppTabs(
		container.NewTabItem("Full Conversions", fcTabContent),
		container.NewTabItem("Custom Stories", csTabContent),
	)
}

// ----------------------- Custom Stories ----------------------- //

func makeCustomStoryListTab() fyne.CanvasObject {
	var data = customStories

	cardContentLabel := widget.NewLabel("")
	cardContentLabel.Wrapping = fyne.TextWrapWord

	defaultTitle := "Select a custom story"
	card := widget.NewCard(defaultTitle, "", cardContentLabel)

	defaultImgRaw, _ := defaultImgFS.Open("default.jpg")
	img, _ := jpeg.Decode(defaultImgRaw)
	var defaultImg = canvas.NewImageFromImage(img)
	//card.Image = defaultImg
	displayImg := defaultImg

	storyViewContainer := container.New(layout.NewMaxLayout(), displayImg, card)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id].name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedMod = data[id]
		card.SetTitle(data[id].name)
		card.SetSubTitle("Author: " + data[id].author)
		cardContentLabel.SetText(makeStoryText(data[id]))

		if data[id].imgFile == "" {
			// card.SetImage(defaultImg)
			displayImg = defaultImg
		} else {
			imgFile := data[id].dir + data[id].imgFile
			displayImg = canvas.NewImageFromFile(imgFile)
			// card.SetImage(displayImg)
		}
		storyViewContainer.Objects[0] = displayImg
	}
	list.OnUnselected = func(id widget.ListItemID) {
		selectedMod = nil
		card.SetTitle(defaultTitle)
		card.SetSubTitle("")
		cardContentLabel.SetText("")
	}
	// listTab := container.NewHSplit(list, container.New(layout.NewVBoxLayout(), card))
	listTab := container.NewHSplit(list, storyViewContainer)
	listTab.SetOffset(0.3)
	return listTab
}

func makeStoryText(cs *CustomStory) string {
	return fmt.Sprintf("Folder:\n%s\nDescription:\n%s", cs.dir, cs.desc)
}

// ----------------------- General ----------------------- //

func displayIfError(err error, w fyne.Window) {
	if err != nil {
		ErrorLogger.Println(err)
		dialog.ShowError(err, w)
	}
}

func showSettings(a fyne.App) {
	w := a.NewWindow("Theme Settings")
	w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
	w.Resize(fyne.NewSize(480, 480))
	w.Show()
}

func refreshMods(w fyne.Window) {
	var err error
	customStories, err = GetCustomStories(csPath)
	displayIfError(err, w)
	fullConversions, err = GetFullConversions()
	displayIfError(err, w)

	windowContent.Objects = []fyne.CanvasObject{makeModTypeTabs()}
	windowContent.Refresh()
}

func deleteSelectedMod(w fyne.Window) {
	if selectedMod == nil {
		displayIfError(errors.New("no mod selected"), w)
		return
	}

	folderList := formatStringList(selectedMod.listFolders())

	warningMessage := "Delete the following folders?\n\n" + folderList
	warningMessage += "\nAll files will be deleted permanently.\n\nMod saves will not be deleted."

	cnf := dialog.NewConfirm("Confirmation", warningMessage, confirmDeleteCallback, w)
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.Show()
}

func confirmDeleteCallback(response bool) {
	if response {
		for _, f := range selectedMod.listFolders() {
			err := os.RemoveAll(f)
			if err != nil {
				ErrorLogger.Println(err)
			}

			// TODO Refresh
		}
	}
}

// ----------------------- FC tab ----------------------- //

func makeFullConversionListTab() fyne.CanvasObject {
	var data = fullConversions

	cardContentLabel := widget.NewLabel("")
	cardContentLabel.Wrapping = fyne.TextWrapWord

	defaultTitle := "Select a full conversion"
	card := widget.NewCard(defaultTitle, "", cardContentLabel)

	defaultImgRaw, _ := defaultImgFS.Open("default.jpg")
	img, _ := jpeg.Decode(defaultImgRaw)
	var defaultImg = canvas.NewImageFromImage(img)
	card.Image = nil

	launchButton := widget.NewButton("Launch", launchFullConversion)
	launchButton.Hide()

	vbox := container.NewVBox(card, launchButton)
	hbox := container.NewHBox(vbox)

	// fcViewContainer := container.New(layout.NewMaxLayout(), defaultImg, card)
	fcViewContainer := container.New(layout.NewMaxLayout(), defaultImg, hbox)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.FileApplicationIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id].name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedMod = data[id]
		// card.SetSubTitle("Author: " + data[id].author)
		// cardContentLabel.SetText(makeStoryText(data[id]))
		// cardContentLabel.SetText("This is a very very long text which should wrap around. White Night is an amazing mod, don't play it")
		// folderString := formatStringList(data[id].uniqueResources)
		// cardContentLabel.SetText("Mod folder(s):\n" + folderString)
		launchButton.Show()

		// InfoLogger.Println("Logo for", data[id].name, "is", data[id].logo)

		if data[id].logo == "" {
			card.SetImage(nil)
			card.SetTitle(data[id].name)
		} else {
			card.SetTitle(data[id].name) // TODO we have the logo, no need to clutter the space further?
			// card.SetSubTitle(getStringSpacer(90)) // to not let the card shrink too much
			// card.SetSubTitle("")

			displayImg := getImageFromFile(data[id].logo)
			displayImg.FillMode = canvas.ImageFillOriginal
			card.SetImage(displayImg)
		}
	}
	list.OnUnselected = func(id widget.ListItemID) {
		selectedMod = nil
		card.SetTitle(defaultTitle)
		card.SetSubTitle("")
		cardContentLabel.SetText("")
		launchButton.Hide()
	}
	// listTab := container.NewHSplit(list, container.New(layout.NewVBoxLayout(), card))
	listTab := container.NewHSplit(list, fcViewContainer)
	listTab.SetOffset(0.3)
	return listTab
}

func loadTGA(path string) image.Image {
	imgRaw, err := os.Open(path)
	if err != nil {
		ErrorLogger.Println(err)
	}
	img, err := tga.Decode(imgRaw)
	if err != nil {
		ErrorLogger.Println(err)
	}
	return img
}

func getImageFromFile(path string) *canvas.Image {
	if strings.Contains(path, ".tga") {
		img := loadTGA(path)
		return canvas.NewImageFromImage(img)
	} else {
		return canvas.NewImageFromFile(path)
	}
}

func launchFullConversion() {
	InfoLogger.Println("Launch button pressed")
}
