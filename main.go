package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
)

type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

type Error struct {
	Message string `json:"message"`
}

type Message struct {
	Message string `json:"message"`
}

var customers = []Customer{
	{1, "John", "Customer", "john@abc.com", "9234562341", false},
	{2, "Paul", "Customer", "paul@pqr.com", "9231928201", true},
	{3, "Max", "Customer", "max@xyz.com", "9241029203", false},
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "intro")
}

func getCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid ID"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid ID"})
		return
	}

	for _, customer := range customers {
		if customer.ID == id {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(customer)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(Error{"Customer not found"})
	return
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(customers)
	return
}

func addCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var customer Customer

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid request body"})
		return
	}

	for _, existingCustomer := range customers {
		if existingCustomer.ID == customer.ID {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(Error{"Customer already exists"})
			return
		}
	}

	customers = append(customers, customer)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(customer)
	return
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	customerId, errID := strconv.Atoi(mux.Vars(r)["id"])

	if errID != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid customer id"})
		return
	}

	var customer Customer

	reqBody, _ := ioutil.ReadAll(r.Body)
	errBody := json.Unmarshal(reqBody, &customer)
	if errBody != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid request body"})
		return
	}

	for i, customer := range customers {
		if customer.ID == customerId {
			customers[i] = customer
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(customer)
			return
		}
	}

}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid customer id"})
		return
	}

	customerID, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid customer id"})
		return
	}

	for i, customer := range customers {
		if customer.ID == customerID {
			customers = slices.Delete(customers, i, i+1)
			w.WriteHeader(http.StatusOK)
			message := fmt.Sprintf("Customer %d is deleted", customerID)
			_ = json.NewEncoder(w).Encode(Message{message})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(Error{"customer is not found"})
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/customers/{id:[0-9]+}", getCustomer).Methods("GET")
	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers", addCustomer).Methods("POST")
	r.HandleFunc("/customers/{id:[0-9]+}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id:[0-9]+}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server is running on port 3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
