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

var developerMode = false
var getSupportedPipelines = applications.GetSupportedPipelines
var getApplicationWithID = applications.GetApplication
var getBuilds = builds.GetBuilds
var getDevelopmentBuilds = builds.GetDevelopmentBuilds

func getApplications(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	list, _ := getSupportedPipelines()

	for i := range list.Application {
		list.Application[i].TeamCityID = ""
		list.Application[i].BuildTypeID = ""
		list.Application[i].NumberOfBuilds = 0
	}

	json.NewEncoder(w).Encode(list)
}

func getApplication(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	ID, _ := strconv.Atoi(path.Base(r.URL.Path))

	application, err := getApplicationWithID(ID)

	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	var startDate = time.Now().AddDate(0, 0, -28).Format("20060102T150405")

	if !developerMode {
		builds, err := getBuilds(application.TeamCityID, application.BuildTypeID, startDate)

		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		parseApplications(application, builds)
	} else {
		parseApplications(application, getDevelopmentBuilds())
	}

	json.NewEncoder(w).Encode(application)
}

func parseApplications(application *applications.Application, build builds.Builds) {
	layout := "20060102T150405+0000"
	parsedLastPublish, err := time.Parse(layout, build.BuildList[0].StartDate)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(parsedLastPublish)

	application.LastPublish = parsedLastPublish.String()
	application.NumberOfBuilds = len(build.BuildList)
}

func main() {
	var developerModeFlag = flag.Bool("dev", false, "Set to true if you do not have direct access to TeamCity")
	flag.Parse()

	developerMode = *developerModeFlag
	http.HandleFunc("/applications", getApplications)
	http.HandleFunc("/applications/", getApplication)

	fmt.Println("Now listening on port http//localhost:4040")
	log.Fatal(http.ListenAndServe(":4040", nil))
}
