package main

import (
	"embed"
	"errors"
	"image/jpeg"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"modmanager/internal/logger"
	"modmanager/internal/mods"
)

const (
	mainTitle      = "Amnesia Mod Manager"
	appInfo        = "Amnesia Mod Manager v1.3.0\nCopyright 2023 - github.com/jkulawik/ a.k.a. Darkfire"
	helpDeleteInfo = "Saves tied to mods currently do not get deleted.\n" +
		"Custom stories can be deleted entirely.\n" +
		"Full conversions can leave leftovers because many of them\n" +
		"do not properly list all of their folders and files in their config."
	deleteWorkshopItemInfo = "deleting Steam Workshop mods directly is disabled because Steam\n" +
		"will redownload the mod automatically.\n" +
		"Please unsubscribe from the Workshop items instead"
	isTestDataBuild = true
	csPath          = "custom_stories"
	workshopPath    = "../../workshop/content/57300"
)

var (
	customStories      []*mods.CustomStory
	fullConversions    []*mods.FullConversion
	selectedStory      *mods.CustomStory
	selectedConversion *mods.FullConversion
	selectedMod        mods.Mod

	defaultImg    *canvas.Image
	windowContent *fyne.Container
	mainWindow    fyne.Window

	//go:embed default.jpg
	defaultImgFS embed.FS
	//go:embed icon.png
	iconBytes []byte
)

func main() {
	defaultImgRaw, _ := defaultImgFS.Open("default.jpg")
	img, _ := jpeg.Decode(defaultImgRaw)
	defaultImg = canvas.NewImageFromImage(img)

	a := app.New()
	a.SetIcon(fyne.NewStaticResource("icon.png", iconBytes))
	mainWindow = a.NewWindow(mainTitle)

	var err error
	start_time := time.Now()
	customStories, err = mods.GetCustomStories(csPath)
	displayIfError(err, mainWindow)
	workshopStories, err := mods.GetCustomStories(workshopPath)
	customStories = append(customStories, workshopStories...)
	displayIfError(err, mainWindow)
	fullConversions, err = mods.GetFullConversions(".")
	displayIfError(err, mainWindow)
	logger.Info.Println("Mods loaded in ", time.Since(start_time))
	err = mods.CheckIsRootDir(".")
	displayIfError(err, mainWindow)

	windowContent = container.NewMax()
	toolbar := makeToolbar(mainWindow, a)
	windowContent.Objects = []fyne.CanvasObject{makeModTypeTabs()}
	windowContent.Refresh()

	mainView := container.NewBorder(
		container.NewVBox(toolbar, widget.NewSeparator()),
		nil, nil, nil, windowContent)

	mainWindow.SetContent(mainView)

	mainWindow.Resize(fyne.NewSize(900, 480))
	mainWindow.ShowAndRun()
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
	// csTabContent.Refresh()  // doesn't seem to be needed

	fcTabContent := container.NewMax()
	fcTabContent.Objects = []fyne.CanvasObject{makeFullConversionListTab()}
	// fcTabContent.Refresh()  // doesn't seem to be needed

	tabs := container.NewAppTabs(
		container.NewTabItem("Custom Stories", csTabContent),
		container.NewTabItem("Full Conversions", fcTabContent),
	)

	tabs.OnSelected = func(currentTab *container.TabItem) {
		if currentTab.Text == "Full Conversions" {
			selectedMod = selectedConversion
		} else if currentTab.Text == "Custom Stories" {
			selectedMod = selectedStory
		}
	}
	return tabs
}

// ----------------------- General ----------------------- //

func displayIfError(err error, w fyne.Window) {
	if err == nil {
		return
	}
	logger.Error.Println(err)
	if strings.Contains(err.Error(), "did not find any valid full conversions") ||
		strings.Contains(err.Error(), "workshop/content/57300 no such file or directory") {
		return
	}
	dialog.ShowError(err, w)
}

func showSettings(a fyne.App) {
	w := a.NewWindow("Theme Settings")
	w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
	w.Resize(fyne.NewSize(480, 480))
	w.Show()
}

func refreshMods(w fyne.Window) {
	selectedConversion = nil
	selectedStory = nil
	selectedMod = nil

	var err error
	customStories, err = mods.GetCustomStories(csPath)
	displayIfError(err, w)
	workshopStories, err := mods.GetCustomStories(workshopPath)
	customStories = append(customStories, workshopStories...)
	displayIfError(err, mainWindow)
	fullConversions, err = mods.GetFullConversions(".")
	displayIfError(err, w)

	windowContent.Objects = []fyne.CanvasObject{makeModTypeTabs()}
	windowContent.Refresh()
}

func deleteSelectedMod(w fyne.Window) {
	if mods.IsModNil(selectedMod) {
		displayIfError(errors.New("no mod selected"), w)
		return
	}

	cs, ok := selectedMod.(*mods.CustomStory)
	if ok && cs.IsSteamWorkshop {
		dialog.ShowError(errors.New(deleteWorkshopItemInfo), w)
		return
	}

	folderList := strings.Join(selectedMod.ListFolders(), "\n")

	warningMessage := "Delete the following folders?\n\n" + folderList + "\n"
	warningMessage += "\nAll files will be deleted permanently.\n\nMod saves will not be deleted."

	cnf := dialog.NewConfirm("Confirmation", warningMessage, confirmDeleteCallback, w)
	cnf.SetDismissText("No")
	cnf.SetConfirmText("Yes")
	cnf.Show()
}

func confirmDeleteCallback(response bool) {
	if response {
		for _, f := range selectedMod.ListFolders() {
			err := mods.DeleteModDir(f)
			displayIfError(err, mainWindow)
		}
		refreshMods(mainWindow)
	}
}

// ----------------------- Custom Stories ----------------------- //

func makeCustomStoryListTab() fyne.CanvasObject {
	var data = customStories

	cardContentLabel := widget.NewLabel("")
	cardContentLabel.Wrapping = fyne.TextWrapWord

	defaultTitle := "Select a custom story"
	card := widget.NewCard(defaultTitle, "", cardContentLabel)

	//card.Image = defaultImg
	displayImg := defaultImg

	launchButton := widget.NewButton("Launch", launchHybridCustomStory)
	launchButton.Hide()
	vbox := container.NewVBox(card, layout.NewSpacer(), launchButton)
	storyViewContainer := container.New(layout.NewMaxLayout(), displayImg, vbox)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.FolderIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// Note: this is an update function that runs in a loop
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id].Name)
			if data[id].IsHybrid {
				item.(*fyne.Container).Objects[0].(*widget.Icon).Resource = theme.FolderNewIcon()
				item.(*fyne.Container).Objects[0].(*widget.Icon).Refresh()
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedStory = data[id]
		selectedMod = selectedStory
		card.SetTitle(data[id].Name)
		card.SetSubTitle("Author: " + data[id].Author)
		description := data[id].GetStoryText()
		cardContentLabel.SetText(description)

		if data[id].ImgFile == "" {
			displayImg = defaultImg
		} else {
			imgFile := data[id].Dir + "/" + data[id].ImgFile
			displayImg = canvas.NewImageFromFile(imgFile)
		}
		storyViewContainer.Objects[0] = displayImg

		if data[id].Logo == "" {
			card.SetImage(nil)
		} else {
			logo := getImageFromFile(data[id].Logo)
			logo.FillMode = canvas.ImageFillContain
			card.SetImage(logo)
		}
		if data[id].IsHybrid {
			launchButton.Show()
		}
	}
	list.OnUnselected = func(id widget.ListItemID) {
		selectedStory = nil
		selectedMod = nil
		card.SetTitle(defaultTitle)
		card.SetSubTitle("")
		cardContentLabel.SetText("")
		launchButton.Hide()
	}
	listTab := container.NewHSplit(list, storyViewContainer)
	listTab.SetOffset(0.3)
	return listTab
}

// ----------------------- FC tab ----------------------- //

func makeFullConversionListTab() fyne.CanvasObject {
	var data = fullConversions

	cardContentLabel := widget.NewLabel("")
	cardContentLabel.Wrapping = fyne.TextWrapWord

	defaultTitle := "Select a full conversion"
	card := widget.NewCard(defaultTitle, "", cardContentLabel)

	card.Image = nil

	launchButton := widget.NewButton("Launch", launchFullConversion)
	launchButton.Hide()

	vbox := container.NewVBox(card, layout.NewSpacer(), launchButton)
	fcViewContainer := container.New(layout.NewMaxLayout(), defaultImg, vbox)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.FolderNewIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id].Name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedConversion = data[id]
		selectedMod = selectedConversion
		cardContentLabel.SetText("Mod config location: " + data[id].MainInitConfig)
		// cardContentLabel.SetText("This is a very very long text which should wrap around. White Night is an amazing mod, don't play it")
		// folderString := formatStringList(data[id].uniqueResources)
		// cardContentLabel.SetText("Mod folder(s):\n" + folderString)
		launchButton.Show()

		// logger.Info.Println(data[id].Name, "logo:", data[id].Logo)

		card.SetTitle(data[id].Name)
		// card.SetSubTitle("")
		if data[id].Logo == "" {
			card.SetImage(nil)
		} else {
			displayImg := getImageFromFile(data[id].Logo)
			displayImg.FillMode = canvas.ImageFillContain
			card.SetImage(displayImg)
		}
	}
	list.OnUnselected = func(id widget.ListItemID) {
		selectedMod = nil
		selectedConversion = nil
		card.SetTitle(defaultTitle)
		card.SetSubTitle("")
		cardContentLabel.SetText("")
		launchButton.Hide()
	}
	listTab := container.NewHSplit(list, fcViewContainer)
	listTab.SetOffset(0.3)
	return listTab
}

func launchHybridCustomStory() {
	launchModFromInit(selectedStory.Dir + "/" + selectedStory.InitCfgFile)
}

func launchFullConversion() {
	launchModFromInit(selectedConversion.MainInitConfig)
}

func launchModFromInit(init_file string) {
	var execMap = map[string]string{
		"windows": ".\\Amnesia_NoSteam.exe",
		"linux":   "./Amnesia_NOSTEAM.bin.x86_64",
	}
	gameExe := execMap[runtime.GOOS]

	cmd := exec.Command(gameExe, init_file)
	// mainWindow.Hide() // TODO try if this fixes the FC issue

	// bar := widget.NewProgressBarInfinite()
	// dialog.ShowCustom(mainTitle, "benu", bar, mainWindow)
	// the above  doesn't use the deprecated NewProgressInfinite, but doesn't pause app usage like desired
	d := dialog.NewProgressInfinite(mainTitle, "Game is running...", mainWindow)
	d.Show()
	err := cmd.Run()
	d.Hide()
	displayIfError(err, mainWindow)
}
