package main

import (
	"fmt"

	"path/filepath"
)

func main() {
	found, _ := filepath.Glob("/usr/share/locale/en/*/*")
	fmt.Printf("%v", found)
}