package flags

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Env struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}

type Configs struct {
	Envs []Env `yaml:"env"`
}

var (
	used bool = false
	currentFileDir string
)

func init() {
	if currentFileDir == ""{
		_, filename, _, _ := runtime.Caller(0)
		currentFileDir = filepath.Dir(filename)
	}
}

func InitServiceFlags() {
	if used {
		fmt.Println("flags init has been used")
		return
	}
	var yamlResult Configs
	yamlData, err := os.ReadFile(currentFileDir + "/values_local.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlData, &yamlResult)
	if err != nil {
		panic(err)
	}

	for _, v := range yamlResult.Envs {
		os.Setenv(v.Name, v.Value)
	}
	used = true
}