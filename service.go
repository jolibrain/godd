package godd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Status struct storing requests execution results
type Status struct {
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

// ServiceInfo structure contain informations
// fetched by the GetServiceInfo function
type ServiceInfo struct {
	Status struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"status"`
	Body struct {
		Mllib       string `json:"mllib"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Jobs        []struct {
			Job    int    `json:"job"`
			Status string `json:"status"`
		} `json:"jobs"`
	} `json:"body"`
}

// NewService create a *ServiceRequest,
// it takes a *ServiceRequest as input,
// initialize it with default values,
// then return a *ServiceRequest structure.
func NewService(service *ServiceRequest) *ServiceRequest {
	if service == nil {
		service = &ServiceRequest{
			Description: "",
			Type:        "supervised",
		}
	}
	return service
}

// CreateService create a service using the
// /services endpoint, it takes the host and
// a *ServiceRequest as input and return a
// *CreationResult structure.
func CreateService(host string, service *ServiceRequest) (result *Status, err error) {
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

// GetServiceInfo fetch service informations using
// the /services/<service_name> endpoint
// It takes the host and the service name as input,
// and return a *ServiceInfo structure.
func GetServiceInfo(host string, service string) (info *ServiceInfo, err error) {
	// Perform GET request on /services/<service_name>
	resp, err := http.Get(host + "/services/" + service)
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

// DeleteService delete a service using the
// /services/<service_name> endpoint, it takes
// and host and a service name as input and return
// a *Status structure.
func DeleteService(host string, service string) (status *Status, err error) {
	// Create DELETE request
	req, err := http.NewRequest("DELETE", host+"/services/"+service, nil)
	if err != nil {
		return status, err
	}

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return status, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return status, err
	}

	// Fill info structure with response data
	err = json.Unmarshal(respBody, &status)
	if err != nil {
		return status, err
	}

	return status, nil
}
