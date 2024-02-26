# AmnesiaModManager
A mod manager for "Amnesia: The Dark Descent".

![](/screenshots/screenshot_02.png)

Features:
* Browse all installed mods
  * Custom Stories
  * Full Conversions
  * Hybrid mods (full conversion custom stories introduced in Amnesia 1.5)
  * Steam Workshop mods
* Delete mods
  * Custom Stories can be deleted entirely (including hybrid mods)
  * Deleting Full Conversions is supported, but how well it works depends on the mod set-up.
  * Note: improperly configured mods (e.g. which ask you to dump files in the base game folders) can leave leftovers
* Launch Full Conversions from one place
  * Currently supported only for the Steam release. The app can only start the NoSteam version due to Steam's DRM.
  * Hybrid mods can also be started from the app
* Theme settings - light/dark, highlight color, font size

![](/screenshots/screenshot_01.png)
![](/screenshots/screenshot_03.png)

# Disclaimer

Windows Defender sometimes flags the application as malicious.

This is mostly due to the fact that Go is a language rarely used on desktops; you can scan the app yourself on [VirusTotal](https://www.virustotal.com/gui/home/upload) and you should see that 99% of the engines don't mark it as malicious.

# Usage

Get the program (modmanger.exe) for your system from the [releases](https://github.com/jkulawik/AmnesiaModManager/releases) section and place it in your Amnesia install folder (next to Amnesia.exe). 

# Changelog

* 1.2.6: first public release
* 1.3.0:
  * Added support for hybrid custom stories added in Amnesia 1.5 (i.e. full conversions that are launched from the game)
  * Add support for detecting custom stories downloaded from the Steam Workshop
  * Improved mod description formatting (previously the app displayed raw description contents)
  * Improved performance by adding an image cache
  * Improved lang file loading by correcting invalid XML comments (Amnesia's tinyXML doesn't mind `--` inside comments, but Go's standard library does)
  * Heavily restructured the app code to use proper Go conventions

# Known issues

On *some* Windows computers launching an FC stops the game registering mouse clicks.
