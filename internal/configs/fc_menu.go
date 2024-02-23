package configs

type MenuXML struct {
	Main MenuXMLMain `xml:"Main"`
}

type MenuXMLMain struct {
	MainFadeInTime      string `xml:"MainFadeInTime,attr"`
	MainFadeOutTimeFast string `xml:"MainFadeOutTimeFast,attr"`
	MainFadeOutTimeSlow string `xml:"MainFadeOutTimeSlow,attr"`

	TopMenuFadeInTime             string `xml:"TopMenuFadeInTime,attr"`
	TopMenuFadeOutTime            string `xml:"TopMenuFadeOutTime,attr"`
	TopMenuFontRelativeSize       string `xml:"TopMenuFontRelativeSize,attr"`
	TopMenuStartRelativePos       string `xml:"TopMenuStartRelativePos,attr"`
	TopMenuStartRelativePosInGame string `xml:"TopMenuStartRelativePosInGame,attr"`
	TopMenuFont                   string `xml:"TopMenuFont,attr"`

	MainMenuLogoStartRelativePos string `xml:"MainMenuLogoStartRelativePos,attr"`
	MainMenuLogoRelativeSize     string `xml:"MainMenuLogoRelativeSize,attr"`

	BGScene            string `xml:"BGScene,attr"`
	BGCamera_FOV       string `xml:"BGCamera_FOV,attr"`
	BGCamera_ZoomedFOV string `xml:"BGCamera_ZoomedFOV,attr"`

	ZoomSound string `xml:"ZoomSound,attr"`
	Music     string `xml:"Music,attr"`

	MenuLogo string `xml:"MenuLogo,attr"`
}
