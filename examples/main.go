package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/CorentinB/godd"
)

const myDD = "http://127.0.0.1:8080"

func main() {
	// Get informations on the instance
	info, err := godd.GetInfo(myDD)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(info)

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

	// Create predict structure for request parameters
	var predict godd.PredictRequest

	// Specify values for your prediction
	predict.Service = "imageserv"
	predict.Data = append(predict.Data, "https://t2.ea.ltmcdn.com/fr/images/9/0/0/les_bienfaits_d_avoir_un_chien_1009_600.jpg")
	predict.Parameters.Output.Bbox = true
	predict.Parameters.Output.ConfidenceThreshold = 0.1

	predictResult, err := godd.Predict(myDD, &predict)
	if err != nil {
		log.Fatal(err)
	}

	if predictResult.Status.Code == 200 {
		// Print the complete JSON result:
		// fmt.Println(string(predictResult))
		fmt.Println("Category: " + predictResult.Body.Predictions[0].Classes[0].Cat)
		fmt.Println("Probability: " + strconv.FormatFloat(predictResult.Body.Predictions[0].Classes[0].Prob, 'f', 6, 64))
		var bbox, _ = json.Marshal(predictResult.Body.Predictions[0].Classes[0].Bbox)
		fmt.Println("Bbox: " + string(bbox))
	} else {
		fmt.Println("Prediction failed: " + predictResult.Status.Msg)
	}

	// Get service informations
	serviceInfoResult, err := godd.GetServiceInfo(myDD, "imageserv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Service informations:")
	fmt.Println(serviceInfoResult)

	// Delete service
	serviceDeleteStatus, err := godd.DeleteService(myDD, "imageserv")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Service deletion:")
	fmt.Println(serviceDeleteStatus)
}
