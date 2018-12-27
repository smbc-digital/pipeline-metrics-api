package applications

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//List represents a list of Application
type List struct {
	Application []Application
}

//Application represents the build details of a project in TeamCity
//LastPublish, NumberOfBuilds, TeamCityID, and BuildTypeID can be empty
type Application struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	TeamCityID     string `json:",omitempty"`
	BuildTypeID    string `json:",omitempty"`
	LastPublish    string `json:"lastPublish,omitempty"`
	NumberOfBuilds int    `json:"numberOfBuilds,omitempty"`
}

var osGetwd = os.Getwd
var ioutilReadFile = ioutil.ReadFile

//New creates a new instance of Application with the specifed parameters
func New(name, teamCityID, buildTypeID, lastPublish string, numberOfBuilds, id int) Application {
	app := Application{
		id,
		name,
		teamCityID,
		buildTypeID,
		lastPublish,
		numberOfBuilds}

	return app
}

//GetSupportedPipelines reads in the preconfigured supported applications
func GetSupportedPipelines() (List, error) {
	var list List

	pwd, _ := osGetwd()
	byteValue, _ := ioutilReadFile(pwd + "/config/supported-pipelines.json")

	json.Unmarshal(byteValue, &list.Application)

	return list, nil
}

//GetApplication return a single application that matches the ID specifed
func GetApplication(ID int) (*Application, error) {
	list, _ := GetSupportedPipelines()
	for _, app := range list.Application {
		if app.ID == ID {
			return &app, nil
		}
	}

	return nil, errors.New("Incorrect value given for ApplicationID")
}
