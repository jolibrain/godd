package godd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// PredictResult struct storing predictions result
type PredictResult struct {
	Status struct {
		Code int    `json:"code,omitempty"`
		Msg  string `json:"msg,omitempty"`
	} `json:"status,omitempty"`
	Head struct {
		Method  string  `json:"method,omitempty"`
		Service string  `json:"service,omitempty"`
		Time    float64 `json:"time,omitempty"`
	} `json:"head,omitempty"`
	Body struct {
		Predictions []struct {
			Classes []struct {
				Prob float64 `json:"prob,omitempty"`
				Last bool    `json:"last,omitempty"`
				Bbox struct {
					Ymax float64 `json:"ymax,omitempty"`
					Xmax float64 `json:"xmax,omitempty"`
					Ymin float64 `json:"ymin,omitempty"`
					Xmin float64 `json:"xmin,omitempty"`
				} `json:"bbox,omitempty"`
				Mask struct {
					Format string `json:"format,omitempty"`
					Width  int    `json:"width,omitempty"`
					Height int    `json:"height,omitempty"`
					Data   []int  `json:"data,omitempty"`
				} `json:"mask,omitempty"`
				Cat string `json:"cat,omitempty"`
			}
			URI string `json:"uri,omitempty"`
		} `json:"predictions,omitempty"`
	} `json:"body,omitempty"`
}

// PredictRequest hold data for the
// prediction request
type PredictRequest struct {
	// General parameters
	Service    string   `json:"service,omitempty"`
	Data       []string `json:"data,omitempty"`
	Parameters struct {
		Input struct {
			// Image - image
			Width      int     `json:"width,omitempty"`
			Height     int     `json:"height,omitempty"`
			CropWidth  int     `json:"crop_width,omitempty"`
			CropHeight int     `json:"crop_height,omitempty"`
			BW         bool    `json:"bw,omitempty"`
			MeanTF     float64 `json:"mean,omitempty"`
			Mean       []int   `json:"mean,omitempty"`
			STD        float64 `json:"std,omitempty"`
			// CSV - csv
			Ignore    []string  `json:"ignore,omitempty"`
			Separator string    `json:"separator,omitempty"`
			ID        string    `json:"id,omitempty"`
			Scale     bool      `json:"scale,omitempty"`
			MinVals   []float64 `json:"min_vals,omitempty"`
			MaxVals   []float64 `json:"max_vals,omitempty"`
			// MISSING: categoricals_mapping
			// Text - txt
			Count         int    `json:"count,omitempty"`
			MinCount      int    `json:"min_count,omitempty"`
			MinWordLength int    `json:"min_word_length,omitempty"`
			TFIDF         bool   `json:"tfidf,omitempty"`
			Sentences     bool   `json:"sentences,omitempty"`
			Characters    bool   `json:"characters,omitempty"`
			Sequence      int    `json:"sequence,omitempty"`
			ReadForward   bool   `json:"read_forward,omitempty"`
			Alphabet      string `json:"alphabet,omitempty"`
			Sparse        bool   `json:"sparse,omitempty"`
		} `json:"input,omitempty"`
		Output struct {
			Best     int    `json:"best,omitempty"`
			Template string `json:"template,omitempty"`
			Network  *struct {
				URL         string `json:"url,omitempty"`
				HTTPMethod  string `json:"http_method,omitempty"`
				ContentType string `json:"content_type,omitempty"`
			} `json:"network,omitempty"`
			Measure             []float64 `json:"measure,omitempty"`
			ConfidenceThreshold float64   `json:"confidence_threshold,omitempty"`
			Bbox                bool      `json:"bbox,omitempty"`
			Mask                bool      `json:"mask,omitempty"`
			Rois                string    `json:"rois,omitempty"`
			Index               bool      `json:"index,omitempty"`
			BuildIndex          bool      `json:"build_index,omitempty"`
			Search              bool      `json:"search,omitempty"`
			MultiboxRois        bool      `json:"multibox_rois,omitempty"`
			CTC                 bool      `json:"ctc,omitempty"`
		} `json:"output,omitempty"`
		Mllib struct {
		      // Caffe / Caffe2
		      GPU          bool   `json:"gpu,omitempty"`
		      GPUID        []int  `json:"gpuid,omitempty"`
		      ExtractLayer string `json:"extract_layer,omitempty"`
		      // Net or TF
		      TestBatchSize int `json:"test_batch_size,omitempty"`
		      // Tensorflow
		      InputLayer  string `json:"inputlayer,omitempty"`
		      OutputLayer string `json:"outputlayer,omitempty"`
		} `json:"mllib,omitempty"` 
	} `json:"parameters,omitempty"`
}

// Predict perform a /predict call and
// return a PredictResult structure
func Predict(host string, predictRequest *PredictRequest) (result PredictResult, err error) {
	// Turn requestPredict structure into a map for request
	requestPredict := map[string]interface{}{
		"service":    predictRequest.Service,
		"data":       predictRequest.Data,
		"parameters": predictRequest.Parameters,
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
	client := http.Client{Timeout: 500 * time.Second}
	resp, err := client.Do(req)
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
