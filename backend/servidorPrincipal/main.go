package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var database *gorm.DB
var priceTarackerService RepitableEventHandler

func main() {
	priceTarackerService = makeRepitableEventHandler()
	database = getDataBaseConnection() //stablish connection with the datasabe

	startTrackingPrices()

	r := mux.NewRouter()
	//serving static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(""))))

	r.HandleFunc("/", mainHandler).Methods("GET") // render the landing page, all the products
	// r.HandleFunc("/{producto}", getOpecionesDeProducto).Methods("GET")
	// r.HandleFunc("/{producto}/{vendedor}", getOferta).Methods("GET")

	r.HandleFunc("/agregarProducto", agregarProducto).Methods("POST")
	r.HandleFunc("/eliminarProducto", eliminarProducto).Methods("DELETE")
	r.HandleFunc("/obtenerProductos", obtenerProductos).Methods("GET")

	r.HandleFunc("/agregarOpcion", agregarOpcion).Methods("POST")
	r.HandleFunc("/eliminarOpcion", eliminarOpcion).Methods("DELETE")
	r.HandleFunc("/obtenerOpciones", obtenerOpciones).Methods("GET")

	r.HandleFunc("/obtenerHistorialPrecios", obtenerHistorialPrecios).Methods("GET")
	r.HandleFunc("/cambiarIntervaloBusqueda", cambiarIntervaloBusqueda).Methods("PUT")
	// r.HandleFunc("/agregarOferta", ).Methods("DELETE")
	http.Handle("/", r)
	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}
}
func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Landing page")
}

func agregarProducto(w http.ResponseWriter, r *http.Request) {
	var producto Product
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &producto)

	database.Create(&producto)

	json.NewEncoder(w).Encode(producto)

}

func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	var producto Product
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &producto)

	//eliminar repitable event to query pproduct details
	priceTarackerService.Stop(producto.ProductID)

	//eliminarProducto de la base de datos
	database.Delete(&producto)
	response := messageReponse{"ok", "ok"}

	json.NewEncoder(w).Encode(response)

}
func obtenerProductos(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	result := database.Find(&products)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar productos"})
		return
	}

	json.NewEncoder(w).Encode(products)

}

func agregarOpcion(w http.ResponseWriter, r *http.Request) {
	var newOption agregarOpcionRequestObject

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newOption)

	UpdateProductInformationFromInternet(&newOption.Option)
	fmt.Println("newOption updated: ", &newOption.Option)

	//query porduct
	var product Product
	result := database.First(&product, newOption.ProductID)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Producto"})
		return
	}
	//add the productOption to the list
	product.Options = append(product.Options, newOption.Option)
	//update the database
	database.Save(product)

	priceTarackerService.RegisterRepitableEvent(product.ProductID, makeExtractPriceCallBack(product), time.Duration(product.TrackingInterval*1000000000))

	json.NewEncoder(w).Encode(newOption.Option)

}

func obtenerOpciones(w http.ResponseWriter, r *http.Request) {
	var producto Product

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &producto)

	result := database.Preload("Options").Find(&producto)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Producto"})
		return
	}

	json.NewEncoder(w).Encode(producto.Options)

}

func eliminarOpcion(w http.ResponseWriter, r *http.Request) {
	var opcion ProductOption

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &opcion)

	result := database.Delete(&opcion)
	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Opcion de product"})
		return
	}
	//dejar de consultar informacion sobre el producto online
	producto := Product{
		ProductID: opcion.ProductID,
	}
	//obtener opciones de la base de datos reconstruir la llamada repetitiva
	result := database.Preload("Options").Find(&producto)
	priceTarackerService.RegisterRepitableEvent(producto.ProductID, makeExtractPriceCallBack(producto), time.Duration(producto.TrackingInterval*1000000000))

	json.NewEncoder(w).Encode(messageReponse{
		"Ok", "Ok"})
}
func obtenerHistorialPrecios(w http.ResponseWriter, r *http.Request) {
	var opcion ProductOption

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &opcion)

	result := database.Preload("Prices").Find(&opcion)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Opcion"})
		return
	}

	json.NewEncoder(w).Encode(opcion.Prices)

}
func cambiarIntervaloBusqueda(w http.ResponseWriter, r *http.Request) {
	var productoNuevo Product
	var producto Product

	// read json to get the new state of producto
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &productoNuevo)

	// read existing product from datasabe
	result := database.Find(&producto, productoNuevo.ProductID)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Producto"})
		return
	}
	// copy TrackingInterval field from productoNuevo
	producto.TrackingInterval = productoNuevo.TrackingInterval
	//update the datasabe
	database.Save(&producto)

	result := database.Preload("Options").Find(&producto)
	priceTarackerService.RegisterRepitableEvent(producto.ProductID, makeExtractPriceCallBack(producto), time.Duration(producto.TrackingInterval*1000000000))

	response := messageReponse{"ok", "ok"}
	json.NewEncoder(w).Encode(response)

}

func startTrackingPrices() {

	var products []Product
	database.Preload("Options.Prices").Preload("Options").Find(&products)

	for _, product := range products {
		fmt.Println("name: ", product.ProductName, "Duration: ", product.TrackingInterval)
		timeDuration := time.Duration(product.TrackingInterval * 1000000000)
		callback := makeExtractPriceCallBack(product)
		priceTarackerService.RegisterRepitableEvent(product.ProductID, callback, timeDuration)
	}
	fmt.Println(priceTarackerService.EventLookUp)

}
func makeExtractPriceCallBack(product Product) func() {
	return func() {
		for _, option := range product.Options {
			UpdateProductInformationFromInternet(&option)
			database.Save(&option)

			fmt.Println("query from internet for option with id: ", option.OptionID)
		}
	}
}
