package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type messageReponse struct {
	Type    string `json:"Type"`
	Message string `json:"Message"`
}

type dateWrapper struct {
	date time.Time
}

func (this *dateWrapper) UnmarshalJSON(data []byte) error {
	var dateString string
	if err := json.Unmarshal(data, &dateString); err != nil {
		return err
	}
	myDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		panic(err)
	}
	this.date = myDate
	fmt.Println("date: ", dateString, myDate, this.date)
	return nil
}

type productOptionRequestObject struct {
	Opinion     float32     `json:"Opinion"`
	ArrivalDate dateWrapper `json:"ArrivalDate"`
	ImageSrc    string      `json:"ImageSrc"`
	Url         string      `json:"Url"`
}
type agregarOpcionRequestObject struct {
	ProductID int           `json:"productID"`
	Option    ProductOption `json:"option"`
}
