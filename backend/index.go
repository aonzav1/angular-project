package main

import (
	"fmt"
	"log"
	"net/http"

	category "example.com/m/categories"
	product "example.com/m/products"
	user "example.com/m/users"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Main function
func main() {

	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling root request")
	})

	router.HandleFunc("/register", user.RegisterHandler).Methods("Post")
	router.HandleFunc("/login", user.LoginHandler).Methods("Post")

	router.HandleFunc("/products/", product.GetProducts).Methods("GET")
	router.HandleFunc("/products/add", product.CreateProduct).Methods("POST")
	router.HandleFunc("/products/delete/{productid}", product.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/products/delete/", product.DeleteProducts).Methods("DELETE")

	router.HandleFunc("/categories/", category.GetCategories).Methods("GET")
	router.HandleFunc("/categories/add", category.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/delete/{productid}", category.DeleteCategory).Methods("DELETE")
	router.HandleFunc("/categories/delete/", category.DeleteCategories).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
