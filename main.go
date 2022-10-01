package main

import (
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/httpd"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/runner"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	go runner.Start()
	httpd.Start()
}
