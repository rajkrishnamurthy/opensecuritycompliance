// This file is autogenerated. Modify as per your task needs.
package main

type UserInputs struct {
	OutputFileName       string `yaml:"OutputFileName"`
	ConfigFile           string `yaml:"ConfigFile"`
	RegoFile             string `yaml:"RegoFile"`
	Query                string `yaml:"Query"`
	OpaConfigurationFile string `yaml:"OpaConfigurationFile"`
	Source               string `yaml:"Source"`
	LogFile              string `yaml:"LogFile"`
	DataFile             string `yaml:"DataFile"`
}

type Outputs struct {
	OpaPolicyReport string
	LogFile         string `yaml:"AuditFile"`
}
