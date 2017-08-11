package imaged

import (
	"github.com/toasterson/opencloud/zfs"
)

const (
	// name of the service
	name        = "imaged"
	description = "OpenCloud Image Service"

	// port which daemon should be listen
	port = ":9999"
)

type Imaged struct {
}

func (this *Imaged) List(pool string, reply *[]string) (err error) {
	*reply, err = zfs.List(pool)
	return err
}

//	dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}
