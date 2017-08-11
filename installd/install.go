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
	"github.com/toasterson/mozaik/logger"
)

const altRootLocation string = "/a"
const altMountLocation string = "/mnt.install"
const solusrfileName string = "solaris.zlib"
const solmediarootfileName string = "boot_archive"

func Install(conf InstallConfiguration) {
	//TODO Switch to switch statement
	if conf.InstallType == "fulldisk" {
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

func CreateAndMountZpool(conf *InstallConfiguration) (err error) {
	_, err = zpool.CreatePool(conf.RPoolName, conf.PoolArgs, true, conf.PoolType, conf.Disks, true)
	if err != nil {
		return
	}
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

func getMediaFiles(conf *InstallConfiguration) {
	util.Must(HTTPDownload(fmt.Sprintf("%s/%s", conf.MediaURL, solusrfileName), "/tmp"))
	//TODO different locations on i86 and amd64
	util.Must(HTTPDownload(fmt.Sprintf("%s/platform/i86pc/%s", conf.MediaURL, solmediarootfileName), "/tmp"))
}

func installOSFromMediaFiles(saveLocation string) {
	os.Mkdir(altMountLocation, os.ModeDir)
	util.Must(mount.MountLoopDevice("ufs", altMountLocation, fmt.Sprintf("%s/%s", saveLocation, solmediarootfileName)))
	util.Must(mount.MountLoopDevice("hsfs", fmt.Sprintf("%s/usr", altMountLocation), fmt.Sprintf("%s/%s", saveLocation, solusrfileName)))
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
	for _, dir := range filelist {
		filepath.Walk(fmt.Sprintf("%s/%s", altMountLocation, dir), walkCopy)
	}
}

func walkCopy(path string, info os.FileInfo, err error) error {
	dstpath := strings.Replace(path, altMountLocation, altRootLocation, -1)
	lsrcinfo, err := os.Lstat(path)
	if os.IsNotExist(err) {
		//Ignore Unexistent Directories
		return nil
	}
	util.Must(err)
	if info.IsDir() {
		logger.Trace(fmt.Sprintf("Mkdir %s", dstpath))
		util.Must(os.Mkdir(dstpath, info.Mode()))
	} else if lsrcinfo.Mode()&os.ModeSymlink != 0 {
		//We have a Symlink thus Create it on the Target
		dstTarget, _ := os.Readlink(path)
		logger.Trace(fmt.Sprintf("Creating Symlink %s -> %s", dstpath, dstTarget))
		util.Must(os.Symlink(dstTarget, dstpath))
	} else {
		//We Have a regular File Copy it
		go copyFileExact(path, info, dstpath)
	}
	return nil
}

func copyFileExact(source string, srcInfo os.FileInfo, dest string) {
	//logger.Trace(fmt.Sprintf("Copy %s -> %s", path, dest))
	src, err := os.Open(source)
	defer src.Close()
	util.Must(err)
	dst, err := os.Create(dest)
	defer dst.Close()
	util.Must(err)
	_, err = io.Copy(dst, src)
	util.Must(err)
	//util.Must(dst.Sync())
	srcStat := srcInfo.Sys().(*syscall.Stat_t)
	util.Must(syscall.Chmod(dest, srcStat.Mode))
	util.Must(syscall.Chown(dest, int(srcStat.Uid), int(srcStat.Gid)))
	//util.Must(os.Chtimes(dest, time.Unix(int64(srcStat.Atim.Sec),int64(srcStat.Atim.Nsec)), time.Unix(int64(srcStat.Mtim.Sec),int64(srcStat.Mtim.Nsec))))
}

/*
func MakeSystemDirectories(conf *InstallConfiguration){
	dirList := []string{
		"tmp",
		"system",
		"system/contract",
		"system/object",
		"system/boot",
		"proc",
		"mnt",
		"dev",
		"devices",
		"devices/pseudo",
		"dev",
		"dev/fd",
		"dev/rmt",
		"dev/swap",
		"dev/dsk",
		"dev/rdsk",
		"dev/net",
		"dev/ipnet",
		"dev/sad",
		"dev/pts",
		"dev/term",
		"dev/vt",
		"dev/zcons",
	}

}
*/
