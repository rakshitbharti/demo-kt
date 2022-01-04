package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CustomerController struct {
	Store CustomerStore
}
type Customer struct {
	ID, Name, Email string
}

type CustomerStore interface {
	Create(Customer) error
	Update(string, Customer) error
	Delete(string) error
	GetById(string) (Customer, error)
	GetAll() ([]Customer, error)
}

func (c *CustomerController) Add(w http.ResponseWriter, r *http.Request) {
	var customer1 Customer
	err1 := json.NewDecoder(r.Body).Decode(&customer1)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err1.Error(), http.StatusBadRequest)
	}
	err := c.Store.Create(customer1)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(customer1.ID + ": exists"))
		return
	}
	w.Write([]byte(customer1.ID + ": Created"))
	fmt.Println(customer1.ID, ":Customer has been created")
}
func (c *CustomerController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err1 := c.Store.Update(key, customer)
	if err1 != nil {
		fmt.Println("Error:", err1)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(key + ": does not Exists"))
		return
	}
	w.Write([]byte(key + ": updated"))
	fmt.Println(key, ":Customer Updated")
}
func (c *CustomerController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	err := c.Store.Delete(key)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(key + ": doesn't Exists"))
		return
	}
	w.Write([]byte(key + ": deleted"))
	fmt.Println(key, ":Customer Deleted")
}
func (c *CustomerController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	cust, err := c.Store.GetById(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(key + ": does not Exist"))
		fmt.Println("Error:", err)
		return
	}
	converted, err := json.Marshal(cust)
	if err != nil {
		w.Write([]byte("Error"))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(converted)
	fmt.Println(cust)
}
func (c *CustomerController) GetAll(w http.ResponseWriter, r *http.Request) {
	customers, err := c.Store.GetAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	converted, err := json.Marshal(customers)
	if err != nil {
		w.Write([]byte("Error"))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(converted)
	fmt.Println(customers)
}

type MapStore struct {
	Store map[string]Customer
}

func NewMapStore() *MapStore {
	return &MapStore{Store: make(map[string]Customer)}
}

func (m *MapStore) Create(c Customer) error {
	if _, ok := m.Store[c.ID]; ok {
		return errors.New("customer already exists")
	}
	m.Store[c.ID] = c
	return nil

}
func (m *MapStore) Update(s string, c Customer) error {
	if _, ok := m.Store[s]; ok {
		m.Store[s] = c
	} else {
		return errors.New("Customer does not exists")
	}
	return nil
}
func (m *MapStore) Delete(s string) error {
	if _, ok := m.Store[s]; ok {
		delete(m.Store, s)
	} else {
		return errors.New("Customer doesn't exists")
	}
	return nil
}
func (m *MapStore) GetById(s string) (Customer, error) {
	if val, ok := m.Store[s]; ok {
		return val, nil
	} else {
		return Customer{}, errors.New("No such customer exists")
	}
}
func (m *MapStore) GetAll() ([]Customer, error) {
	var Customers []Customer
	for _, v := range m.Store {
		Customers = append(Customers, v)
	}
	return Customers, nil
}
func InitializeRoutes(h *CustomerController) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/customer", h.GetAll).Methods("GET")
	r.HandleFunc("/api/customer/{id}", h.Get).Methods("GET")
	r.HandleFunc("/api/customer", h.Add).Methods("POST")
	r.HandleFunc("/api/customer/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/api/customer/{id}", h.Delete).Methods("DELETE")
	return r
}
func main() {
	controller := &CustomerController{
		Store: NewMapStore(),
	}
	r := InitializeRoutes(controller)
	log.Println("Listening...")
	//b := os.Getenv("PORT")
	//s := ":" + b
	http.ListenAndServe(":8080", r)
}
