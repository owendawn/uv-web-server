package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Root string
	Port string
}

func SaveConfig(configPath string, config Config) {
	data, err := json.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(configPath, data, 0660)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadConfig(configPath string) (config *Config) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		return nil
	}
	config = &Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		println(err)
		return nil
	}
	return config
}
