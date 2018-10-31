package godd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// CreationResult struct storing service creation result
type CreationResult struct {
	Status struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"status"`
}

// ServiceRequest hold data for the
// service creation request
type ServiceRequest struct {
	Name        string
	Description string
	Mllib       string
	Type        string
	Connector   string
	Width       int
	Height      int
	Mean        []float64
	Nclasses    int
	GPU         bool
	GPUID       int
	Repository  string
	Extensions  []string
}

// NewService create a *ServiceRequest with default values
func NewService(service *ServiceRequest) *ServiceRequest {
	if service == nil {
		service = &ServiceRequest{
			Description: "",
			Type:        "supervised",
		}
	}
	return service
}

// CreateService create a service
func CreateService(host string, service *ServiceRequest) (result CreationResult, err error) {
	requestCreate := map[string]interface{}{
		"mllib":       service.Mllib,
		"description": service.Description,
		"type":        service.Type,
		"parameters": map[string]interface{}{
			"input": map[string]interface{}{
				"connector": service.Connector,
				"mean":      service.Mean,
				"width":     service.Width,
				"height":    service.Height,
			},
			"mllib": map[string]interface{}{
				"nclasses": service.Nclasses,
				"gpu":      service.GPU,
				"gpuid":    service.GPUID,
			},
		},
		"model": map[string]interface{}{
			"repository": service.Repository,
			"extensions": service.Extensions,
		},
	}

	bytesReq, err := json.Marshal(requestCreate)
	if err != nil {
		return result, err
	}

	// Send HTTP request
	req, err := http.NewRequest("PUT", host+"/services/"+service.Name, bytes.NewBuffer(bytesReq))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// Fill info structure with response data
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, err
}
