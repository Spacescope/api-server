package utils

import (
	"syscall"

	"github.com/jinzhu/configor"
)

var CNF TomlConfig

type TomlConfig struct {
	Title     string
	APIServer APIServer `toml:"block-explorer-api"`
}

type APIServer struct {
	Addr   string `toml:"listen"`
	DB     string `toml:"db"`
	BusiDB string `toml:"busi_db"`
}

func InitConfFile(file string, cf *TomlConfig) error {
	err := syscall.Access(file, syscall.O_RDONLY)
	if err != nil {
		return err
	}
	err = configor.Load(cf, file)
	if err != nil {
		return err
	}

	return nil
}
