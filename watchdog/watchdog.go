package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	log "github.com/palette-software/insight-tester/common/logging"

	"gopkg.in/yaml.v2"
	"path/filepath"
)

type Webservice struct {
	Endpoint     string `yaml:"Endpoint"`
	UseProxy     bool   `yaml:"UseProxy"`
	ProxyAddress string `yaml:"ProxyAddress"`
}

type Config struct {
	Webservice Webservice `yaml:"Webservice"`
}

func setupUpdateServer() (string, error) {
	config, err := parseConfig()
	if err != nil {
		return "", err
	}

	// Do the proxy setup, if necessary
	err = setupProxy(config)
	if err != nil {
		return "", err
	}

	return config.Webservice.Endpoint, nil
}

func parseConfig() (Config, error) {
	var config Config

	configFilePath, err := findAgentConfigFile()
	if err != nil {
		return config, err
	}

	// Open agent's .yml config file
	input, err := os.Open(configFilePath)
	if err != nil {
		log.Error.Println("Error opening file: ", err)
		return config, err
	}
	defer input.Close()
	b, err := ioutil.ReadAll(input)
	if err != nil {
		log.Error.Println("Error reading file: ", err)
		return config, err
	}

	// Parse the .yml config file
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Error.Println("Error parsing yaml: ", err)
		return config, err
	}

	return config, nil
}

func setupProxy(config Config) error {
	// Set the proxy address, if there is any
	if config.Webservice.UseProxy {
		if len(config.Webservice.ProxyAddress) == 0 {
			err := fmt.Errorf("Missing proxy address from config file!")
			log.Error.Println(err)
			return err
		}
		proxyUrl, err := url.Parse(config.Webservice.ProxyAddress)
		if err != nil {
			log.Error.Printf("Could not parse proxy settings: %s from Config.yml. Error message: %s", config.Webservice.ProxyAddress, err)
			return err
		}
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
		log.Info.Println("Default Proxy URL is set to: ", proxyUrl)
	}

	return nil
}

// FIXME: locating the config file is not generic! This means this way is not going to be okay if we wanted to use this service as an auto-updater for the insight-server
// NOTE: This only works as long as the watchdog service runs from the very same folder as the agent.
// But they are supposed to be in the same folder by design.
func findAgentConfigFile() (string, error) {
	configPath := filepath.Join(baseFolder, "Config", "Config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Error.Println("Agent config file does not exist! Error message: ", err)
		return "", err
	}

	// Successfully located agent config file
	return configPath, nil
}
