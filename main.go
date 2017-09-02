package main

import (
	"fmt"

	"github.com/toasterson/uxfiletool"
)

func main() {
	err := uxfiletool.ExactCopy("/usr/share/locale/en/LC_MESSAGES/dino.mo", "/tmp/installd")
	fmt.Println(err)
}