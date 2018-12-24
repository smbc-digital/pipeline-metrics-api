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
		byteArray := []byte("{ \"build\": [ { \"id\": 1001, \"number\": \"build_number_1001\", \"status\": \"test_status\", \"startDate\": \"2018_01_04\" }, { \"id\": 1002, \"number\": \"build_number_1002\", \"status\": \"test_status\", \"startDate\": \"2018_01_03\" }, { \"id\": 1003, \"number\": \"build_number_1003\", \"status\": \"test_status\", \"startDate\": \"2018_01_02\" }, { \"id\": 1004, \"number\": \"build_number_1004\", \"status\": \"test_status\", \"startDate\": \"2018_01_01\" } ] }")
		return byteArray, nil
	}
	get = mockGet

	result, _ := GetBuilds("test", "test", "test")

	fmt.Println(result.BuildList[0])
	if result.BuildList[0].ID != 1001 ||
		result.BuildList[0].Number != "build_number_1001" ||
		result.BuildList[0].StartDate != "2018_01_04" ||
		result.BuildList[0].Status != "test_status" ||
		len(result.BuildList) != 4 {
		t.Error("Incorrent data returned")
	}
}
