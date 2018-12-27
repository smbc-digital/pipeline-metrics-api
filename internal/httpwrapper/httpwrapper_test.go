package httpwrapper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	testCase := []struct {
		url                  string
		headers              *[]Header
		expectedError        error
		expectedResponseBody []byte
	}{
		{"successfull", &[]Header{Header{Name: "", Value: ""}}, nil, []byte("example body response")},
		{"failure", nil, errors.New("Http response recieved indicates an error has occured"), nil},
	}

	makeRequest = func(request *http.Request) (*http.Response, error) {
		body := ioutil.NopCloser(bytes.NewBufferString("example body response"))

		response := http.Response{
			StatusCode: http.StatusFound,
			Request:    nil,
			Body:       body,
		}

		if request.URL.String() == "failure" {
			return nil, errors.New("example failed response")
		}
		return &response, nil
	}

	for _, tc := range testCase {
		result, err := Get(tc.url, tc.headers)

		if !reflect.DeepEqual(result, tc.expectedResponseBody) {
			t.Error("Incorrect responce returned from httpwrapper.Get")
		}

		if err != nil &&
			tc.expectedError != nil &&
			err.Error() != tc.expectedError.Error() {
			t.Error("Incorrect response returned from httpwrapper.Get")
		}
	}
}

func TestGenerateAuthenticationHeader(t *testing.T) {
	result := GenerateAuthenticationHeader("test_username", "test_password")

	if result.Name != "Authorization" || result.Value != "Basic "+base64.StdEncoding.EncodeToString([]byte("test_username:test_password")) {
		t.Error("Incorrect authentication headers returend")
	}
}
