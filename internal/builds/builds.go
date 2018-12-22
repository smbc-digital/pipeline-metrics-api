package builds

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

//GetBuilds calls the TeamCityAPI and returns the parsed response
func GetBuilds(teamCityID, buildTypeID, startDate string) Builds {
	requestURL := strings.Join([]string{
		"http://pipelines.stockport.gov.uk:1980/httpAuth/app/rest/builds?locator=project:",
		teamCityID,
		",buildType:",
		buildTypeID,
		",sinceDate:",
		startDate,
		"%2B0000",
		"&fields=build(id,number,status,startDate)"}, "")

	return parseResponse(makeRequest(requestURL))
}

//GetDevelopmentBuilds returns the parsed development data for use when there is no access to TeamCity
func GetDevelopmentBuilds() Builds {
	response, _ := ioutil.ReadFile("config/development-data.json")
	return parseResponse(response)
}

func makeRequest(url string) []byte {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(os.Getenv("TeamCityUsername"), os.Getenv("TeamCityPassword"))

	response, _ := client.Do(request)

	defer response.Body.Close()

	byteResponse, _ := ioutil.ReadAll(response.Body)

	return byteResponse
}

func parseResponse(response []byte) Builds {
	var builds Builds
	json.Unmarshal(response, &builds)
	return builds
}
