package builds

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pipeline-metrics-api/internal/httpwrapper"
)

//Builds represents the build log for a single project in TeamCity
type Builds struct {
	BuildList []struct {
		ID        int    `json:"id"`
		Number    string `json:"number"`
		Status    string `json:"status"`
		StartDate string `json:"startDate"`
	} `json:"build"`
}

var readFile = ioutil.ReadFile
var get = httpwrapper.Get

//GetBuilds calls the TeamCityAPI and returns the parsed response
func GetBuilds(teamCityID, buildTypeID, startDate string) (Builds, error) {
	requestURL := strings.Join([]string{
		"http://pipelines.stockport.gov.uk:1980/httpAuth/app/rest/builds?locator=project:",
		teamCityID,
		",buildType:",
		buildTypeID,
		",sinceDate:",
		startDate,
		"%2B0000",
		"&fields=build(id,number,status,startDate)"}, "")

	headers := []httpwrapper.Header{
		httpwrapper.Header{
			Name:  "Accept",
			Value: "application/json"},
		httpwrapper.GenerateAuthenticationHeader(os.Getenv("TeamCityUsername"), os.Getenv("TeamCityPassword")),
	}

	response, err := get(requestURL, &headers)
	return parseResponse(response), err
}

//GetDevelopmentBuilds returns the parsed development data for use when there is no access to TeamCity
func GetDevelopmentBuilds() Builds {
	pwd, _ := os.Getwd()
	response, _ := readFile(pwd + "/config/development-data.json")
	return parseResponse(response)
}

func parseResponse(response []byte) Builds {
	var builds Builds
	json.Unmarshal(response, &builds)
	return builds
}
