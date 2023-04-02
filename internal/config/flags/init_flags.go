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

// название файла, если "ENV_FILE_NAME" c секретами и прочими ENV переменными пустой
const defaultEnvFileName = "values_local.yaml"

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
	// считываем название файла, которое находится в пакете flags, откуда мы достаем секреты и другие env переменные
	// недопустимо использовать другое расположение файла
	envFileName := os.Getenv("ENV_FILE_NAME")
	if envFileName == "" {
		envFileName = defaultEnvFileName
	}
	if used {
		fmt.Println("flags init has been used")
		return
	}
	var yamlResult Configs
	yamlData, err := os.ReadFile(currentFileDir + "/" + envFileName)
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