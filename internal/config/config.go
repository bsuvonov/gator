package config

import (
	"encoding/json"
	"os"
	"fmt"
)

const configFileName = ".gatorconfig.json"


type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName *string `json:"current_user_name,omitempty"`
}


func ReadConfig () (Config, error) {
	configFile, err := os.Open(configFileName)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	var config Config

	decoder := json.NewDecoder(configFile)

	if err = decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}


func (c *Config) SetUser(username string) error {
	c.CurrentUserName = &username

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	os.WriteFile(configFileName, jsonData, 0644)

	fmt.Println("The user has been set.")

	return nil
}