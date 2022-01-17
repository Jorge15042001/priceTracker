package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"io/ioutil"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var datasabe *gorm.DB

func main() {
	datasabe = getDataBaseConnection() //stablish connection with the datasabe

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

	datasabe.Create(&producto)

	json.NewEncoder(w).Encode(producto)

}

func eliminarProducto(w http.ResponseWriter, r *http.Request) {
	var producto Product
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &producto)

	datasabe.Delete(&producto)
	response := messageReponse{"ok", "ok"}

	json.NewEncoder(w).Encode(response)

}
func obtenerProductos(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	result := datasabe.Find(&products)

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
	result := datasabe.First(&product, newOption.ProductID)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Producto"})
		return
	}
	//add the productOption to the list
	product.Options = append(product.Options, newOption.Option)
	//update the database
	datasabe.Save(product)

	json.NewEncoder(w).Encode(newOption.Option)

}

func obtenerOpciones(w http.ResponseWriter, r *http.Request) {
	var producto Product

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &producto)

	result := datasabe.Preload("Options").Find(&producto)

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

	result := datasabe.Delete(&opcion)
	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Opcion de product"})
		return
	}

	json.NewEncoder(w).Encode(messageReponse{
		"Ok", "Ok"})
}
func obtenerHistorialPrecios(w http.ResponseWriter, r *http.Request) {
	var opcion ProductOption

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &opcion)

	result := datasabe.Preload("Prices").Find(&opcion)

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
	result := datasabe.Find(&producto, productoNuevo.ProductID)

	if result.Error != nil {
		json.NewEncoder(w).Encode(messageReponse{
			"Error", "No es posible encontrar Producto"})
		return
	}
	// copy TrackingInterval field from productoNuevo
	producto.TrackingInterval = productoNuevo.TrackingInterval
	//update the datasabe
	datasabe.Save(&producto)

	response := messageReponse{"ok", "ok"}
	json.NewEncoder(w).Encode(response)

}
