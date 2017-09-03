package main

import (
	"fmt"

	"github.com/toasterson/uxfiletool"
)

func main() {
	fmt.Println(uxfiletool.FindByGlob("libnss*"))
}