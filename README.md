# AmnesiaModManager
A simple mod manager for "Amnesia: The Dark Descent".

![](/screenshots/screenshot_01.png)

Features:
* Browse mods
  * Custom Stories: displays all the in-game data: title, author, description, background image, as well as the folder name that the CS is installed in
  * Full Conversions: displays the mod title and its logo (if one is used)
  * Hybrid mods: The app now supports the hybrid mod type added in Amnesia 1.5. These show up in the custom story tab, but show more info than regular custom stories
* Delete mods
  * Custom Stories can be deleted entirely (including hybrid mods)
  * Deleting Full Conversions is supported, but how well it works depends on the mod set-up. Properly configured mods get deleted correctly, but messy mods might leave leftovers
* Launch Full Conversions from one place
  * Currently supported only for the Steam release. The app can only start the NoSteam version due to Steam's DRM.
* Theme settings - light/dark, highlight color, font size

![](/screenshots/screenshot_02.png)

This program was made mostly as practice with test-driven development and a test of the Fyne GUI package,
but it can prove useful to people who play a lot of Amnesia mods.

# Disclaimer

Windows Defender sometimes flags the application as malicious.
This is mostly due to the fact that Go is a language rarely used on desktops; you can scan the app yourself on [VirusTotal](https://www.virustotal.com/gui/home/upload) and you should see that none of the engines mark it as malicious.

# Usage

Get the program (modmanger.exe) for your system from the releases section and place it in your Amnesia install folder (next to Amnesia.exe). 

# Changelog

* 1.2.6: first public release
* 1.3.0:
  * Added support for hybrid custom stories added in Amnesia 1.5 (i.e. full conversions that are launched from the game)
  * Add support for detecting custom stories downloaded from the Steam Workshop
  * Improved mod description formatting (previously the app displayed raw description contents)
  * Improved performance by adding an image cache
  * Heavily restructured the app code to use proper Go conventions

# Known issues

On *some* Windows computers launching an FC stops the game registering mouse clicks.
