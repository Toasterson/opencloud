package installd

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/toasterson/opencloud/zpool"
	"github.com/toasterson/opencloud/zfs"
	"github.com/toasterson/mozaik/util"
)

func CreateAndMountZpool(conf *InstallConfiguration) (err error) {
	_, err = zpool.CreatePool(conf.RPoolName, conf.PoolArgs, true, conf.PoolType, conf.Disks, true)
	if err != nil {
		return
	}
	return
}

func CreateDatasets(conf *InstallConfiguration) {
	if conf.SwapSize == "" {
		conf.SwapSize = "2g"
	}
	if conf.DumpSize == "" {
		conf.DumpSize = conf.SwapSize
	}
	if conf.BEName == "" {
		conf.BEName = "openindiana"
	}
	var err error
	if conf.InstallType != "bootenv" {
		_, err = zfs.CreateDataset(fmt.Sprintf("%s/ROOT", conf.RPoolName), zfs.DatasetTypeFilesystem, map[string]string{"mountpoint": "legacy"}, true)
		util.Must(err)
		_, err = zfs.CreateDataset(fmt.Sprintf("%s/swap", conf.RPoolName), zfs.DatasetTypeVolume, map[string]string{"blocksize": "4k", "size": conf.SwapSize}, true)
		util.Must(err)
		_, err = zfs.CreateDataset(fmt.Sprintf("%s/dump", conf.RPoolName), zfs.DatasetTypeVolume, map[string]string{"size": conf.DumpSize}, true)
		util.Must(err)
		//TODO Zfs Layout Creation
	}
	if conf.MediaType != MediaTypeZImage {
		bootenv, err := zfs.CreateDataset(fmt.Sprintf("%s/ROOT/%s", conf.RPoolName, conf.BEName), zfs.DatasetTypeFilesystem, map[string]string{"mountpoint": altRootLocation}, true)
		u1 := uuid.NewV4()
		bootenv.SetProperty("org.opensolaris.libbe:uuid", u1.String())
		util.Must(err)
	}
	rpool := zpool.OpenPool(conf.RPoolName)
	err = rpool.SetProperty("bootfs", fmt.Sprintf("%s/ROOT/%s", conf.RPoolName, conf.BEName))
}