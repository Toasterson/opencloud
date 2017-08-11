package zpool

import (
	"fmt"
	"strings"
	"github.com/toasterson/opencloud/zfs"
)

// SetProperty set ZFS pool property to value. Not all properties can be set,
// some can be set only at creation time and some are read only.
// Always check if returned error and its description.
func (p *Pool) SetProperty(prop string, value string) (err error) {
	if _, err = zpoolExec([]string{"set", fmt.Sprintf("%s=%s", prop, value), p.Name}); err != nil {
		return
	}
	p.Properties[prop], err = p.GetProperty(prop)
	return
}

// GetProperty reload and return single specified property. This also reloads requested
// property in Properties map.
func (p *Pool) GetProperty(propName string) (prop zfs.Property, err error) {
	var retVal []string
	if retVal, err = zpoolExec([]string{"get", propName, p.Name}); err != nil {
		return
	}
	//TODO switch to use -Hp as this does not print first line
	propLine := strings.Fields(retVal[0])
	prop.Value = propLine[1]
	prop.Source = propLine[2]
	return
}
