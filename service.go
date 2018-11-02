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

// ServiceRequest hold data for the service creation request
// See: https://deepdetect.com/api/#create-a-service
type ServiceRequest struct {
	// General parameters
	Name        string
	Mllib       string `json:"mllib"`
	Type        string `json:"type"`
	Description string `json:"description"`
	// Model
	Model struct {
		Repository       string   `json:"repository"`
		Templates        string   `json:"templates"`
		Weights          string   `json:"weights"`
		CreateRepository bool     `json:"create_repository"`
		IndexPreload     bool     `json:"index_preload"`
		Extensions       []string `json:"extensions"`
	} `json:"model"`
	Parameters struct {
		Input struct {
			// Input
			Connector string `json:"connector"`
			// Input - image
			Width         int       `json:"width"`
			Height        int       `json:"height"`
			BW            bool      `json:"bw"`
			MeanTS        float64   `json:"mean"`
			Mean          []float64 `json:"mean"`
			STD           float64   `json:"std"`
			Segmentation  bool      `json:"segmentation"`
			MultiLabel    bool      `json:"multi_label"`
			RootFolder    string    `json:"root_folder"`
			CTC           bool      `json:"ctc"`
			UnchangedData bool      `json:"unchanged_data"`
			Bbox          bool      `json:"bbox"`
			// Input - CSV
			Label        string   `json:"label"`
			Ignore       []string `json:"ignore"`
			LabelOffset  int      `json:"label_offset"`
			Separator    string   `json:"separator"`
			ID           string   `json:"id"`
			Scale        bool     `json:"scale"`
			Categoricals []string `json:"categoricals"`
			DB           bool     `json:"db"`
			// Input - TXT
			Sentences   bool   `json:"sentences"`
			Characters  bool   `json:"characters"`
			Sequence    int    `json:"sequence"`
			ReadForward bool   `json:"read_forward"`
			Alphabet    string `json:"alphabet"`
			Sparse      bool   `json:"sparse"`
		} `json:"input"`
		Mllib struct {
			// Caffe and Caffe2
			Nclasses           int      `json:"nclasses"`
			Ntargets           int      `json:"ntargets"`
			GPU                bool     `json:"gpu"`
			GPUID              []int      `json:"gpuid"`
			Template           string   `json:"template"`
			LayersMLP          []int    `json:"layers"`
			LayersConvnet      []string `json:"layers"`
			Activation         string   `json:"activation"`
			Dropout            float64  `json:"dropout"`
			Regression         bool     `json:"regression"`
			Autoencoder        bool     `json:"autoencoder"`
			CropSize           int      `json:"crop_size"`
			Rotate             bool     `json:"rotate"`
			Mirror             bool     `json:"mirror"`
			Finetuning         bool     `json:"finetuning"`
			DB                 bool     `json:"db"`
			ScalingTemperature float64  `json:"scaling_temperature"`
			Loss               string   `json:"loss"`
			// Noise - images only
			Prob         float64 `json:"prob"`
			AllEffects   bool    `json:"all_effects"`
			Decolorize   bool    `json:"decolorize"`
			HistEQ       bool    `json:"hist_eq"`
			Inverse      bool    `json:"inverse"`
			GaussBlur    bool    `json:"gauss_blur"`
			Posterize    bool    `json:"posterize"`
			Erode        bool    `json:"erode"`
			Saltpepper   bool    `json:"saltpepper"`
			Clahe        bool    `json:"clahe"`
			ConvertToHSV bool    `json:"convert_to_hsv"`
			ConvertToLAB bool    `json:"convert_to_lab"`
			// Distort - images only
			Brightness     bool `json:"brightness"`
			Contrast       bool `json:"contrast"`
			Saturation     bool `json:"saturation"`
			HUE            bool `json:"HUE"`
			RandomOrdering bool `json:"random ordering"`
			// TensorFlow
			InputLayer  string `json:"inputlayer"`
			OutputLayer string `json:"outputlayer"`
		}
		Output struct {
			StoreConfig bool `json:"store_config"`
		} `json:"output"`
	} `json:"parameters"`
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
		"parameters":  service.Parameters,
		"model":       service.Model,
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
