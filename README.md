[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

[![Build Status](https://travis-ci.com/jolibrain/godd.svg?token=RUPYCGKsp5yMHL2ydcwd&branch=master)](https://travis-ci.com/jolibrain/godd) [![Go Report Card](https://goreportcard.com/badge/github.com/jolibrain/godd)](https://goreportcard.com/report/github.com/jolibrain/godd) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/0d3f90ca8f7146248520913e89e37c9e)](https://www.codacy.com/app/jolibrain/godd?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jolibrain/godd&amp;utm_campaign=Badge_Grade) [![GoDoc](https://godoc.org/github.com/jolibrain/godd?status.svg)](https://godoc.org/github.com/jolibrain/godd)

# GoDD
🧠 [DeepDetect](https://github.com/jolibrain/deepdetect) package for easy integration in any Go project

GoDD offer a simple way to use [DeepDetect](https://github.com/jolibrain/deepdetect) in your Go software, by providing a simple interface to communicate with the different API endpoints supported by DeepDetect.

**GoDD currrently only support prediction, not training.**

# Install

`go get -u github.com/jolibrain/godd`

# Examples

[DeepDetect](https://github.com/jolibrain/deepdetect) quickstart with Docker:

`docker pull beniz/deepdetect_cpu`

`docker run -d -p 8080:8080 -v $HOME/deepdetect-models:/opt/my-models beniz/deepdetect_cpu`

`wget https://deepdetect.com/models/voc0712_dd.tar.gz`

`sudo mkdir -p $HOME/deepdetect-models/voc0712 && sudo tar -xvf voc0712_dd.tar.gz -C $HOME/deepdetect-models/voc0712`

---

Get informations on a DeepDetect instance:

```go
// Set DeepDetect host informations
const myDD = "127.0.0.1:8080"

// Retrieve informations
info, err := godd.GetInfo(myDD)
if err != nil {
	fmt.Println(err.Error())
	os.Exit(1)
}

// Display informations
fmt.Println(info)

// Display only the services field
fmt.Println(info.Head.Services)
```

---

Create a service:

```go
// Create a service request structure
var service godd.ServiceRequest

// Specify values for your service creation
service.Name = "imageserv"
service.Description = "object detection service"
service.Type = "supervised"
service.Mllib = "caffe"
service.Parameters.Input.Connector = "image"
service.Parameters.Input.Width = 300
service.Parameters.Input.Height = 300
service.Parameters.Mllib.Nclasses = 21
service.Model.Repository = "/opt/my-models/voc0712/"

// Send the service creation request
creationResult, err := godd.CreateService(myDD, &service)
if err != nil {
	log.Fatal(err)
}

// Check if the service is created
if creationResult.Status.Code == 200 {
	fmt.Println("Service creation: " + creationResult.Status.Msg)
} else {
	fmt.Println("Service creation: " + creationResult.Status.Msg)
}
```

---

Predict:

```go
// Create predict structure for request parameters
var predict godd.PredictRequest

// Specify values for your prediction
predict.Service = "imageserv"
predict.Data = append(predict.Data, "https://t2.ea.ltmcdn.com/fr/images/9/0/0/les_bienfaits_d_avoir_un_chien_1009_600.jpg")
predict.Parameters.Output.Bbox = true
predict.Parameters.Output.ConfidenceThreshold = 0.1

// Execute the prediction
predictResult, err := godd.Predict(myDD, &predict)
if err != nil {
	log.Fatal(err)
}

// Print data of the first object detected
if predictResult.Status.Code == 200 {
	// Print the complete JSON result:
	// fmt.Println(string(predictResult))
	fmt.Println("Category: " + predictResult.Body.Predictions[0].Classes[0.Cat)
	fmt.Println("Probability: " + strconv.FormatFloa(predictResult.Body.Predictions[0].Classes[0].Prob, 'f', 6, 64))
	var bbox, _ = json.Marshal(predictResult.Body.Predictions[0].Classes[0.Bbox)
	fmt.Println("Bbox: " + string(bbox))
} else {
	fmt.Println("Prediction failed: " + predictResult.Status.Msg)
}
```

---

Delete a service:

```go
// Delete service
serviceDeleteStatus, err := godd.DeleteService(myDD, "imageserv")
if err != nil {
	log.Fatal(err)
}

fmt.Println("Service deletion:")
fmt.Println(serviceDeleteStatus)
```

**You can see the full examples in the [examples](https://github.com/jolibrain/godd/tree/master/examples) folder.**
