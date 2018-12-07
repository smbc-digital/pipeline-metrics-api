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
	"fmt"
)

type Application struct {
	Id int `json:"id"`
	Name string `json:"name"`
	TeamCityId string `json:",omitempty"`
	BuildTypeId string `json:",omitempty"`
}

type ApplicationsList struct {
	Applications []Application
}

func initConfig() (ApplicationsList, error) {
	var applications ApplicationsList

	applicationsConfig, err := os.Open("applications.json")
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

	requestUrl := strings.Join([]string{ "http://pipelines.stockport.gov.uk:1980/httpAuth/app/rest/builds?locator=project:", application.TeamCityId, ",buildType:", application.BuildTypeId, "&fields=build(id,number,status,startDate)"}, "")
	
	client := &http.Client{}

	request, _ := http.NewRequest("GET", requestUrl, nil)
	request.Header.Add("Authorization", "not a cookie")

	resp, _ := client.Do(request)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}

func main() {
	applications, err := initConfig()

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/applications", applications.getAll)
	http.HandleFunc("/applications/", applications.getOne)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
