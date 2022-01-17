package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type WebStoreData struct {
	Link        string  `json:"link"`
	ImgSrc      string  `json:"imgSrc"`
	Opinion     float32 `json:"opinion"`
	ArrivalDate string  `json:"arrivalDate"`
	Price       float32 `json:"price"`
}

func UpdateProductInformationFromInternet(option *ProductOption) {
	fmt.Println(option.Url)

	params := url.Values{}
	params.Add("url", option.Url)

	resp, err := http.PostForm("http://localhost:5000/", params)

	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()

	var data WebStoreData
	reqBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(reqBody, &data)
	fmt.Println("data: ", data)
	//Log the request body
	bodyString := string(reqBody)
	log.Print("json ", bodyString)

	arrivalDate, err := time.Parse("2006-01-02", data.ArrivalDate)
	if err != nil {
		panic(err)
	}

	var priceObj PricePoint
	priceObj.Price = data.Price
	priceObj.CheckDate = time.Now()

	option.Url = data.Link
	option.ImageSrc = data.ImgSrc
	option.Opinion = data.Opinion
	option.Prices = append(option.Prices, priceObj)
	option.ArrivalDate = arrivalDate

}
