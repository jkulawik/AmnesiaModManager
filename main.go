package main

import (
	"embed"
	"errors"
	"log"
	"os"

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
)

const appInfo = "Amnesia Mod Manager v1.1\nCopyright 2023 - github.com/jkulawik/ a.k.a. Darkfire"
const workdir = "workdir"

var csPath string
var customStories []*CustomStory
var selectedMod Mod

//go:embed default.jpg
var defaultImgFS embed.FS

//go:embed icon.png
var iconBytes []byte

var windowContent *fyne.Container

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func initLoggers() {
	WarningLogger = log.New(os.Stderr, "WARNING: ", log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Lshortfile)
}

func main() {

	initLoggers()
	a := app.New()
	a.SetIcon(fyne.NewStaticResource("amm_icon", iconBytes))
	w := a.NewWindow("Amnesia Mod Manager")

	err := CheckIsRootDir(workdir)
	displayIfError(err, w)

	if workdir == "" {
		csPath = "custom_stories"
	} else {
		csPath = workdir + "/custom_stories"
	}
	customStories, err = GetCustomStories(csPath)
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

func displayIfError(err error, w fyne.Window) {
	if err != nil {
		ErrorLogger.Println(err)
		dialog.ShowError(err, w)
	}
}

func makeModTypeTabs() fyne.CanvasObject {

	content := container.NewMax()
	content.Objects = []fyne.CanvasObject{makeCustomStoryListTab()}
	content.Refresh()

	return container.NewAppTabs(
		container.NewTabItem("Custom Stories", content),
		container.NewTabItem("Full Conversions", widget.NewLabel("Coming soon (maybe)")),
	)
}

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

func makeToolbar(window fyne.Window, app fyne.App) fyne.CanvasObject {
	t := widget.NewToolbar(
		widget.NewToolbarAction(theme.InfoIcon(), func() { dialog.ShowInformation("About", appInfo, window) }),
		widget.NewToolbarAction(theme.SettingsIcon(), func() { showSettings(app) }),
		//widget.NewToolbarAction(theme.ConfirmIcon(), func() { fmt.Println("Mark") }),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { refreshCustomStories(window) }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() { deleteSelectedMod(window) }),
	)
	return t
}

func showSettings(a fyne.App) {
	w := a.NewWindow("Theme Settings")
	w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
	w.Resize(fyne.NewSize(480, 480))
	w.Show()
}

func refreshCustomStories(w fyne.Window) {
	var err error
	customStories, err = GetCustomStories(csPath)
	displayIfError(err, w)

	windowContent.Objects = []fyne.CanvasObject{makeModTypeTabs()}
	windowContent.Refresh()
}

func deleteSelectedMod(w fyne.Window) {
	if selectedMod == nil {
		displayIfError(errors.New("no mod selected"), w)
		return
	}

	folderList := ""
	for _, f := range selectedMod.listFolders() {
		folderList += f + "\n"
	}

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
