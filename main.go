package main

import (
	"github.com/toasterson/opencloud/zfs"
	"log"
)

func main() {
	datasets , err:= zfs.List("rpool")
	if err != nil{
		log.Fatal(err)
	}
	log.Println(datasets)
	for _, dataset := range(datasets){
		size, serr := zfs.UsedIncludingChildren(dataset)
		if serr != nil {
			log.Fatal(serr)
		}
		log.Println(size)
	}
	log.Println("Test")
}
