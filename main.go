package main

import (
	"fmt"

	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Before")
	err := filepath.Walk("/usr/share/locale/en", localFunc)
	fmt.Println("After")
	fmt.Printf("%s", err)
}

func localFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir(){
		fmt.Println(path)
	}
	return nil
}