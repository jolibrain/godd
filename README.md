[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

[![Build Status](https://travis-ci.com/CorentinB/godd.svg?token=RUPYCGKsp5yMHL2ydcwd&branch=master)](https://travis-ci.com/CorentinB/godd) [![Go Report Card](https://goreportcard.com/badge/github.com/CorentinB/godd)](https://goreportcard.com/report/github.com/CorentinB/godd) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/0d3f90ca8f7146248520913e89e37c9e)](https://www.codacy.com/app/CorentinB/godd?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=CorentinB/godd&amp;utm_campaign=Badge_Grade) [![GoDoc](https://godoc.org/github.com/CorentinB/godd?status.svg)](https://godoc.org/github.com/CorentinB/godd)

# GoDD
🧠 [DeepDetect](https://github.com/jolibrain/deepdetect) package for easy integration in any Go project

GoDD offer a simple way to use [DeepDetect](https://github.com/jolibrain/deepdetect) in your Go software, by providing a simple interface to communicate with the different API endpoints supported by DeepDetect.

# Install

`go get -u -d github.com/CorentinB/godd`

# Examples

Get informations on a DeepDetect instance:

```go
// Retrieve informations
info, err := godd.GetInfo(myHost)
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
service.Name = "imagenet"
service.Description = "Example of service creation with godd"
service.Mllib = "caffe"
service.Parameters.Input.Connector = "image"
service.Parameters.Input.Width = 300
service.Parameters.Input.Height = 300
service.Parameters.Mllib.Nclasses = 601
service.Parameters.Mllib.GPU = true
service.Model.Repository = "/home/corentin/my_model/"

// Send the service creation request
creationResult, err := godd.CreateService(myHost, &service)
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
predict.Service = "imagenet"
predict.Data = append(predict.Data, "https://t2.ea.ltmcdn.com/fr/images/9/0/0/les_bienfaits_d_avoir_un_chien_1009_600.jpg")

predictResult, err := godd.Predict(myHost, &predict)
if err != nil {
	log.Fatal(err)
}

if predictResult.Status.Code == 200 {
    fmt.Println("Category: " + predictResult.Body.Predictions[0].Classes.Cat)
    fmt.Println("Probability: " + predictResult.Body.Predictions[0].Classes.Prob)
} else {
	fmt.Println("Prediction failed: " + predictResult.Status.Msg)
}
```

---

Delete a service:

```go
// Delete service
serviceDeleteStatus, err := godd.DeleteService(myHost, "mask")
if err != nil {
	log.Fatal(err)
}

fmt.Println("Service deletion:")
fmt.Println(serviceDeleteStatus)
```

**You can see more examples in the [examples](https://github.com/CorentinB/godd/tree/master/examples) folder.**

# Todo list

- [X] /services
- [X] /predict
- [X] /info
- [ ] /train