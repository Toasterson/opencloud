package installd

import (
	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/opencloud/zpool"
	"github.com/toasterson/opencloud/zfs"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/toasterson/opencloud/common"
	"github.com/toasterson/opencloud/mount"
	"os"
	"path/filepath"
	"strings"
	"io"
	"syscall"
	"time"
)

const altRootLocation string = "/a"
const altMountLocation string = "/mnt.install"

func Install(conf InstallConfiguration) {
	//TODO Switch to switch statement
	if conf.InstallType == "fulldisk"{
		util.Must(formatDrives(&conf))
	}
	util.Must(CreateAndMountZpool(&conf))
	CreateDatasets(&conf)
	util.Must(InstallOS(&conf))
}

func fixZfsMountPoints(conf *InstallConfiguration) {
	bootenv := zfs.OpenDataset(fmt.Sprintf("%s/ROOT/%s", conf.RPoolName, conf.BEName))
	bootenv.SetProperty("canmount", "noauto")
	bootenv.SetProperty("mountpoint", "/")
}

func formatDrives(conf *InstallConfiguration) (err error) {
	//TODO sanity checks if drive exists
	//TODO Formating
	return
}

func CreateAndMountZpool(conf *InstallConfiguration) (err error){
	if conf.BEName == "" {
		conf.BEName = "openindiana"
	}
	rpool, err := zpool.CreatePool(conf.RPoolName, conf.PoolArgs, true, conf.PoolType, conf.Disks, false)
	if err != nil {
		return
	}
	err = rpool.SetProperty("bootfs", fmt.Sprintf("%s/ROOT/%s", conf.RPoolName, conf.BEName))
	return
}

func InstallOS(conf *InstallConfiguration) (err error) {
	switch conf.MediaType {
	case MediaTypeSolNetBoot:
		//Get the files Needed to /tmp
		getMediaFiles(conf)
		installOSFromMediaFiles("/tmp")
	case MediaTypeSolCDrom:
	case MediaTypeSolUSB:
		//Assume everything needed is located under /.cdrom
		installOSFromMediaFiles("/.cdrom")
	case MediaTypeIPS:
		return common.NotSupportedError("IPS installation")
	case MediaTypeZAP:
		return common.NotSupportedError("ZAP installation")
	case MediaTypeZImage:
		return common.NotSupportedError("Image installation")
	default:
		return common.InvalidConfiguration("MediaType")
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
	var err error
	if conf.InstallType != "bootenv"{
		_, err = zfs.CreateDataset(fmt.Sprintf("%s/ROOT", conf.RPoolName), zfs.DatasetTypeFilesystem, map[string]string{"mountpoint":"legacy"},true)
		util.Must(err)
		_, err = zfs.CreateDataset(fmt.Sprintf("%s/swap", conf.RPoolName), zfs.DatasetTypeVolume, map[string]string{"blocksize":"4k", "size": conf.SwapSize},true)
		util.Must(err)
		_, err = zfs.CreateDataset(fmt.Sprintf("%s/dump", conf.RPoolName), zfs.DatasetTypeVolume, map[string]string{"size": conf.DumpSize},true)
		util.Must(err)
		//TODO Zfs Layout Creation
	}
	if conf.MediaType != MediaTypeZImage {
		bootenv, err := zfs.CreateDataset(fmt.Sprintf("%s/ROOT/%s", conf.RPoolName, conf.BEName), zfs.DatasetTypeFilesystem, map[string]string{"mountpoint": altRootLocation},true)
		u1 := uuid.NewV4()
		bootenv.SetProperty("org.opensolaris.libbe:uuid", u1.String())
		util.Must(err)
	}
}

func getMediaFiles(conf *InstallConfiguration){
	util.Must(HTTPDownload(fmt.Sprintf("%s/solaris.zlib", conf.MediaURL), "/tmp"))
	//TODO different locations on i86 and amd64
	util.Must(HTTPDownload(fmt.Sprintf("%s/platform/i86pc/boot_archive", conf.MediaURL), "/tmp"))
}

func installOSFromMediaFiles(saveLocation string) {
	os.Mkdir(altMountLocation, 0755)
	util.Must(mount.MountLoopDevice("ufs", altMountLocation, fmt.Sprintf("%s/boot_archive", saveLocation)))
	util.Must(mount.MountLoopDevice("ufs", fmt.Sprintf("%s/usr", altMountLocation), fmt.Sprintf("%s/solaris.zlib", saveLocation)))
	filelist := []string{
		"boot",
		"kernel",
		"lib",
		"platform",
		"root",
		"sbin",
		"usr",
		"etc",
		"var",
		"opt",
		"zonelib",
	}
	for _, dir := range filelist{
		filepath.Walk(fmt.Sprintf("%s/%s", altMountLocation, dir), walkCopy)
	}
}

func walkCopy(path string, info os.FileInfo, err error) error {
	dstpath := strings.Replace(path, altMountLocation, altRootLocation, -1)
	if info.IsDir(){
		util.Must(os.Mkdir(dstpath, info.Mode()))
	} else {
		src, err := os.Open(path)
		defer src.Close()
		util.Must(err)
		dst, err := os.Create(dstpath)
		defer dst.Close()
		util.Must(err)
		_, err = io.Copy(dst, src)
		util.Must(err)
		util.Must(dst.Sync())
		srcStat := info.Sys().(*syscall.Stat_t)
		util.Must(syscall.Chmod(dstpath, srcStat.Mode))
		util.Must(syscall.Chown(dstpath, int(srcStat.Uid), int(srcStat.Gid)))
		util.Must(os.Chtimes(dstpath, time.Unix(int64(srcStat.Atim.Sec),int64(srcStat.Atim.Nsec)), time.Unix(int64(srcStat.Mtim.Sec),int64(srcStat.Mtim.Nsec))))
	}
	return nil
}