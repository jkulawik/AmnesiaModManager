go build -o bin/modmanager.bin -ldflags "-s -w" modmanager/cmd/modmanager &&  zip ./bin/modmanager_Linux.zip ./bin/modmanager.bin
