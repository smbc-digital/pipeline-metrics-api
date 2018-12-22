package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/pipeline-metrics-api/internal/applications"
	"github.com/pipeline-metrics-api/internal/builds"
)

var developerMode *bool

func getApplications(w http.ResponseWriter, r *http.Request) {
	list, _ := applications.GetSupportedPipelines()

	for i := range list.Application {
		list.Application[i].TeamCityID = ""
		list.Application[i].BuildTypeID = ""
	}

	json.NewEncoder(w).Encode(list)
}

func getApplication(w http.ResponseWriter, r *http.Request) {
	ID, _ := strconv.Atoi(path.Base(r.URL.Path))
	application, err := applications.GetApplication(ID)

	fmt.Println(err)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	var startDate = time.Now().AddDate(0, 0, -28).Format("20060102T150405")

	if !*developerMode {
		parseApplications(application, builds.GetBuilds(application.TeamCityID, application.BuildTypeID, startDate))
	} else {
		parseApplications(application, builds.GetDevelopmentBuilds())
	}

	json.NewEncoder(w).Encode(application)
}

func parseApplications(application *applications.Application, build builds.Builds) {
	application.LastPublish = build.BuildList[0].StartDate
	application.NumberOfBuilds = len(build.BuildList)
}

func main() {
	developerMode = flag.Bool("dev", false, "Set to true if you do not have direct access to TeamCity")
	flag.Parse()

	http.HandleFunc("/applications", getApplications)
	http.HandleFunc("/applications/", getApplication)

	fmt.Println("Now listening on port http//localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
