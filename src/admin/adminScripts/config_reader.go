package adminScripts

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Database struct {
	IP       string `yaml:"ip"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Base     string `yaml:"Base"`
}

type Confing struct {
	SetupComplet string   `yaml:"setup_complet"`
	Database     Database `yaml:"database"`
}

func SetupComleteRead() bool {
	var conf Confing
	yamlFile, err := os.ReadFile("admin/adminScripts/configs/config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}

	if conf.SetupComplet == "true" {
		return true
	} else {
		return false
	}
}

func SetupComleteWrite() {

	// Создаем структуру с новым значением
	conf := Confing{
		SetupComplet: "true",
	}

	// Сериализуем структуру в YAML
	yamlFile, err := yaml.Marshal(&conf)
	if err != nil {
		panic(err)
	}

	// Записываем данные в файл
	err = os.WriteFile("admin/adminScripts/configs/config.yml", yamlFile, 0644)
	if err != nil {
		panic(err)
	}
}

func DatabaseSettingsReader() (string, string, string, string) {
	var conf Confing
	yamlFile, err := os.ReadFile("admin/adminScripts/configs/config.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
	return conf.Database.IP, conf.Database.User, conf.Database.Password, conf.Database.Base
}

func DatabaseSettingsWriter(ip string, user string, pass string, db string) {
	conf := Confing{
		Database: Database{
			IP:       ip,
			User:     user,
			Password: pass,
			Base:     db,
		},
	}

	yamlFile, err := yaml.Marshal(&conf)
	if err != nil {
		panic(err)
	}

	// Записываем данные в файл
	err = os.WriteFile("admin/adminScripts/configs/config.yml", yamlFile, 0644)
	if err != nil {
		panic(err)
	}
}
