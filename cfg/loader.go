package cfg

import (
	"os"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"fmt"
)

var DEFUALT_TOML_CONFIG string = 
`[config]
aliased = true
`

// TODO:
// - make the $HOME/.nate-editor folder if it doesn't exist
// - make the $HOME/.nate-editor/config.toml file if it doesn't exist
// - write a default toml file
// 

const (
	CONFIG_DIR_PATH = "/.nate-editor/"
	CONFIG_TOML_FILE = "config.toml"
)

func Setup() TomlConfig {
	CONFIG_DIR := os.Getenv("HOME") + CONFIG_DIR_PATH
	CONFIG_PATH := CONFIG_DIR + CONFIG_TOML_FILE

	// if the user doesn't have a /.nate-editor
	// directory we create it for them.
	if _, err := os.Stat(CONFIG_DIR); os.IsNotExist(err) {
		if err := os.Mkdir(CONFIG_DIR, 0775); err != nil {
			panic(err)
		}
	}

	// make sure a config.toml file exists in the
	// nate-editor directory.
	if _, err := os.Stat(CONFIG_PATH); os.IsNotExist(err) {
		configFile, fileCreateErr := os.Create(CONFIG_PATH)
		if fileCreateErr != nil {
			panic(fileCreateErr)
		}
		defer configFile.Close()

		// write some stuff
		_, writeErr := configFile.Write([]byte("[config]"))
		if writeErr != nil {
			panic(writeErr)
		}
		configFile.Sync()
	}

	if _, err := os.Open(CONFIG_PATH); err != nil {
		panic(err)
	}

	configTomlData, err := ioutil.ReadFile(CONFIG_PATH)
    if err != nil {
        panic(err)
    }

	var conf TomlConfig
	if _, err := toml.Decode(string(configTomlData), &conf); err != nil {
		panic(err)
	}

	fmt.Println("Loaded '" + CONFIG_PATH + "'.")
	fmt.Println(conf)
	return conf
}
