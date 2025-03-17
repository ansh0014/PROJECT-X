package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Offer model
type Offer struct {
	ServiceName  string  `json:"ServiceName"`
	Price        float64 `json:"Price"`
	Offer        string  `json:"Offer"`
	DeliveryTime int     `json:"DeliveryTime"` // in minutes
}

// Restaurant represents a food establishment
type Restaurant struct {
	Name     string
	Location string
}

// Delivery service data structure
var deliveryServices = map[string][]Offer{
	"mcdonalds:new york:ny": {
		{ServiceName: "Zomato", Price: 250.00, Offer: "20% off", DeliveryTime: 30},
		{ServiceName: "Swiggy", Price: 240.00, Offer: "15% off", DeliveryTime: 25},
		{ServiceName: "UberEats", Price: 260.00, Offer: "Free delivery", DeliveryTime: 28},
	},
	"burger king:new york:ny": {
		{ServiceName: "Zomato", Price: 220.00, Offer: "10% off", DeliveryTime: 35},
		{ServiceName: "Swiggy", Price: 215.00, Offer: "₹50 off", DeliveryTime: 40},
		{ServiceName: "UberEats", Price: 230.00, Offer: "Buy 1 Get 1", DeliveryTime: 32},
	},
	"dominos:new york:ny": {
		{ServiceName: "Zomato", Price: 300.00, Offer: "30% off", DeliveryTime: 25},
		{ServiceName: "Swiggy", Price: 320.00, Offer: "Free drink", DeliveryTime: 28},
		{ServiceName: "UberEats", Price: 290.00, Offer: "₹100 off", DeliveryTime: 35},
	},
	"pizza hut:chicago:il": {
		{ServiceName: "Zomato", Price: 280.00, Offer: "25% off", DeliveryTime: 40},
		{ServiceName: "Swiggy", Price: 275.00, Offer: "Free sides", DeliveryTime: 35},
		{ServiceName: "DoorDash", Price: 290.00, Offer: "10% cashback", DeliveryTime: 30},
	},
	"kfc:chicago:il": {
		{ServiceName: "Zomato", Price: 230.00, Offer: "15% off", DeliveryTime: 28},
		{ServiceName: "Swiggy", Price: 220.00, Offer: "₹70 off", DeliveryTime: 25},
		{ServiceName: "GrubHub", Price: 240.00, Offer: "Free dessert", DeliveryTime: 30},
	},
}

func main() {
	fmt.Println("Starting Food Delivery Comparator API")
	r := mux.NewRouter()

	// API Routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/offers", getOffers).Methods("GET")
	api.HandleFunc("/offers/{restaurant}/{city}/{state}", getSpecificOffers).Methods("GET")

	// Serve static files (frontend)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend"))))

	// Start server
	fmt.Println("Server is running on http://0.0.0.0:8000")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", r))
}

// Get all sample offers (for demo purposes)
func getOffers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get query parameters
	restaurant := strings.ToLower(r.URL.Query().Get("restaurant"))
	city := strings.ToLower(r.URL.Query().Get("city"))
	state := strings.ToLower(r.URL.Query().Get("state"))

	// If all parameters are provided, search for specific offers
	if restaurant != "" && city != "" && state != "" {
		key := fmt.Sprintf("%s:%s:%s", restaurant, city, state)
		if offers, exists := deliveryServices[key]; exists {
			json.NewEncoder(w).Encode(offers)
			return
		}
	}

	// If no matching offers or parameters missing, return sample data
	sampleOffers := []Offer{
		{ServiceName: "Zomato", Price: 250.00, Offer: "20% off", DeliveryTime: 30},
		{ServiceName: "Swiggy", Price: 240.00, Offer: "15% off", DeliveryTime: 25},
		{ServiceName: "UberEats", Price: 260.00, Offer: "Free delivery", DeliveryTime: 28},
	}
	json.NewEncoder(w).Encode(sampleOffers)
}

// Get specific offers for a restaurant in a location
func getSpecificOffers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	restaurant := strings.ToLower(vars["restaurant"])
	city := strings.ToLower(vars["city"])
	state := strings.ToLower(vars["state"])

	key := fmt.Sprintf("%s:%s:%s", restaurant, city, state)
	
	if offers, exists := deliveryServices[key]; exists {
		json.NewEncoder(w).Encode(offers)
		return
	}

	// If no matching offers, return empty array with proper status
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": fmt.Sprintf("No offers found for %s in %s, %s", restaurant, city, state),
	})
}
