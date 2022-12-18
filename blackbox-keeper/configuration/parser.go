package configuration

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func NewYamlConfiguration(filename string) (ConfigYaml, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ConfigYaml{}, err
	}
	defer file.Close()

	var cfg ConfigYaml
	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(false)
	err = decoder.Decode(&cfg)
	if err != nil {
		return ConfigYaml{}, err
	}

	return cfg, nil
}

func NewXmlConfiguration(filename string) (ConfigXml, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error while opening xml file:", err)
		return ConfigXml{}, err
	}
	defer file.Close()

	var cfg ConfigXml
	byteArray, err := io.ReadAll(file)
	err = xml.Unmarshal(byteArray, &cfg)

	if err != nil {
		fmt.Println("Error while unmarshalling xml file:", err)
		return ConfigXml{}, err
	}

	return cfg, nil
}
