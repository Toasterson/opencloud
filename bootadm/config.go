package bootadm

import (
	"github.com/toasterson/opencloud/uname"
	"github.com/toasterson/opencloud/installd"
	"text/template"
	"os"
	"fmt"
	"bytes"
)

const (
	BootLoaderTypeLoader = 0
	BootLoaderTypeGrub = 1
)

const loaderBootConfig string = `title {.BEName}
bootfs {.RPoolName}/ROOT/{.BEName}`

const loaderConfFile string = "/%s/boot/menu.lst"

const grubConfig string = `default 0
timeout 3
title {.BEName}
findroot (pool_{.RPoolName},0,a)
bootfs {.RPoolName}/ROOT/{.BEName}
kernel$ /platform/i86pc/kernel/$ISADIR/unix -B $ZFS-BOOTFS
module$ /platform/i86pc/$ISADIR/boot_archive`

const grubConfFile string = "/%s/boot/grub/menu.lst"

const xenBootConfig string = `default 0
timeout 3
title {.BEName}
findroot (pool_{.RPoolName},1,a)
bootfs {.RPoolName}/ROOT/{.BEName}
kernel$ /platform/i86pc/kernel/amd64/unix -B $ZFS-BOOTFS
module$ /platform/i86pc/amd64/boot_archive`


type bootLoaderType int

func CreateBootConfiguration(loaderType bootLoaderType, rootDir string, bootOptions []string, conf *installd.InstallConfiguration) (err error){
	if rootDir == "" {
		rootDir = "/"
	}
	hplatform := uname.GetHardwarePlatform()
	config := loaderBootConfig
	confLocation := loaderConfFile
	if hplatform == uname.HardwarePlatformXen || loaderType == BootLoaderTypeGrub {
		config = grubConfig
		confLocation = grubConfFile
	}
	tmplConfig, err := template.New("BootConfig").Parse(config)
	if err != nil {
		return
	}
	var out bytes.Buffer
	err = tmplConfig.Execute(&out, conf)
	if err != nil {
		return
	}
	confFile, err := os.Create(fmt.Sprintf(confLocation, conf.RPoolName))
	if err != nil {
		return
	}
	_, err = confFile.Write(out.Bytes())
	return
}
