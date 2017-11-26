/*
Fetch the NameAPI.org REST API and turn JSON response into a Go struct.

Sent data have to be JSON data encoded into request body.
Send request headers must be set to 'application/json'.
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// url of the NameAPI.org endpoint:
const (
	url = "http://rc50-api.nameapi.org/rest/v5.0/parser/personnameparser?" +
		"apiKey=<API-KEY>"
)

func main() {

	// JSON string to be sent to NameAPI.org:
	jsonString := `{
        "inputPerson": {
            "type": "NaturalInputPerson",
            "personName": {
                "nameFields": [
                    {
                        "string": "Petra",
                        "fieldType": "GIVENNAME"
                    }, {
                        "string": "Meyer",
                        "fieldType": "SURNAME"
                    }
                ]
            },
            "gender": "UNKNOWN"
        }
    }`
	// Convert JSON string to NewReader (expected by NewRequest)
	jsonBody := strings.NewReader(jsonString)

	// Need to create a client in order to modify headers
	// and set content-type to 'application/json':
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, jsonBody)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	// Proceed only if no error:
	switch {
	default:
		// Create a struct dedicated to receiving the fetched
		// JSON content:
		type Level5 struct {
			String   string `json:"string"`
			TermType string `json:"termType"`
		}
		type Level41 struct {
			Gender     string  `json:"gender"`
			Confidence float64 `json:"confidence"`
		}
		type Level42 struct {
			Terms []Level5 `json:"terms"`
		}
		type Level3 struct {
			Gender           Level41 `json:"gender"`
			OutputPersonName Level42 `json:"outputPersonName"`
		}
		type Level2 struct {
			ParsedPerson Level3 `json:"parsedPerson"`
		}
		type RespContent struct {
			Matches []Level2 `json:"matches"`
		}

		// Decode fetched JSON and put it into respContent:
		respContentBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		var respContent RespContent
		err = json.Unmarshal(respContentBytes, &respContent)
		if err != nil {
			log.Println(err)
		}
		log.Println(respContent)
	case err != nil:
		log.Println("Network error:", err)
	case resp.StatusCode != 200:
		log.Println("Bad HTTP status code:", err)
	}

}
