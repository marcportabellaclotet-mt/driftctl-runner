package config

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type DritfctlRun struct {
	Description        string    `yaml:"description"`
	Group              string    `yaml:"group"`
	Provider           string    `yaml:"provider"`
	AWSConfig          AWSConfig `yaml:"aws"`
	TFStateList        []string  `yaml:"tfStateList"`
	ScanFilter         string    `yaml:"scanFilter"`
	DatadogIntegration bool      `yaml:"datadogIntegration"`
	Result             int
	ReportHTML         []byte
	Summary            Summary
}

type AWSConfig struct {
	AWSAccountId    string `yaml:"awsAccountId"`
	TFStateMethod   string `yaml:"tfStateMethod"`   // Options : default|awsProfile|assumeRole
	InfraScanMethod string `yaml:"infraScanMethod"` // Options : default|awsProfile|assumeRole
	TFStateRole     string `yaml:"tfStateRole"`
	InfraScanRole   string `yaml:"infraScanRole"`
	TFStateProfile  string `yaml:"tfStateProfile"`
}

type Summary struct {
	Coverage       float64
	TotalResources float64
	TotalChanged   float64
	TotalUnmanaged float64
	TotalDeleted   float64
	TotalManaged   float64
}

var DritfctlRunMap = map[string]DritfctlRun{}

func ReadConfig() {
	yfile, err := ioutil.ReadFile("driftctlRunner.yaml")
	if err != nil {
		logrus.Fatal(err)
	}
	err = yaml.Unmarshal(yfile, DritfctlRunMap)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("driftctl config has been loaded")
}
