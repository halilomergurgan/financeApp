package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	dsn := "financeuser:password@tcp(127.0.0.1:3306)/financeApp"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	router.HandleFunc("/categories", createCategory).Methods("POST")
	router.HandleFunc("/categories", getCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", getCategoryByID).Methods("GET")
	router.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := map[string]string{"error": "Invalid request payload"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		response := map[string]string{"error": "Username, email, and password are required"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		response := map[string]string{"error": "Database error"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if exists {
		response := map[string]string{"error": "Email address already in use"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	stmt, err := db.Prepare("INSERT INTO users(username, email, password) VALUES(?, ?, ?)")
	if err != nil {
		response := map[string]string{"error": "Database error"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Username, user.Email, user.Password)
	if err != nil {
		response := map[string]string{"error": "Failed to create user"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve user ID"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	user.ID = int(id)

	response := map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"users": users,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var user User
	err := db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"user": user,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" {
		http.Error(w, "Username and email are required", http.StatusBadRequest)
		return
	}

	user.ID = id

	stmt, err := db.Prepare("UPDATE users SET username = ?, email = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, id)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO categories(name) VALUES(?)")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(category.Name)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve category ID", http.StatusInternalServerError)
		return
	}

	category.ID = int(id)

	response := map[string]interface{}{
		"message":  "Category created successfully",
		"category": category,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			http.Error(w, "Error scanning category", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data": map[string]interface{}{
			"categories": categories,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category Category
	err = db.QueryRow("SELECT id, name FROM categories WHERE id = ?", id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"category": category,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var category Category
	err = json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if category.Name == "" {
		http.Error(w, "Category name is required", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE categories SET name = ? WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.Name, id)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	category.ID = id

	response := map[string]interface{}{
		"message":  "Category updated successfully",
		"category": category,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("DELETE FROM categories WHERE id = ?")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Category deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
