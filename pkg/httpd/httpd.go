package httpd

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/config"
	"github.com/marcportabellaclotet-mt/driftctl-runner/pkg/helpers"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

//go:embed assets/*
var assets embed.FS

type Data struct {
	DrifctlRuns []string
}

func Start() {
	logrus.Info("Starting driftctl-runner web server")
	http.HandleFunc("/", Server)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

func Server(w http.ResponseWriter, r *http.Request) {
	parsedUrl := strings.Split(r.URL.Path, "/")[1]
	keyExists := helpers.KeyExists(parsedUrl)
	if r.URL.Path == "/" {
		indexFile, err := assets.ReadFile("assets/index.html")
		if err != nil {
			logrus.Warn(err)
			return
		}
		tmpl, _ := template.New("").Parse(string(indexFile))
		//tmpl := template.Must(template.ParseFiles(string(indexFile)))
		tmpl.Execute(w, Data{DrifctlRuns: maps.Keys(config.DritfctlRunMap)})
		return
	} else if r.URL.Path == "/reloadconfig" {
		config.ReadConfig()
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if !keyExists {
		fmt.Fprintf(w, "Scan for %s is not defined", parsedUrl)
		return
	} else if len(config.DritfctlRunMap[parsedUrl].ReportHTML) > 0 {
		w.Write(config.DritfctlRunMap[parsedUrl].ReportHTML)
		return
	} else {
		fmt.Fprintf(w, "Scan for %s is not ready yet", parsedUrl)
		return
	}
}
