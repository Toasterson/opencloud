package bootadm

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/toasterson/opencloud/uname"
)

const (
	BootLoaderTypeLoader = 0
	BootLoaderTypeGrub = 1
)

const loaderConfFile string = "/%s/boot/menu.lst"

const loaderBootConfig string = `title {.BEName}
bootfs {.RPoolName}/ROOT/{.BEName}`

const grubBootConfig string = `default 0
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

type loaderType int

type BootConfig struct {
	Type loaderType
	RPoolName string
	BEName string
	BootOptions []string //TODO Implement
}

func CreateBootConfigurationFiles(rootDir string, conf BootConfig) (err error){
	if rootDir == "" {
		rootDir = "/"
	}

	hplatform := uname.GetHardwarePlatform()
	config := loaderBootConfig
	confLocation := loaderConfFile
	if hplatform == uname.HardwarePlatformXen {
		config = xenBootConfig
		confLocation = grubConfFile
	} else if conf.Type == BootLoaderTypeGrub {
		config = grubBootConfig
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
	if err = os.Mkdir(fmt.Sprintf("/%s/boot", conf.RPoolName), os.ModeDir); err != nil {
		return
	}
	confFile, err := os.Create(fmt.Sprintf(confLocation, conf.RPoolName))
	if err != nil {
		return
	}
	_, err = confFile.Write(out.Bytes())
	return
}
