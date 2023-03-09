# AmnesiaModManager
A simple mod manager for "Amnesia: The Dark Descent".

Features:
* Browse mods
  * Custom Stories tab displays all the in-game data: title, author, description, background image, as well as the folder name that the CS is installed in
  * Full Conversions tab displays the mod title and its logo (if one is used)
* Delete mods
  * Custom Stories can be deleted entirely
  * Deleting Full Conversions is supported, but how well it works depends on the mod set-up. Properly configured mods get deleted correctly, but messy mods might leave leftovers
* Launch Full Conversions from one place (this currently has a bug on some Windows machines which stops the game from registering mouse clicks)
* Theme settings - light/dark, highlight color, font size

This program was made mostly as practice with test-driven development and a test of the Fyne GUI package,
but it can prove useful to people who play a lot of Amnesia mods.

# Disclaimer

Windows Defender sometimes flags the application as malicious.
This is mostly due to the fact that Go is a language rarely used on desktops; you can scan the app yourself on [VirusTotal](https://www.virustotal.com/gui/home/upload) and you should see that none of the engines mark it as malicious.

# Usage

Get the program for your system from the releases section and place it in your Amnesia install folder (next to Amnesia.exe).

If you have Go installed, you can also compile it yourself by running `go build .` in the project folder.
