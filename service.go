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
		Code int    `json:"code,omitempty"`
		Msg  string `json:"msg,omitempty"`
	} `json:"status,omitempty"`
}

// ServiceRequest hold data for the service creation request
// See: https://deepdetect.com/api/#create-a-service
type ServiceRequest struct {
	// General parameters
	Name        string
	Mllib       string `json:"mllib,omitempty"`
	Type        string `json:"type,default=supervised"`
	Description string `json:"description,omitempty"`
	// Model
	Model struct {
		Repository       string   `json:"repository,omitempty"`
		Templates        string   `json:"templates,omitempty"`
		Init             string   `json:"init,omitempty"`
		Weights          string   `json:"weights,omitempty"`
		CreateRepository bool     `json:"create_repository,omitempty"`
		IndexPreload     bool     `json:"index_preload,omitempty"`
		Extensions       []string `json:"extensions,omitempty"`
	} `json:"model,omitempty"`
	Parameters struct {
		Input struct {
			// Input
			Connector string `json:"connector,omitempty"`
			// Input - image
			Width         int       `json:"width,omitempty"`
			Height        int       `json:"height,omitempty"`
			BW            bool      `json:"bw,omitempty"`
			MeanTS        float64   `json:"mean,omitempty"`
			Mean          []float64 `json:"mean,omitempty"`
			STD           float64   `json:"std,omitempty"`
			Segmentation  bool      `json:"segmentation,omitempty"`
			MultiLabel    bool      `json:"multi_label,omitempty"`
			RootFolder    string    `json:"root_folder,omitempty"`
			CTC           bool      `json:"ctc,omitempty"`
			UnchangedData bool      `json:"unchanged_data,omitempty"`
			Bbox          bool      `json:"bbox,omitempty"`
			// Input - CSV
			Label        string   `json:"label,omitempty"`
			Ignore       []string `json:"ignore,omitempty"`
			LabelOffset  int      `json:"label_offset,omitempty"`
			Separator    string   `json:"separator,omitempty"`
			ID           string   `json:"id,omitempty"`
			Scale        bool     `json:"scale,omitempty"`
			Categoricals []string `json:"categoricals,omitempty"`
			DB           bool     `json:"db,omitempty"`
			// Input - TXT
			Sentences   bool   `json:"sentences,omitempty"`
			Characters  bool   `json:"characters,omitempty"`
			Sequence    int    `json:"sequence,omitempty"`
			ReadForward bool   `json:"read_forward,omitempty"`
			Alphabet    string `json:"alphabet,omitempty"`
			Sparse      bool   `json:"sparse,omitempty"`
		} `json:"input,omitempty"`
		Mllib struct {
			// Caffe and Caffe2
			Nclasses           int      `json:"nclasses,omitempty"`
			Ntargets           int      `json:"ntargets,omitempty"`
			GPU                bool     `json:"gpu,omitempty"`
			GPUID              []int    `json:"gpuid,omitempty"`
			Template           string   `json:"template,omitempty"`
			LayersMLP          []int    `json:"layers,omitempty"`
			LayersConvnet      []string `json:"layers,omitempty"`
			Activation         string   `json:"activation,omitempty"`
			Dropout            float64  `json:"dropout,omitempty"`
			Regression         bool     `json:"regression,omitempty"`
			Autoencoder        bool     `json:"autoencoder,omitempty"`
			CropSize           int      `json:"crop_size,omitempty"`
			Rotate             bool     `json:"rotate,omitempty"`
			Mirror             bool     `json:"mirror,omitempty"`
			Finetuning         bool     `json:"finetuning,omitempty"`
			DB                 bool     `json:"db,omitempty"`
			ScalingTemperature float64  `json:"scaling_temperature,omitempty"`
			Loss               string   `json:"loss,omitempty"`
			// Noise - images only
			Prob         float64 `json:"prob,omitempty"`
			AllEffects   bool    `json:"all_effects,omitempty"`
			Decolorize   bool    `json:"decolorize,omitempty"`
			HistEQ       bool    `json:"hist_eq,omitempty"`
			Inverse      bool    `json:"inverse,omitempty"`
			GaussBlur    bool    `json:"gauss_blur,omitempty"`
			Posterize    bool    `json:"posterize,omitempty"`
			Erode        bool    `json:"erode,omitempty"`
			Saltpepper   bool    `json:"saltpepper,omitempty"`
			Clahe        bool    `json:"clahe,omitempty"`
			ConvertToHSV bool    `json:"convert_to_hsv,omitempty"`
			ConvertToLAB bool    `json:"convert_to_lab,omitempty"`
			// Distort - images only
			Brightness     bool `json:"brightness,omitempty"`
			Contrast       bool `json:"contrast,omitempty"`
			Saturation     bool `json:"saturation,omitempty"`
			HUE            bool `json:"HUE,omitempty"`
			RandomOrdering bool `json:"random ordering,omitempty,omitempty"`
			// TensorFlow
			InputLayer  string `json:"inputlayer,omitempty"`
			OutputLayer string `json:"outputlayer,omitempty"`
			// Memory
			Datatype 	 string   `json:"datatype,omitempty"`
			MaxBatchSize 	 int      `json:"maxBatchSize,omitempty"`
			MaxWorkspaceSize int      `json:"maxWorkspaceSize,omitempty"`
		} `json:"mllib,omitempty"`
		Output struct {
			StoreConfig bool `json:"store_config,omitempty"`
		} `json:"output,omitempty"`
	} `json:"parameters,omitempty"`
}

// ServiceInfo structure contain informations
// fetched by the GetServiceInfo function
type ServiceInfo struct {
	Status struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"msg,omitempty"`
	} `json:"status,omitempty"`
	Body struct {
		Mllib       string `json:"mllib,omitempty"`
		Description string `json:"description,omitempty"`
		Name        string `json:"name,omitempty"`
		Jobs        []struct {
			Job    int    `json:"job,omitempty"`
			Status string `json:"status,omitempty"`
		} `json:"jobs,omitempty"`
	} `json:"body,omitempty"`
}

// CreateService create a service using the
// /services endpoint, it takes the host and
// a *ServiceRequest as input and return a
// *CreationResult structure.
func CreateService(host string, service *ServiceRequest) (result *Status, err error) {
	if service.Type == "" {
		service.Type = "supervised"
	}
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
// the host and a service name as input and return
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
