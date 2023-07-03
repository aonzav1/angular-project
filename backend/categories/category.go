package category

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
	Type    string     `json:"type"`
	Data    []Category `json:"data"`
	Message string     `json:"message"`
}

type Category struct {
	CategoryID int    `json:"categoryid"`
	Name       string `json:"name"`
}

var _db *sql.DB

// Get all categorys

// response and request handlers
func GetCategories(w http.ResponseWriter, r *http.Request) {
	_db := db.SetupDB()

	utils.PrintMessage("Getting categories...")

	// Get all categorys from categorys table that don't have categoryID = "1"
	rows, err := _db.Query("SELECT * FROM categories")

	// check errors
	utils.CheckErr(err)

	// var response []JsonResponse
	var categorys []Category

	// Foreach category
	for rows.Next() {
		var categoryID int
		var name string

		err = rows.Scan(&categoryID, &name)

		// check errors
		utils.CheckErr(err)
		categorys = append(categorys, Category{CategoryID: categoryID, Name: name})
	}

	var response = JsonResponse{Type: "success", Data: categorys}

	json.NewEncoder(w).Encode(response)
}

// Create a category

// response and request handlers
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	var response = JsonResponse{}

	if name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing some parameter. require (name)"}
	} else {
		_db := db.SetupDB()

		utils.PrintMessage("Inserting category into DB")

		fmt.Println("Inserting new category name: " + name)

		var lastInsertID int
		err := _db.QueryRow("INSERT INTO categories(name) VALUES($1) returning categoryid;", name).Scan(&lastInsertID)

		// check errors
		utils.CheckErr(err)

		response = JsonResponse{Type: "success", Message: "The category has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a category

// response and request handlers
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	categoryID := params["categoryid"]

	var response = JsonResponse{}

	if categoryID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing categoryID parameter."}
	} else {
		_db := db.SetupDB()

		utils.PrintMessage("Deleting category from DB")

		_, err := _db.Exec("DELETE FROM categories where categoryID = $1", categoryID)

		// check errors
		utils.CheckErr(err)

		response = JsonResponse{Type: "success", Message: "The category has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all categorys

// response and request handlers
func DeleteCategories(w http.ResponseWriter, r *http.Request) {
	_db := db.SetupDB()

	utils.PrintMessage("Deleting all categorys...")

	_, err := _db.Exec("DELETE FROM categories")

	// check errors
	utils.CheckErr(err)

	utils.PrintMessage("All categorys have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All categorys have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
