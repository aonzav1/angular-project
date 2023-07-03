package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/m/config"
	"example.com/m/db"
	"example.com/m/utils"
	"github.com/dgrijalva/jwt-go"
)

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []byte `json:"data"`
	Message string `json:"message"`
}

var _db *sql.DB

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}

	// Get username and password from the form
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	email := r.Form.Get("email")

	var response = JsonResponse{}

	_db = db.SetupDB()

	if !canRegister(username, email, password) {
		w.WriteHeader(http.StatusBadRequest)
		response = JsonResponse{Type: "error", Message: "Invalid username or password."}
	} else {
		hashedPassword, pwdErr := HashPassword(password)
		if pwdErr != nil {
			fmt.Println("Hash error")
		}

		utils.PrintMessage("Inserting user into DB")

		fmt.Println("Inserting new user with ID: " + username + " and hashed password: " + hashedPassword)

		var lastInsertID int
		err := _db.QueryRow("INSERT INTO users(username, password, email) VALUES($1, $2, $3) returning userid;", username, hashedPassword, email).Scan(&lastInsertID)

		// check errors
		utils.CheckErr(err)

		response = JsonResponse{Type: "success", Message: "User has been created successfully!"}
		fmt.Fprint(w, "User created successfully")
	}

	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}

	_db = db.SetupDB()
	// Get username and password from the form
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	var response = JsonResponse{}

	userid, loginErr := doLogin(username, password)
	if loginErr != nil || userid == -1 {
		response = JsonResponse{Type: "error", Message: "Invalid username or password."}
		fmt.Fprint(w, "Invalid credentials")
	} else {
		// Create a new JWT token
		fmt.Fprint(w, "Login as userId "+string(userid))
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Userid:   userid,
			Username: username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(config.GetJwtKey())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to generate token")
			return
		}
		response = JsonResponse{Type: "success", Data: []byte(tokenString), Message: "Logged in!"}
		fmt.Fprint(w, "Login success")
	}

	json.NewEncoder(w).Encode(response)
}

func canRegister(username string, email string, password string) bool {
	utils.PrintMessage(username + ", " + password + ", " + email)
	if username == "" || password == "" || email == "" {
		return false
	}
	var userid int
	err := _db.QueryRow("SELECT userid FROM users where username = $1 OR email = $2", username, email).Scan(&userid)

	if err == nil {
		utils.PrintMessage("Repeat username or email, detected")
		return false
	}

	//Check security level of password
	return true
}

func doLogin(username string, password string) (int, error) {
	var storedPassword string = ""
	var userid int = 0
	err := _db.QueryRow("SELECT userid,password FROM users where username = $1", username).Scan(&userid, &storedPassword)

	if storedPassword == "" {
		utils.PrintMessage("Can't find " + username)
	} else {
		utils.PrintMessage("Compare " + password + " with " + storedPassword)
		if ComparePassword(password, storedPassword) {
			return userid, err
		}
	}
	return -1, err
}
