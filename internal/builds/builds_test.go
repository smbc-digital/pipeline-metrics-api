package builds

import (
	"fmt"
	"testing"

	"github.com/pipeline-metrics-api/internal/httpwrapper"
)

func TestParseResponse(t *testing.T) {
	result := parseResponse([]byte("{ \"build\": [ { \"id\": 500, \"number\": \"testnumber\", \"status\": \"SUCCESS\", \"startDate\": \"2018\" }]}"))

	if result.BuildList[0].ID != 500 ||
		result.BuildList[0].Number != "testnumber" ||
		result.BuildList[0].StartDate != "2018" ||
		result.BuildList[0].Status != "SUCCESS" {
		t.Error("Build has not been parsed correctly")
	}
}

func TestGetDevelopmentBuilds(t *testing.T) {
	mockIoutilReadFile := func(url string) ([]byte, error) {
		byteArray := []byte("{ \"build\": [ { \"id\": 500, \"number\": \"testnumber\", \"status\": \"SUCCESS\", \"startDate\": \"2018\" }]}")
		return byteArray, nil
	}
	readFile = mockIoutilReadFile

	result := GetDevelopmentBuilds()

	if result.BuildList[0].ID != 500 ||
		result.BuildList[0].Number != "testnumber" ||
		result.BuildList[0].StartDate != "2018" ||
		result.BuildList[0].Status != "SUCCESS" {
		t.Error("Build has not been parsed correctly")
	}
}

func TestGetBuilds(t *testing.T) {
	mockGet := func(url string, headers *[]httpwrapper.Header) ([]byte, error) {
		byteArray := []byte("{ \"build\": [ { \"id\": 232010, \"number\": \"AWS(315).FE(1.2.3660)\", \"status\": \"SUCCESS\", \"startDate\": \"20181220T111046+0000\" }, { \"id\": 230771, \"number\": \"AWS(312).FE(1.2.3596)\", \"status\": \"SUCCESS\", \"startDate\": \"20181212T155611+0000\" }, { \"id\": 230218, \"number\": \"AWS(312).FE(1.2.3582)\", \"status\": \"SUCCESS\", \"startDate\": \"20181210T103318+0000\" }, { \"id\": 228359, \"number\": \"AWS(303).FE(1.2.3549)\", \"status\": \"SUCCESS\", \"startDate\": \"20181127T141427+0000\" } ] }")
		return byteArray, nil
	}
	get = mockGet

	result, _ := GetBuilds("test", "test", "test")

	fmt.Println(result.BuildList[0])
	if result.BuildList[0].ID != 232010 ||
		result.BuildList[0].Number != "AWS(315).FE(1.2.3660)" ||
		result.BuildList[0].StartDate != "20181220T111046+0000" ||
		result.BuildList[0].Status != "SUCCESS" ||
		len(result.BuildList) != 4 {
		t.Error("Incorrent data returned")
	}
}
