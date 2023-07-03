package product

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/m/db"
	"example.com/m/utils"
	"github.com/gorilla/mux"
)

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Product `json:"data"`
	Message string    `json:"message"`
}

type Product struct {
	ProductID   int     `json:"productid"`
	Name        string  `json:"name"`
	ImageUrl    string  `json:"imageUrl"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
}

var _db *sql.DB

// Get all products

// response and request handlers
func GetProducts(w http.ResponseWriter, r *http.Request) {
	_db := db.SetupDB()

	utils.PrintMessage("Getting products...")

	rows, err := _db.Query("SELECT productid, name, imageUrl, description, price, stock, categoryid FROM products")

	// check errors
	utils.CheckErr(err)

	//Should be able to be query by category

	// var response []JsonResponse
	var products []Product

	// Foreach product
	for rows.Next() {
		var productID, stock, category int
		var name, description, imageUrl string
		var price float32
		err = rows.Scan(&productID, &name, &imageUrl, &description, &price, &stock, &category)

		// check errors
		utils.CheckErr(err)
		products = append(products, Product{ProductID: productID, Name: name, Description: description, Price: price, Stock: stock, ImageUrl: imageUrl})
	}

	var response = JsonResponse{Type: "success", Data: products}

	json.NewEncoder(w).Encode(response)
}

// Create a product

// response and request handlers
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	description := r.FormValue("description")
	stock := r.FormValue("stock")
	price := r.FormValue("price")
	imageUrl := r.FormValue("imageUrl")
	categoryid := r.FormValue("categoryid")

	var response = JsonResponse{}

	if name == "" || description == "" || stock == "" || price == "" || imageUrl == "" || categoryid == "" {
		response = JsonResponse{Type: "error", Message: "You are missing some parameter. require (name, description, price, stock, imageUrl, categoryid)"}
	} else {
		_db := db.SetupDB()

		utils.PrintMessage("Inserting product into DB")

		fmt.Println("Inserting new product name: " + name)

		var lastInsertID int
		err := _db.QueryRow("INSERT INTO products(name, description, price, stock, imageUrl,categoryid) VALUES($1, $2, $3, $4, $5, $6) returning productid;", name, description, price, stock, imageUrl, categoryid).Scan(&lastInsertID)

		// check errors
		utils.CheckErr(err)

		response = JsonResponse{Type: "success", Message: "The product has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a product

// response and request handlers
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	productID := params["productid"]

	var response = JsonResponse{}

	if productID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing productID parameter."}
	} else {
		_db := db.SetupDB()

		utils.PrintMessage("Deleting product from DB")

		_, err := _db.Exec("DELETE FROM products where productID = $1", productID)

		// check errors
		utils.CheckErr(err)

		response = JsonResponse{Type: "success", Message: "The product has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all products

// response and request handlers
func DeleteProducts(w http.ResponseWriter, r *http.Request) {
	_db := db.SetupDB()

	utils.PrintMessage("Deleting all products...")

	_, err := _db.Exec("DELETE FROM products")

	// check errors
	utils.CheckErr(err)

	utils.PrintMessage("All products have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All products have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
