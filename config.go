package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Title   string         `yaml:"title" default:"goemon - Traffic Viewer"`
	Map     string         `yaml:"map" default:""`
	Devices []DeviceConfig `yaml:"devices"`
}

type DeviceConfig struct {
	Name       string                  `yaml:"name" default:""`
	IpAddress  string                  `yaml:"ip_address"`
	Port       uint16                  `yaml:"port" default:"161"`
	Community  string                  `yaml:"community" default:"public"`
	Interfaces []DeviceInterfaceConfig `yaml:"interfaces"`
}

type DeviceInterfaceConfig struct {
	Name      string `yaml:"name"`
	PositionX int    `yaml:"position_x" default:"0"`
	PositionY int    `yaml:"position_y" default:"0"`
}

func LoadConfigs(fileName string, config *Config) error {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}
