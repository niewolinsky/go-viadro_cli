package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/spf13/viper"
)

func getUserHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Can't get your home directory.")
		os.Exit(1)
	}

	return usr.HomeDir
}

func getConfigDir() string {
	return path.Join(getUserHomeDir(), ".config")
}

func getConfigPath() string {
	return path.Join(getConfigDir(), "viadro.json")
}

func Init() {
	if _, err := os.Stat(getConfigDir()); os.IsNotExist(err) {
		err = os.Mkdir(getConfigDir(), os.ModeDir|0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(getConfigPath()); os.IsNotExist(err) {
		err = os.WriteFile(getConfigPath(), []byte("{}"), 0600)
		if err != nil {
			panic(err)
		}
	}

	viper.SetConfigFile(getConfigPath())
	viper.SetDefault("endpoint", "http://localhost:4000/v1/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
