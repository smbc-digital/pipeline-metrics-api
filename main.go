package main

import (
	"net/http"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"time"
	"fmt"
	"flag"
)

var developerMode *bool

type Application struct {
	Id int `json:"id"`
	Name string `json:"name"`
	TeamCityId string `json:",omitempty"`
	BuildTypeId string `json:",omitempty"`
	LastPublish string `json:"lastPublish,omitempty"`
	NumberOfBuilds int `json:"numberOfBuilds,omitempty"`
}

type ApplicationsList struct {
	Applications []Application
}

type Build struct {
	BuildList []struct {
		ID        int    `json:"id"`
		Number    string `json:"number"`
		Status    string `json:"status"`
		StartDate string `json:"startDate"`
	} `json:"build"`
}

func initConfig() (ApplicationsList, error) {
	var applications ApplicationsList

	applicationsConfig, err := os.Open("config/supported-pipelines.json")
	if err != nil {
		return applications, err
	}
	byteValue, _ := ioutil.ReadAll(applicationsConfig)

	json.Unmarshal(byteValue, &applications.Applications)

	return applications, nil
}

func (applications *ApplicationsList) getAll(w http.ResponseWriter, r *http.Request) {
	var list = *applications

	for i, _ := range list.Applications {
		list.Applications[i].TeamCityId = ""
		list.Applications[i].BuildTypeId = ""
	}

	json.NewEncoder(w).Encode(list)
}

func (applications *ApplicationsList) getOne(w http.ResponseWriter, r *http.Request) {
	applicationId, _ := strconv.Atoi(path.Base(r.URL.Path))
	var application *Application

	for _, app := range applications.Applications {
		if app.Id == applicationId {
			application = &app
			break
		}
	}

	if application == nil {
		json.NewEncoder(w).Encode("Application with this ID doesn't exist.")
		return
	}

	var startDate = time.Now().AddDate(0, 0, -28).Format("20060102T150405")

	response := doRequest(application.TeamCityId, application.BuildTypeId, startDate)
	
	var build Build
	err := json.Unmarshal(response, &build)

	if err != nil {
		panic(err)
	}

	parseApplications(application, build)

	json.NewEncoder(w).Encode(application)
}

func parseApplications(application *Application, build Build){
	application.LastPublish = build.BuildList[0].StartDate
	application.NumberOfBuilds = len(build.BuildList)
}

func doRequest(teamCityId, buildTypeId, startDate string) []byte{
	if *developerMode {
		lan, _ := ioutil.ReadFile("config/development-data.json")
		return lan
	}

	requestUrl := strings.Join([]string{ "http://pipelines.stockport.gov.uk:1980/httpAuth/app/rest/builds?locator=project:", teamCityId, ",buildType:", buildTypeId,",sinceDate:", startDate, "%2B0000", "&fields=build(id,number,status,startDate)"}, "")
	client := &http.Client{}
	
	request, _ := http.NewRequest("GET", requestUrl, nil)
	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(os.Getenv("TeamCityUsername"),os.Getenv("TeamCityPassword"))
	
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	
	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes
}

func main() {
	applications, err := initConfig()
	developerMode = flag.Bool("dev", false, "Set to true if you do not have direct access to TeamCity")
	flag.Parse()

	if *developerMode == true {
		fmt.Println("You are a developer!")
	}

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/applications", applications.getAll)
	http.HandleFunc("/applications/", applications.getOne)

	fmt.Println("Now listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
