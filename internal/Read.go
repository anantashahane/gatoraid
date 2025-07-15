package internal

import (
	"encoding/json"
	"os"

	"github.com/anantashahane/gatoraid/internal/config"
)

const configFileName = ".gatoraidconfig.json"

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path + "/" + configFileName, nil
}

func Read() (config config.Config, err error) {
	fileName, err := getConfigFilePath()

	if err != nil {
		return config, err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return config, err
	}
	data := make([]byte, 400)
	size, err := file.Read(data)
	if err != nil {
		return config, err
	}
	data = data[:size]
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func SetUser(configuration config.Config, username string) (err error) {
	configuration.CurrentUserName = username
	data, err := json.Marshal(configuration)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
