package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     int64  `json:"phone" len:"10"`
	Contacted bool   `json:"contacted"`
}

type CustomerRequest struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     int64  `json:"phone" len:"10"`
	Contacted bool   `json:"contacted"`
}

type Error struct {
	Message string `json:"message"`
}

type Message struct {
	Message string `json:"message"`
}

var customers = []Customer{
	{"9f7d3f64-5f06-4e28-9d34-ff4f2e5d0db1", "John", "Customer", "john@abc.com", 9234562341, false},
	{"c8b7f0e3-3476-4b21-a2c4-3d5b5f5b8a6a", "Paul", "Customer", "paul@pqr.com", 9231928201, true},
	{"2b4c1fb5-51cb-4ff7-8d4e-9b7df673f7a2", "Max", "Customer", "max@xyz.com", 9241029203, false},
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "intro")
}

func getCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	if id == "" {
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

	var customerReq CustomerRequest

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &customerReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid request body"})
		return
	}

	id := uuid.New().String()

	var customer = Customer{id, customerReq.Name, customerReq.Role, customerReq.Email, customerReq.Phone, customerReq.Contacted}
	customers = append(customers, customer)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(customer)
	return
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	customerId := mux.Vars(r)["id"]

	if customerId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid customer id"})
		return
	}

	var customerReq CustomerRequest

	reqBody, _ := ioutil.ReadAll(r.Body)
	errBody := json.Unmarshal(reqBody, &customerReq)
	if errBody != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid request body"})
		return
	}

	for i, customer := range customers {
		if customer.ID == customerId {
			updatedCustomer := Customer{
				ID:        customer.ID,
				Name:      customerReq.Name,
				Role:      customerReq.Role,
				Email:     customerReq.Email,
				Phone:     customerReq.Phone,
				Contacted: customerReq.Contacted,
			}

			customers[i] = updatedCustomer
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(updatedCustomer)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(Error{"Customer not found"})
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Error{"Invalid customer id"})
		return
	}

	for i, customer := range customers {
		if customer.ID == id {
			customers = slices.Delete(customers, i, i+1)
			w.WriteHeader(http.StatusOK)
			message := fmt.Sprintf("Customer %s is deleted", id)
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
	r.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers", addCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	fmt.Println("Server is running on port 3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}
