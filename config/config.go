package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigOption struct {
	PostgreDB string
}

type MailServerInfo struct {
	MailFrom       string
	SendGridApiKey string
}

var (
	config *ConfigOption
	mail   *MailServerInfo
)

//GetConfigOption returns the ConfigOption
func GetConfigOption() *ConfigOption {
	if config != nil {
		return config
	}
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return config
}

func GetMailOption() *MailServerInfo {
	if mail != nil {
		return mail
	}
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&mail)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return mail
}
