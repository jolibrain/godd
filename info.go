package godd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Info structure contain informations fetched by the
// GetInfo function
type Info struct {
	Status InfoStatus
	Head   InfoHead
}

// InfoStatus structure contain the status part of the
// /info response
type InfoStatus struct {
	Code int
	Msg  string
}

// InfoHead structure contain the head part of the
// /info response
type InfoHead struct {
	Method   string
	Version  string
	Branch   string
	Commit   string
	Services []InfoServices
}

// InfoServices structure contain the services array
// of the head part of the /info response
type InfoServices struct {
	Mltype      string
	Name        string
	Description string
	Mllib       string
	Predict     bool
}

// GetInfo return an object containing informations from /info
func GetInfo(host string, port string) (info Info, err error) {
	// Perform GET request on /info
	resp, err := http.Get("http://" + host + ":" + port + "/info")
	if err != nil {
		return info, err
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, err
	}

	// Fill info structure with response data
	err = json.Unmarshal(body, &info)
	if err != nil {
		return info, err
	}

	return info, nil
}
