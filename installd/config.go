package installd

const (
	MediaTypeSolNetBoot = "solnetboot"
	MediaTypeSolCDrom   = "solcdrom"
	MediaTypeSolUSB     = "solusb"
	MediaTypeIPS        = "ips"
	MediaTypeZAP        = "zap"
	MediaTypeZImage     = "zimage"
)

type InstallConfiguration struct {
	InstallType         string `json:"install_type"`         //Possible options are efi, bootenv, fulldisk, efifulldisk
	RootFSType          string `json:"root_fs_type"`         //ufs, zfs
	Disks               []string `json:"disks"`              //The Disks that shall be used
	MediaType           string `json:"media_type"`           //Valid Values are SolNetboot, SolCdrom, SolUSB, IPS, ZAP, ZImage
	UseBootEnvironments bool `json:"use_boot_environments:"` //Whether to use boot environments or not
	ZFSLayout           interface{} `json:"zfs_layout"`      //The Partition Layout for ZFS. e.g Where is /var /etc and others located
	PoolArgs            map[string]string `json:"pool_args"` //Enable things like compression etc.
	PoolType            string `json:"pool_type"`            //The Type of pool eg mirrored, raidz, single(default)
	RPoolName           string `json:"rpool_name"`           //Name of the root pool
	MediaURL            string `json:"media_url"`            //The URL the media can be found at. uses install_media if ommitted
	BEName              string `json:"be_name"`              //Name of the new Boot Environment defaults to openindiana
	SwapSize            string `json:"swap_size"`            //Size of the SWAP Partition defaults to 2g
	DumpSize            string `json:"dump_size"`            //Size of the Dump Partition defaults to swap_size
}
