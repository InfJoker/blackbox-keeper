package configuration

import (
	"fmt"
	"github.com/sbabiv/xml2map"
	"gopkg.in/yaml.v2"
	"os"
)

func NewConfiguration(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	decoder.SetStrict(false)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	fmt.Printf("YAML config:\t%v\n", cfg)
	return cfg, nil
}

func NewXmlConfiguration(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := xml2map.NewDecoder(file)
	mymap, err := decoder.Decode()

	if err != nil {
		return Config{}, err
	}
	fmt.Printf("XML config:\t%v\n", mymap)
	return cfg, nil
}
