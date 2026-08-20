package version

import (
	"fmt"
	"os"
)

// Version is set by the linker flags in the Makefile.
var Version string

func PrintVersionAndExit() {
	fmt.Printf("%s\n", Version)
	os.Exit(0)
}
