package godd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var httpTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 10 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

var httpClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: httpTransport,
}

// PredictResult struct storing predictions result
type PredictResult struct {
	Status struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"status"`
	Head struct {
		Method  string  `json:"method"`
		Service string  `json:"service"`
		Time    float64 `json:"time"`
	} `json:"head"`
	Body struct {
		Predictions []struct {
			Classes []struct {
				Prob float64 `json:"prob"`
				Last bool    `json:"last"`
				Bbox struct {
					Ymax float64 `json:"ymax"`
					Xmax float64 `json:"xmax"`
					Ymin float64 `json:"ymin"`
					Xmin float64 `json:"xmin"`
				} `json:"bbox"`
				Mask struct {
					Format string `json:"format"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
					Data   []int  `json:"data"`
				} `json:"mask"`
				Cat string `json:"cat"`
			}
			URI string `json:"uri"`
		} `json:"predictions"`
	} `json:"body"`
}

// PredictRequest hold data for the /predict
// request
type PredictRequest struct {
	Service             string
	Width               int
	Height              int
	Best                int
	Bbox                bool
	Mask                bool
	Data                []string
	ConfidenceThreshold float64
}

// NewPredict create a PredictRequest object with default values
func NewPredict(predict *PredictRequest) *PredictRequest {
	if predict == nil {
		predict = &PredictRequest{
			Service:             "imageserv",
			Width:               227,
			Height:              227,
			Best:                1,
			Bbox:                false,
			Mask:                false,
			ConfidenceThreshold: 0.1,
		}
	}
	return predict
}

// Predict perform a /predict call and
// return a PredictResult structure
func Predict(host string, predictRequest *PredictRequest) (result PredictResult, err error) {
	// Turn requestPredict structure into a map for request
	requestPredict := map[string]interface{}{
		"service": predictRequest.Service,
		"parameters": map[string]interface{}{
			"input": map[string]interface{}{
				"width":  predictRequest.Width,
				"height": predictRequest.Height,
			},
			"output": map[string]interface{}{
				"best":                 predictRequest.Best,
				"bbox":                 predictRequest.Bbox,
				"mask":                 predictRequest.Mask,
				"confidence_threshold": predictRequest.ConfidenceThreshold,
			},
		},
		"data": predictRequest.Data,
	}

	// Marshal data
	bytesReq, err := json.Marshal(requestPredict)
	if err != nil {
		return result, err
	}

	// Send HTTP request
	req, err := http.NewRequest("POST", host+"/predict", bytes.NewBuffer(bytesReq))
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
