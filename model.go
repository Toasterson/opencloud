package main

import (
	"net"
	"time"
)

type Host struct {
	Hostname   string
	Ip         net.IPAddr
	Domainname string
	HostGroup  string
}

type OperatingSystem struct {
	Name string `json:"name"`
}

type RealmType struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type Machine struct {
	Hostname   string
	Created    time.Time
	Ip         net.IPAddr
	Domainname string
	Os         OperatingSystem
	Realm      Realm
}

type Realm struct {
	Name     string
	Created  time.Time
	Creator  User
	Type     RealmType
	Users    []User
	Admins   []User
	state    int
	Location string
}

type User struct {
	Name string
}

type App struct {
	name       string
	Type       int
	Maintainer User
}

type Instance struct {
	Name  string
	Realm Realm
	App   App
}

const (
	Initializing = 0
	Working
	Ready
	Shutdown
)
