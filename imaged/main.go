package imaged

import (
	"github.com/takama/daemon"
	"os"
	"fmt"
	"github.com/toasterson/opencloud/common"
)

func main() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		common.Errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := common.Server{srv, port, name, description}
	status, err := common.RunService(service, client)
	if err != nil {
		common.Errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}