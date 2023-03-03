package main

import "os"

func init() {
	initLoggers()
	os.Chdir("testdata")
}
