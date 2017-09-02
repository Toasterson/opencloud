package image

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

const default_profile_filename = "profile.json"

type Type string

func (t Type) Validate() error {
	switch t {
	case
		TypeChroot,
		TypeZfs,
		TypeUfs,
		TypeTar,
		TypeACI:
	default:
		return fmt.Errorf("Image Type not known use one of ZFS|UFS|Chroot|Tar")
	}
	return nil
}

type Profile struct {
	Type Type `json:"type"`
	FileSets []string `json:"file_sets"`
	Manifest schema.ImageManifest `json:"manifest"`
	Files []string `json:"files,omitempty"`
	Users []string `json:"users,omitempty"`
	Groups []string `json:"groups,omitempty"`
	Devices []string `json:"devices,omitempty"`
}

func LoadProfile(file string) (*Profile, error) {
	buff, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	p := Profile{}
	return &p, json.Unmarshal(buff, &p)
}

func (p Profile) ResolveFiles(config *Config){
	p.Files = config.GetFiles(p.FileSets)
}

func (p Profile) Save(dir string) error {
	buff, err := json.Marshal(p)
	if err == nil {
		filep := filepath.Join(dir, default_profile_filename)
		return ioutil.WriteFile(filep, buff, 0644)
	}
	return err
}

func NewProfile(name string) (p Profile, err error) {
	p = Profile{
		Type: TypeChroot,
		Manifest: schema.ImageManifest{ACKind: schema.ImageManifestKind, ACVersion: schema.AppContainerVersion},
	}
	if name != ""{
		ident, err := types.NewACIdentifier(name)
		if err != nil {
			return p, err
		}
		p.Manifest.Name = *ident
	}
	return
}