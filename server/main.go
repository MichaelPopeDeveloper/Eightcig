package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Employee struct {
	ID          int
	Name        string
	Performance int
	Date        time.Time
}

var db *sql.DB

func main() {
	var err error
	db, err = getDBConnection()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(enableCors)

	r.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sortBy")
		employees, err := queryEmployees(sortBy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respondJSON(w, employees)
	}).Methods("GET")
	r.HandleFunc("/", rootHandler)

	// Start the server on port 8080
	fmt.Println("Server listening on http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func queryEmployees(sortBy string) ([]Employee, error) {
	query := "SELECT id, name, performance, date FROM employees"
	if sortBy == "performance" {
		query += " ORDER BY performance DESC"
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		var dateStr string
		err := rows.Scan(&emp.ID, &emp.Name, &emp.Performance, &dateStr)
		if err != nil {
			return nil, err
		}

		emp.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}

		employees = append(employees, emp)
	}

	return employees, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy") // Get the sortBy query parameter
	employees, err := queryEmployees(sortBy)
	if err != nil {
		http.Error(w, "Error fetching employees: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal employees data to JSON and send it in the response
	jsonData, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func getDBConnection() (*sql.DB, error) {
	dbUser := "root"
	dbPass := "password"
	dbHost := "eightcig.cprqn4l9e060.us-west-2.rds.amazonaws.com"
	dbPort := "3306"
	dbName := "eightcig"

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CORS middleware
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Update with your React app's URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
