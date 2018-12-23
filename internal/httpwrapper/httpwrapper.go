package httpwrapper

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
)

//Header is a map of values to attach to a request
type Header struct {
	Name  string
	Value string
}

var client = &http.Client{}

//GenerateAuthenticationHeader returns a Basic Authentication Header
//with passed in username and password
func GenerateAuthenticationHeader(username, password string) Header {
	return Header{
		Name:  "Authorization",
		Value: "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))}
}

//Get performs a get request
func Get(url string, headers *[]Header) ([]byte, error) {
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Add("Accept", "application/json")

	if headers != nil {
		for i := 0; i < len(*headers); i++ {
			request.Header.Add((*headers)[i].Name, (*headers)[i].Value)
		}
	}

	response, err := client.Do(request)

	if err != nil {
		return nil, errors.New("Http response recieved indicates an error has occured")
	}

	defer response.Body.Close()

	byteResponse, _ := ioutil.ReadAll(response.Body)

	return byteResponse, nil
}
