package main

import (
	"net/http"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"
)

type Application struct {
	Id int
	Name string
	TeamCityId string `json:",omitempty"`
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

func (applications *ApplicationsList) get(w http.ResponseWriter, r *http.Request) {
	var list = *applications

	for i, _ := range list.Applications {
		list.Applications[i].TeamCityId = ""
	}

	json.NewEncoder(w).Encode(list)
}

func main() {
	applications, err := initConfig()

	if err != nil {
		panic(err)
	}

	fmt.Println(applications)

	http.HandleFunc("/applications", applications.get)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
