package crash

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Variables map[string]string `yaml:"vars,omitempty"`
	Inputs    *IOConfigs        `yaml:"inputs,omitempty"`
	Outputs   *IOConfigs        `yaml:"outputs,omitempty"`
	Checks    map[string]string `yaml:"checks,omitempty"`
	Plans     *PlanConfigs      `yaml:"plans"`
}

func NewConfig() *Config {
	config := Config{}
	return &config
}

func (config *Config) UnmarshalYAML(bytes []byte) error {
	err := yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) UnmarshalYAMLFile(file string) error {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return config.UnmarshalYAML(fileBytes)
}

func (config *Config) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(config)
}

func (config *Config) Dump() {
	d, err := config.MarshalYAML()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", string(d))
}

type IOConfig struct {
	Name string
	Type string
	Params map[string]interface{}
}

type IOConfigs []IOConfig

type PlanConfig struct {
	Name  string       `yaml:"plan"`
	Steps *StepConfigs `yaml:"steps"`
}

type PlanConfigs []PlanConfig

type StepConfig struct {
	Run      *ActionConfig `yaml:"run,omitempty"`
	Serial   *StepConfigs  `yaml:"serial,omitempty"`
	Parallel *StepConfigs  `yaml:"parallel,omitempty"`

	Success *StepConfigs `yaml:"success,omitempty"`
	Failure *StepConfigs `yaml:"failure,omitempty"`
	Always  *StepConfigs `yaml:"always,omitempty"`

	Checks  []string     `yaml:"check,omitempty"`
	Timeout string       `yaml:"timeout,omitempty"`
}

type StepConfigs []StepConfig

type ActionConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Params map[string]string `yaml:"params"`
}
