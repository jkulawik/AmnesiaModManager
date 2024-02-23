package main

import "os"

func init() {
	os.Chdir("testdata") // FIXME this was a lazy workaround, replace with proper fixture usage
}
