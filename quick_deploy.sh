#!/bin/bash
destination_path="/home/jacob/.var/app/com.valvesoftware.Steam/.local/share/Steam/steamapps/common/Amnesia The Dark Descent/modmanager.bin"
go build -o bin/modmanager.bin modmanager/cmd/modmanager && cp bin/modmanager.bin "$destination_path"