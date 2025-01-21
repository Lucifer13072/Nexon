package scripts

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Confing struct {
	SetupComplet bool `yaml:"setupComplet"`
}

func setup_complete_read() bool {
	var conf Confing
	yamlFile, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	return conf.SetupComplet
}

func setup_complete_write(marp bool) {
	conf := Confing{
		SetupComplet: marp,
	}

	yamlFile, err := yaml.Marshal(&conf)
	if err != nil {
		panic(err)
	}

	// Записываем данные в файл
	err = os.WriteFile("configs/config.yaml", yamlFile, 0644)
	if err != nil {
		panic(err)
	}
}
