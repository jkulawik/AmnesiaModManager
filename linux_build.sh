go build -o bin/modmanager.bin -ldflags "-s -w" . &&  zip ./bin/modmanager_Linux.zip ./bin/modmanager.bin
