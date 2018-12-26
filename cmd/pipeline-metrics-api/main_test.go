package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/pipeline-metrics-api/internal/builds"

	"github.com/pipeline-metrics-api/internal/applications"
)

func TestGetApplications(t *testing.T) {
	getSupportedPipelines = func() (applications.List, error) {
		list := applications.List{
			Application: []applications.Application{
				applications.Application{ID: 1, Name: "test-name-1", BuildTypeID: "test-build-1", TeamCityID: "test-city-1", LastPublish: "test-publish-1", NumberOfBuilds: 100},
				applications.Application{ID: 2, Name: "test-name-2", BuildTypeID: "test-build-2", TeamCityID: "test-city-2", LastPublish: "test-publish-2", NumberOfBuilds: 101},
			},
		}
		return list, nil
	}

	req, _ := http.NewRequest("GET", "/applications", nil)

	recorder := httptest.NewRecorder()

	http.HandlerFunc(getApplications).ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	expected := []byte(`{"Application":[{"id":1,"name":"test-name-1","lastPublish":"test-publish-1"},{"id":2,"name":"test-name-1","lastPublish":"test-publish-2"}]}`)

	if body := recorder.Body.Bytes(); reflect.DeepEqual(body, expected) {
		t.Errorf("Response body differs. Expected %s .\n Got %s instead", expected, body)
	}
}

func TestGetApplication(t *testing.T) {
	buildList := &builds.Builds{
		BuildList: []struct {
			ID        int    `json:"id"`
			Number    string `json:"number"`
			Status    string `json:"status"`
			StartDate string `json:"startDate"`
		}{
			{1, "test-number-1", "test-status", "test-startdate"},
			{2, "test-number-1", "test-status", "test-startdate"},
		},
	}
	testCases := []struct {
		returnedApplication *applications.Application
		returnedBuilds      *builds.Builds
		applicationError    error
		buildsError         error
		mockDeveloperMode   bool
		expected            interface{}
	}{
		{nil, &builds.Builds{}, errors.New("example error"), nil, false, "example error"},
		{&applications.Application{ID: 1, BuildTypeID: "test-build-type", TeamCityID: "test-city-id", Name: "test-name", LastPublish: "", NumberOfBuilds: 0}, &builds.Builds{}, nil, errors.New("example error"), false, "example error"},
		{&applications.Application{ID: 1, BuildTypeID: "test-build-type", TeamCityID: "test-city-id", Name: "test-name", LastPublish: "", NumberOfBuilds: 0}, buildList, nil, nil, false, applications.Application{BuildTypeID: "test-build-type", TeamCityID: "test-city-id", ID: 1, Name: "test-name", LastPublish: "test-startdate", NumberOfBuilds: 2}},
		{&applications.Application{ID: 1, BuildTypeID: "test-build-type", TeamCityID: "test-city-id", Name: "test-name", LastPublish: "", NumberOfBuilds: 0}, buildList, nil, nil, true, applications.Application{BuildTypeID: "test-build-type", TeamCityID: "test-city-id", ID: 1, Name: "test-name", LastPublish: "test-startdate", NumberOfBuilds: 2}},
	}

	for _, tc := range testCases {
		getApplicationWithID = func(ID int) (*applications.Application, error) {
			return tc.returnedApplication, tc.applicationError
		}
		getBuilds = func(teamCityID, buildTypeID, startDate string) (builds.Builds, error) {
			return *tc.returnedBuilds, tc.buildsError
		}
		getDevelopmentBuilds = func() builds.Builds {
			return *tc.returnedBuilds
		}
		developerMode = tc.mockDeveloperMode

		req, _ := http.NewRequest("GET", "/applications/1", nil)
		recorder := httptest.NewRecorder()

		http.HandlerFunc(getApplication).ServeHTTP(recorder, req)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
		}

		response := []byte(strings.TrimSpace(recorder.Body.String()))
		bytes, _ := json.Marshal(tc.expected)

		if !reflect.DeepEqual(response, bytes) {
			t.Errorf("Somthing went wrong expect %s but recieved this %s", response, bytes)
		}
	}
}

func TestParseApplications(t *testing.T) {
	application := &applications.Application{
		ID:             1,
		BuildTypeID:    "",
		TeamCityID:     "",
		Name:           "",
		NumberOfBuilds: 0,
	}

	builds := builds.Builds{
		BuildList: []struct {
			ID        int    `json:"id"`
			Number    string `json:"number"`
			Status    string `json:"status"`
			StartDate string `json:"startDate"`
		}{
			{1, "test_number_1", "test_status", "test_publishdate"},
			{2, "test_number_2", "test_status", "test_publishdate"},
			{3, "test_number_3", "test_status", "test_publishdate"},
			{4, "test_number_4", "test_status", "test_publishdate"},
		},
	}

	parseApplications(application, builds)

	if application.NumberOfBuilds != 4 || application.LastPublish != "test_publishdate" {
		t.Error("Parsing of application and build data failed")
	}
}
