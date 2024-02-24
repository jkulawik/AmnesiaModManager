CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows go build -o bin/modmanager.exe -ldflags "-s -w -H=windowsgui" modmanager/cmd/modmanager &&  zip ./bin/modmanager_Windows.zip ./bin/modmanager.exe
