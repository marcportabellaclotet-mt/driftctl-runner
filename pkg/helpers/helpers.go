package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/config"
	"github.com/sirupsen/logrus"
)

func EnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func KeyExists(key string) bool {
	if _, ok := config.DritfctlRunMap[key]; ok {
		return true
	}
	return false
}

func LookupEnvOrDefaultInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			logrus.Fatalf("LookupEnvOrDefaultInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func LookupEnvOrDefaultString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func GenerateCMDArgs(item string, config config.DritfctlRun) []string {

	cmdArgs := []string{"scan", "-o",
		fmt.Sprintf("html:///tmp/%s.html", item),
		"-o", fmt.Sprintf("json:///tmp/%s.json", item)}
	if strings.Contains(config.ScanFilter, "Attr.") {
		cmdArgs = append(cmdArgs, "--deep")
	}

	tfStateNew := []string{}
	for _, v := range config.TFStateList {
		tfStateNew = append(tfStateNew, "-f")
		tfStateNew = append(tfStateNew, v)
	}
	cmdArgs = append(cmdArgs, tfStateNew...)
	return cmdArgs
}
