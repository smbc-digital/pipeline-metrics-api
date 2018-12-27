package applications

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	application := New(
		"Front End",
		"TeamCityID",
		"BuildTypeID",
		"01122018",
		100,
		1)

	if application.ID != 1 ||
		application.TeamCityID != "TeamCityID" ||
		application.BuildTypeID != "BuildTypeID" ||
		application.NumberOfBuilds != 100 ||
		application.LastPublish != "01122018" {
		t.Error("Application varibles not set correctly")
	}
}

func TestGetSupportedPipelines(t *testing.T) {
	mockReadFile := func(filePath string) ([]byte, error) {
		byteArray := []byte("[{\"id\":1,\"name\":\"DTS-Frontend\",\"teamCityId\":\"Dbd_DtsDeploy\",\"buildTypeId\":\"Dbd_DtsDeploy_6Production\"}]")
		return byteArray, nil
	}

	ioutilReadFile = mockReadFile
	supported, _ := GetSupportedPipelines()

	if supported.Application[0].ID != 1 || supported.Application[0].Name != "DTS-Frontend" {
		t.Error("Supported piplines where not found")
	}
}

func TestGetApplication(t *testing.T) {
	testCases := []struct {
		ID                  int
		expectedApplication *Application
		expectedError       error
	}{
		{1, &Application{1, "DTS-Frontend", "Dbd_DtsDeploy", "Dbd_DtsDeploy_6Production", "", 0}, nil},
		{2, nil, errors.New("Incorrect value given for ApplicationID")},
	}

	mockReadFile := func(filePath string) ([]byte, error) {
		byteArray := []byte("[{\"id\":1,\"name\":\"DTS-Frontend\",\"teamCityId\":\"Dbd_DtsDeploy\",\"buildTypeId\":\"Dbd_DtsDeploy_6Production\"}]")
		return byteArray, nil
	}
	ioutilReadFile = mockReadFile

	for _, tc := range testCases {
		application, err := GetApplication(tc.ID)

		if application != nil && tc.expectedApplication != nil {
			if application.ID != tc.expectedApplication.ID ||
				application.BuildTypeID != tc.expectedApplication.BuildTypeID ||
				application.TeamCityID != tc.expectedApplication.TeamCityID {
				t.Error("GetApplication returned an incorrect result")
			}
		}

		if err != nil && tc.expectedError != nil {
			if err.Error() != tc.expectedError.Error() {
				t.Error("GetApplication returned an incorrect error")
			}
		}
	}
}
