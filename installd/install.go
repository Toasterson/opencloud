package installd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/toasterson/mozaik/util"
	"github.com/toasterson/opencloud/bootadm"
	"github.com/toasterson/opencloud/common"
	"github.com/toasterson/opencloud/mount"
	"github.com/toasterson/opencloud/zfs"
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
	rootDir := altRootLocation
	if conf.InstallType == InstallTypeBootEnv {
		rootDir = GetPathOfBootEnv(conf.BEName)
	}
	MakeSystemDirectories(rootDir, []DirConfig{})
	MakeDeviceLinks(rootDir, []LinkConfig{})
	util.Must(CreateDeviceLinks(rootDir, []string{}))
	bconf := bootadm.BootConfig{Type: bootadm.BootLoaderTypeLoader, RPoolName: conf.RPoolName, BEName: conf.BEName, BootOptions: []string{}}
	util.Must(bootadm.CreateBootConfigurationFiles(rootDir, bconf))
	util.Must(bootadm.UpdateBootArchive(rootDir))
	util.Must(bootadm.InstallBootLoader(rootDir, conf.RPoolName))
	fixZfsMountPoints(&conf)
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
