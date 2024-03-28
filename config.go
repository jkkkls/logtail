package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name      string `yaml:"name"`
	File      string `yaml:"file"`
	Separator string `yaml:"separator"`
	Fields    []struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"fields"`
	Out struct {
		Kafka struct {
			Hosts    []string `yaml:"hosts"`
			Topic    string   `yaml:"topic"`
			Sasl     string   `yaml:"sasl"`
			Username string   `yaml:"username"`
			Password string   `yaml:"password"`
		} `yaml:"kafka"`
	} `yaml:"out"`
}

func LoadConfig(fileName string) (*Config, error) {
	temp := &Config{}
	//
	buff, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buff, temp)
	if err != nil {
		return nil, err
	}

	return temp, nil
}
