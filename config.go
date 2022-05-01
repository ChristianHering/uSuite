package main

import (
	"time"

	goessentials "github.com/ChristianHering/GoEssentials"
)

type Configuration struct {
	ListenAddr string
	DataDir    string
	SessionTTL time.Duration
	HTTPS      TLS
	FSCheck    bool
	DevMode    bool //Bypasses user && FS checks
}

type TLS struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

var configuration = Configuration{ //This is the default configuration
	ListenAddr: ":2119",
	DataDir:    "./data",
	SessionTTL: 2628000000000000, //1 month
	HTTPS: TLS{
		Enabled:  false,
		CertFile: "",
		KeyFile:  "",
	},
	FSCheck: true,
	DevMode: false,
}

func init() {
	err := goessentials.GetConfig("config.json", &configuration)
	if err != nil {
		panic(err)
	}
}
