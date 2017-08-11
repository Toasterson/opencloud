package imaged

import (
	"net/rpc"
	"fmt"
)

func client() (string, error) {
	c, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return "Error:", err
	}
	var result []string
	ferr := c.Call("Imaged.List", "rpool", &result)
	if ferr != nil {
		fmt.Println(ferr)
	} else {
		fmt.Println("Got Result: ")
		fmt.Println(result)
	}
	return "Sucess:", nil
}
