package godd

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var httpTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 60 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 60 * time.Second,
}

var httpClient = &http.Client{
	Timeout:   time.Second * 60,
	Transport: httpTransport,
}

// Info structure contain informations fetched by the
// GetInfo function
type Info struct {
	Status struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"status"`
	Head struct {
		Method   string `json:"method"`
		Version  string `json:"version"`
		Branch   string `json:"branch"`
		Commit   string `json:"commit"`
		Services []struct {
			Mltype      string `json:"mltype"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Mllib       string `json:"mllib"`
			Predict     bool   `json:"predict"`
		} `json:"services"`
	} `json:"head"`
}

// GetInfo return an object containing informations from /info
func GetInfo(host string) (info *Info, err error) {
	// Perform GET request on /info
	resp, err := http.Get(host + "/info")
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
