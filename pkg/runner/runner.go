package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/aws"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/config"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/datadog"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/helpers"

	"github.com/sirupsen/logrus"
)

var defaultDelay int = 300

type scan struct {
	Coverage float64 `json:"coverage"`
	Summary  struct {
		TotalChanged   float64 `json:"total_changed"`
		TotalManaged   float64 `json:"total_managed"`
		TotalMissing   float64 `json:"total_missing"`
		TotalResources float64 `json:"total_resources"`
		TotalUnmanaged float64 `json:"total_unmanaged"`
	} `json:"summary"`
}

func parseJsonReport(jsonReport []byte) (summary config.Summary) {
	var scanRun scan
	err := json.Unmarshal(jsonReport, &scanRun)
	if err != nil {
		logrus.Error("Error when unmarshalling driftctl json report")
		return summary
	}

	summary.Coverage = scanRun.Coverage
	summary.TotalManaged = scanRun.Summary.TotalManaged
	summary.TotalUnmanaged = scanRun.Summary.TotalUnmanaged
	summary.TotalResources = scanRun.Summary.TotalResources
	summary.TotalChanged = scanRun.Summary.TotalChanged
	summary.TotalDeleted = scanRun.Summary.TotalMissing

	return summary
}

func Start() (err error) {
	logrus.Info("Starting driftctl-runner analysis daemon")
	config.ReadConfig()
	delay := helpers.LookupEnvOrDefaultInt(os.Getenv("DCTLRUNNER_DELAY_BETWEEN_SCANS"), defaultDelay)

	for {
		for k, v := range config.DritfctlRunMap {
			result, htmlReport, JsonReport := driftctlRun(k, v)
			if entry, ok := config.DritfctlRunMap[k]; ok {
				entry.Result = result
				if result == 1 || result == 0 {
					summary := parseJsonReport(JsonReport)
					entry.ReportHTML = htmlReport
					entry.Summary = summary
					if v.DatadogIntegration {
						datadog.SendMetrics(k)
					}
				}
				config.DritfctlRunMap[k] = entry
			}
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func newDriftctlScan(item string, config config.DritfctlRun) (exit int) {
	logrus.Infof("Starting a new driftctl scan for %s", item)
	cmdArgs := helpers.GenerateCMDArgs(item, config)
	logrus.Info(cmdArgs)
	stderr := new(strings.Builder)
	cmd := exec.Command("driftctl", cmdArgs...)
	cmd.Stderr = stderr
	cmd.Run()
	exit = cmd.ProcessState.ExitCode()
	if exit == 2 || exit == -1 {
		logrus.Println(stderr.String())
	}
	return exit
}

func driftctlRun(item string, config config.DritfctlRun) (exitCode int, htmlReport []byte, jsonReport []byte) {

	if config.Provider == "aws" {
		aws.SetupAWS(config)
	}

	if config.ScanFilter == "" {
		os.Unsetenv("DCTL_FILTER")
	} else {
		os.Setenv("DCTL_FILTER", config.ScanFilter)
	}

	exitCode = newDriftctlScan(item, config)
	var err error
	if exitCode == 1 || exitCode == 0 {

		htmlReport, err = os.ReadFile(fmt.Sprintf("/tmp/%s.html", item))
		if err != nil {
			logrus.Error(err.Error())
		}
		jsonReport, err = os.ReadFile(fmt.Sprintf("/tmp/%s.json", item))
		if err != nil {
			logrus.Error(err.Error())
		}
	}
	return exitCode, htmlReport, jsonReport
}
