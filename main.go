package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ServiceOffer struct {
	ServiceName  string  `json:"ServiceName"`
	Price        float64 `json:"Price"`
	Offer        string  `json:"Offer"`
	DeliveryTime int     `json:"DeliveryTime,omitempty"`
	Duration     int     `json:"Duration,omitempty"`
}

// Available categories
const (
	CategoryTaxi          = "taxi"
	CategoryRestaurant    = "restaurant"
	CategoryQuickCommerce = "quickcommerce"
)

var taxiServices = map[string][]ServiceOffer{
	"india:delhi:india:mumbai": {
		{ServiceName: "Uber", Price: 6500.00, Offer: "10% cashback", Duration: 1260},
		{ServiceName: "Ola", Price: 7000.00, Offer: "Free waiting", Duration: 1200},
	},
	"india:punjab:india:himachal pradesh": {
		{ServiceName: "Uber", Price: 1700.00, Offer: "₹100 off", Duration: 240},
		{ServiceName: "Ola", Price: 1600.00, Offer: "20% off first ride", Duration: 210},
	},
}

var restaurantServices = map[string][]ServiceOffer{
	"india:punjab:patiala:dominos": {
		{ServiceName: "Zomato", Price: 350.00, Offer: "20% off", DeliveryTime: 30},
		{ServiceName: "Swiggy", Price: 320.00, Offer: "Free drink", DeliveryTime: 25},
	},
	"india:delhi:delhi:burger king": {
		{ServiceName: "Zomato", Price: 250.00, Offer: "30% off", DeliveryTime: 35},
		{ServiceName: "Swiggy", Price: 240.00, Offer: "₹50 off", DeliveryTime: 40},
	},
}

var quickCommerceServices = map[string][]ServiceOffer{
	"india:punjab:patiala:thapar university": {
		{ServiceName: "Zepto", Price: 120.00, Offer: "Free delivery", DeliveryTime: 10},
		{ServiceName: "Blinkit", Price: 110.00, Offer: "15% off", DeliveryTime: 12},
	},
	"india:delhi:delhi:india gate": {
		{ServiceName: "Zepto", Price: 150.00, Offer: "₹30 cashback", DeliveryTime: 15},
		{ServiceName: "Blinkit", Price: 140.00, Offer: "Buy 1 Get 1", DeliveryTime: 20},
	},
}

var locationOptions = map[string]interface{}{
	"countries": []string{"India"},
	"states": map[string][]string{
		"India": {
			"Andhra Pradesh", "Arunachal Pradesh", "Assam", "Bihar",
			"Chhattisgarh", "Delhi", "Goa", "Gujarat", "Haryana",
			"Himachal Pradesh", "Jharkhand", "Karnataka", "Kerala",
			"Madhya Pradesh", "Maharashtra", "Manipur", "Meghalaya",
			"Mizoram", "Nagaland", "Odisha", "Punjab", "Rajasthan",
			"Sikkim", "Tamil Nadu", "Telangana", "Tripura", "Uttar Pradesh",
			"Uttarakhand", "West Bengal",
		},
	},
	"cities": map[string]map[string][]string{
		"India": {
			"Andhra Pradesh":    {"Visakhapatnam", "Vijayawada", "Guntur", "Nellore", "Kurnool"},
			"Arunachal Pradesh": {"Itanagar", "Naharlagun", "Pasighat", "Tawang"},
			"Assam":             {"Guwahati", "Silchar", "Dibrugarh", "Jorhat", "Nagaon"},
			"Bihar":             {"Patna", "Gaya", "Muzaffarpur", "Bhagalpur", "Darbhanga"},
			"Chhattisgarh":      {"Raipur", "Bhilai", "Bilaspur", "Korba", "Durg"},
			"Delhi":             {"Delhi", "New Delhi", "Dwarka", "Rohini", "Pitampura"},
			"Goa":               {"Panaji", "Margao", "Vasco da Gama", "Mapusa", "Ponda"},
			"Gujarat":           {"Ahmedabad", "Surat", "Vadodara", "Rajkot", "Bhavnagar"},
			"Haryana":           {"Gurgaon", "Faridabad", "Hisar", "Panipat", "Ambala"},
			"Himachal Pradesh":  {"Shimla", "Dharamshala", "Manali", "Solan", "Kullu"},
			"Jharkhand":         {"Ranchi", "Jamshedpur", "Dhanbad", "Bokaro", "Hazaribagh"},
			"Karnataka":         {"Bangalore", "Mysore", "Hubli", "Mangalore", "Belgaum"},
			"Kerala":            {"Thiruvananthapuram", "Kochi", "Kozhikode", "Thrissur", "Kollam"},
			"Madhya Pradesh":    {"Indore", "Bhopal", "Jabalpur", "Gwalior", "Ujjain"},
			"Maharashtra":       {"Mumbai", "Pune", "Nagpur", "Thane", "Nashik"},
			"Manipur":           {"Imphal", "Thoubal", "Kakching", "Ukhrul", "Chandel"},
			"Meghalaya":         {"Shillong", "Tura", "Jowai", "Nongstoin", "Baghmara"},
			"Mizoram":           {"Aizawl", "Lunglei", "Champhai", "Saiha", "Kolasib"},
			"Nagaland":          {"Kohima", "Dimapur", "Mokokchung", "Tuensang", "Wokha"},
			"Odisha":            {"Bhubaneswar", "Cuttack", "Rourkela", "Berhampur", "Sambalpur"},
			"Punjab":            {"Ludhiana", "Amritsar", "Jalandhar", "Patiala", "Bathinda"},
			"Rajasthan":         {"Jaipur", "Jodhpur", "Udaipur", "Kota", "Ajmer"},
			"Sikkim":            {"Gangtok", "Namchi", "Mangan", "Gyalshing", "Rangpo"},
			"Tamil Nadu":        {"Chennai", "Coimbatore", "Madurai", "Tiruchirappalli", "Salem"},
			"Telangana":         {"Hyderabad", "Warangal", "Nizamabad", "Karimnagar", "Khammam"},
			"Tripura":           {"Agartala", "Udaipur", "Dharmanagar", "Kailashahar", "Belonia"},
			"Uttar Pradesh":     {"Lucknow", "Kanpur", "Agra", "Varanasi", "Meerut"},
			"Uttarakhand":       {"Dehradun", "Haridwar", "Roorkee", "Haldwani", "Rudrapur"},
			"West Bengal":       {"Kolkata", "Howrah", "Durgapur", "Asansol", "Siliguri"},
		},
	},
	"restaurants": map[string]map[string][]string{},
	"addresses":   map[string]map[string][]string{},
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for now
	},
}

type RealTimeRequest struct {
	Category    string `json:"category"`
	FromCountry string `json:"fromCountry,omitempty"`
	FromState   string `json:"fromState,omitempty"`
	ToCountry   string `json:"toCountry,omitempty"`
	ToState     string `json:"toState,omitempty"`
	Country     string `json:"country,omitempty"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	Restaurant  string `json:"restaurant,omitempty"`
	Address     string `json:"address,omitempty"`
	GroceryItem string `json:"groceryItem,omitempty"`
}

type RealTimeResponse struct {
	Category  string         `json:"category"`
	Route     string         `json:"route,omitempty"`
	Location  string         `json:"location,omitempty"`
	Offers    []ServiceOffer `json:"offers"`
	Timestamp int64          `json:"timestamp"`
}

type ClientSubscription struct {
	request RealTimeRequest
	conn    *websocket.Conn
}

var (
	clients       = make(map[*websocket.Conn]bool)
	subscriptions = make(map[*websocket.Conn]*ClientSubscription)
	clientsMutex  = sync.Mutex{}
	seed          = rand.NewSource(time.Now().UnixNano())
	rnd           = rand.New(seed)
)

func getDynamicRestaurantOptions() map[string]map[string][]string {
	options := make(map[string]map[string][]string)

	chains := []string{
		"Dominos", "Pizza Hut", "McDonald's", "Burger King",
		"KFC", "Subway", "Haldiram's", "Barbeque Nation",
		"Biryani Blues", "Wow! Momo", "Paradise Biryani",
		"Faasos", "Behrouz Biryani", "Truffles", "Theobroma",
	}

	states := locationOptions["states"].(map[string][]string)["India"]
	for _, state := range states {
		options[state] = make(map[string][]string)

		if cities, ok := locationOptions["cities"].(map[string]map[string][]string)["India"][state]; ok {
			for _, city := range cities {

				numRestaurants := 3 + rand.Intn(5)
				cityRestaurants := make([]string, 0, numRestaurants)

				selectedIndexes := make(map[int]bool)
				for i := 0; i < numRestaurants; i++ {
					idx := rand.Intn(len(chains))

					for selectedIndexes[idx] {
						idx = rand.Intn(len(chains))
					}
					selectedIndexes[idx] = true
					cityRestaurants = append(cityRestaurants, chains[idx])
				}

				options[state][city] = cityRestaurants
			}
		}
	}

	return options
}

// Generate dynamic address options for every city
func getDynamicAddressOptions() map[string]map[string][]string {
	options := make(map[string]map[string][]string)

	// Common address patterns across India
	addressPatterns := []string{
		"Main Market", "City Center", "Railway Station", "Airport",
		"Central Mall", "Bus Stand", "University Campus", "City Park",
		"District Hospital", "Tech Park", "Industrial Area", "Metro Station",
		"Stadium", "Government Complex", "Town Hall", "Central Library",
	}

	// For each state in India
	states := locationOptions["states"].(map[string][]string)["India"]
	for _, state := range states {
		options[state] = make(map[string][]string)

		// For each city in that state
		if cities, ok := locationOptions["cities"].(map[string]map[string][]string)["India"][state]; ok {
			for _, city := range cities {
				// Add 3-6 addresses for each city
				numAddresses := 3 + rand.Intn(4) // 3 to 6 addresses
				cityAddresses := make([]string, 0, numAddresses)

				// Select random addresses without duplicates
				selectedIndexes := make(map[int]bool)
				for i := 0; i < numAddresses; i++ {
					idx := rand.Intn(len(addressPatterns))
					// Avoid duplicates
					for selectedIndexes[idx] {
						idx = rand.Intn(len(addressPatterns))
					}
					selectedIndexes[idx] = true
					cityAddresses = append(cityAddresses, addressPatterns[idx])
				}

				options[state][city] = cityAddresses
			}
		}
	}

	return options
}

// Function to generate grocery items options
func getDynamicGroceryOptions() map[string][]string {
	groceryOptions := make(map[string][]string)

	// Common grocery items available in India
	groceryItems := []string{
		"Rice (5kg)", "Wheat Flour (1kg)", "Toor Dal (1kg)", "Cooking Oil (1L)",
		"Sugar (1kg)", "Salt (1kg)", "Milk (1L)", "Bread (400g)",
		"Eggs (12)", "Potatoes (1kg)", "Onions (1kg)", "Tomatoes (1kg)",
		"Tea Leaves (250g)", "Coffee Powder (250g)", "Biscuits (Assorted)",
		"Breakfast Cereal", "Noodles Pack", "Spices Set", "Ghee (500g)",
		"Paneer (200g)", "Coconut Oil (500ml)", "Mustard Oil (1L)", "Honey (250g)",
		"Jam (300g)", "Sauce (200g)", "Curd (400g)", "Butter (100g)",
		"Cheese (200g)", "Fresh Fruits Pack", "Fresh Vegetables Pack",
	}

	// Add all grocery items to the options
	groceryOptions["items"] = groceryItems

	return groceryOptions
}

// Initialize the dynamic options for restaurants, addresses, and grocery items
func initializeDynamicOptions() {
	// Initialize restaurant options
	restaurants := getDynamicRestaurantOptions()
	locationOptions["restaurants"] = restaurants

	// Initialize address options
	addresses := getDynamicAddressOptions()
	locationOptions["addresses"] = addresses

	// Initialize grocery items
	groceryItems := getDynamicGroceryOptions()
	locationOptions["groceryItems"] = groceryItems
}

func main() {
	fmt.Println("Starting Multi-Service Price Comparator API")

	// Initialize dynamic location options
	initializeDynamicOptions()

	r := mux.NewRouter()

	// API Routes
	api := r.PathPrefix("/api").Subrouter()

	// Get options for form fields
	api.HandleFunc("/options", getOptions).Methods("GET")

	// Compare services by category
	api.HandleFunc("/compare/taxi", compareTaxi).Methods("GET")
	api.HandleFunc("/compare/restaurant", compareRestaurant).Methods("GET")
	api.HandleFunc("/compare/quickcommerce", compareQuickCommerce).Methods("GET")

	// WebSocket endpoint for real-time updates
	r.HandleFunc("/ws", handleWebSocket)

	// Serve static files (frontend)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend"))))

	// Start real-time price update goroutine
	go updatePricesRoutine()

	// Start server
	fmt.Println("Server is running on http://localhost:5000")
	fmt.Println("WebSocket server available at ws://localhost:5000/ws")
	log.Fatal(http.ListenAndServe("localhost:5000", r))
}

// Handle WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	log.Printf("New WebSocket connection established: %s", conn.RemoteAddr())

	// Remove client when connection closes
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		delete(subscriptions, conn)
		clientsMutex.Unlock()
		log.Printf("WebSocket connection closed: %s", conn.RemoteAddr())
	}()

	// Handle incoming messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process subscription request
		var request RealTimeRequest
		if err := json.Unmarshal(message, &request); err != nil {
			log.Printf("Error unmarshaling WebSocket message: %v", err)
			continue
		}

		// Register subscription
		clientsMutex.Lock()
		subscriptions[conn] = &ClientSubscription{
			request: request,
			conn:    conn,
		}
		clientsMutex.Unlock()

		// Send initial data immediately
		sendRealTimeResponse(conn, request)
	}
}

// Update prices routine
func updatePricesRoutine() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Apply small random fluctuations to prices
			applyPriceFluctuations()

			// Send updates to all clients
			clientsMutex.Lock()
			for conn, sub := range subscriptions {
				if conn.WriteMessage(websocket.PingMessage, nil) != nil {
					// Connection is dead
					delete(clients, conn)
					delete(subscriptions, conn)
					conn.Close()
					continue
				}

				// Send updated data
				sendRealTimeResponse(conn, sub.request)
			}
			clientsMutex.Unlock()
		}
	}
}

// Apply random price fluctuations to service offers (simulating real-time changes)
func applyPriceFluctuations() {
	// Apply to taxi services
	for key, offers := range taxiServices {
		for i := range offers {
			// Random fluctuation between -5% and +5%
			fluctuation := 1.0 + (rnd.Float64()*0.1 - 0.05)
			offers[i].Price = offers[i].Price * fluctuation
			// Round to 2 decimal places
			offers[i].Price = float64(int(offers[i].Price*100)) / 100
		}
		taxiServices[key] = offers
	}

	// Apply to restaurant services
	for key, offers := range restaurantServices {
		for i := range offers {
			fluctuation := 1.0 + (rnd.Float64()*0.1 - 0.05)
			offers[i].Price = offers[i].Price * fluctuation
			offers[i].Price = float64(int(offers[i].Price*100)) / 100
		}
		restaurantServices[key] = offers
	}

	// Apply to quick commerce services
	for key, offers := range quickCommerceServices {
		for i := range offers {
			fluctuation := 1.0 + (rnd.Float64()*0.1 - 0.05)
			offers[i].Price = offers[i].Price * fluctuation
			offers[i].Price = float64(int(offers[i].Price*100)) / 100
		}
		quickCommerceServices[key] = offers
	}
}

// Send real-time response to a specific client
func sendRealTimeResponse(conn *websocket.Conn, request RealTimeRequest) {
	var (
		offers   []ServiceOffer
		route    string
		location string
	)

	switch request.Category {
	case CategoryTaxi:
		key := fmt.Sprintf("%s:%s:%s:%s",
			strings.ToLower(request.FromCountry),
			strings.ToLower(request.FromState),
			strings.ToLower(request.ToCountry),
			strings.ToLower(request.ToState))

		// Check if we have pre-existing data
		if data, exists := taxiServices[key]; exists {
			offers = data
		} else {
			// Generate dynamic offers for any route
			offers = generateDynamicTaxiOffers(request.FromState, request.ToState)
			taxiServices[key] = offers
		}

		route = fmt.Sprintf("%s to %s", request.FromState, request.ToState)

	case CategoryRestaurant:
		key := fmt.Sprintf("%s:%s:%s:%s",
			strings.ToLower(request.Country),
			strings.ToLower(request.State),
			strings.ToLower(request.City),
			strings.ToLower(request.Restaurant))

		// Check if we have pre-existing data
		if data, exists := restaurantServices[key]; exists {
			offers = data
		} else {
			// Generate dynamic offers for any restaurant in any city
			offers = generateDynamicRestaurantOffers(request.Restaurant, request.City)
			restaurantServices[key] = offers
		}

		location = fmt.Sprintf("%s, %s", request.City, request.State)

	case CategoryQuickCommerce:
		var baseKey string
		if request.GroceryItem != "" {
			// If grocery item is specified, include it in the key
			baseKey = fmt.Sprintf("%s:%s:%s:%s:%s",
				strings.ToLower(request.Country),
				strings.ToLower(request.State),
				strings.ToLower(request.City),
				strings.ToLower(request.Address),
				strings.ToLower(request.GroceryItem))
		} else {
			// Otherwise use the standard key
			baseKey = fmt.Sprintf("%s:%s:%s:%s",
				strings.ToLower(request.Country),
				strings.ToLower(request.State),
				strings.ToLower(request.City),
				strings.ToLower(request.Address))
		}

		// Check if we have pre-existing data
		if data, exists := quickCommerceServices[baseKey]; exists {
			offers = data
		} else {
			// Generate dynamic offers
			if request.GroceryItem != "" {
				// Generate offers for specific grocery item
				offers = generateDynamicGroceryItemOffers(request.GroceryItem, request.Address)
			} else {
				// Generate general quick commerce offers
				offers = generateDynamicQuickCommerceOffers(request.Address, request.City)
			}
			quickCommerceServices[baseKey] = offers
		}

		location = fmt.Sprintf("%s, %s", request.City, request.State)
	}

	// Skip if no offers found
	if len(offers) == 0 {
		return
	}

	// Send response
	response := RealTimeResponse{
		Category:  request.Category,
		Route:     route,
		Location:  location,
		Offers:    offers,
		Timestamp: time.Now().Unix(),
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling WebSocket response: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
		log.Printf("Error sending WebSocket message: %v", err)
	}
}

// Get location options for form fields
func getOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := r.URL.Query().Get("category")
	country := r.URL.Query().Get("country")
	state := r.URL.Query().Get("state")
	city := r.URL.Query().Get("city")
	address := r.URL.Query().Get("address") // For getting grocery items

	var result interface{}

	// Return appropriate options based on the query parameters
	switch {
	case category == "":
		// Return available categories
		result = map[string][]string{
			"categories": {CategoryTaxi, CategoryRestaurant, CategoryQuickCommerce},
		}
	case country == "":
		// Return available countries
		result = map[string][]string{
			"countries": locationOptions["countries"].([]string),
		}
	case state == "":
		// Return available states for the country
		if states, ok := locationOptions["states"].(map[string][]string)[country]; ok {
			result = map[string][]string{
				"states": states,
			}
		} else {
			result = map[string][]string{
				"states": {},
			}
		}
	case city == "":
		// Return available cities for the state
		if cities, ok := locationOptions["cities"].(map[string]map[string][]string)[country][state]; ok {
			result = map[string][]string{
				"cities": cities,
			}
		} else {
			result = map[string][]string{
				"cities": {},
			}
		}
	default:
		// Return available restaurants or addresses based on category
		if category == CategoryRestaurant {
			if restaurants, ok := locationOptions["restaurants"].(map[string]map[string][]string)[state][city]; ok {
				result = map[string][]string{
					"restaurants": restaurants,
				}
			} else {
				result = map[string][]string{
					"restaurants": {},
				}
			}
		} else if category == CategoryQuickCommerce {
			// First return address options
			if address == "" {
				if addresses, ok := locationOptions["addresses"].(map[string]map[string][]string)[state][city]; ok {
					result = map[string][]string{
						"addresses": addresses,
					}
				} else {
					result = map[string][]string{
						"addresses": {},
					}
				}
			} else {
				// If address is specified, return grocery items
				if groceryItems, ok := locationOptions["groceryItems"].(map[string][]string)["items"]; ok {
					result = map[string][]string{
						"groceryItems": groceryItems,
					}
				} else {
					result = map[string][]string{
						"groceryItems": {},
					}
				}
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

// Compare taxi services
// Generate dynamic taxi offers between any two states
func generateDynamicTaxiOffers(fromState, toState string) []ServiceOffer {
	// Calculate base price based on state names
	// This creates a predictable but unique price for each route
	fromLen := len(fromState)
	toLen := len(toState)

	// Use state name length to create base pricing patterns
	basePrice := float64(fromLen*100 + toLen*120)

	// For intra-state travel, reduce price
	if strings.EqualFold(fromState, toState) {
		basePrice = float64(fromLen * 80)
	}

	// Adjust for popular routes
	popularStates := map[string]bool{
		"delhi": true, "maharashtra": true, "karnataka": true,
		"tamil nadu": true, "telangana": true, "west bengal": true,
	}

	if popularStates[strings.ToLower(fromState)] && popularStates[strings.ToLower(toState)] {
		basePrice *= 1.2 // Premium for routes between major states
	}

	// Calculate trip duration based on states
	// Assume 60 minutes per letter in state names (silly but deterministic)
	duration := (fromLen + toLen) * 60
	if duration < 120 {
		duration = 120 // Minimum 2 hours
	}

	// Different pricing and offers for different services
	uberPrice := basePrice * (1.0 + (rnd.Float64() * 0.1))
	olaPrice := basePrice * (0.95 + (rnd.Float64() * 0.1)) // Slightly cheaper on average

	// Round to 2 decimal places
	uberPrice = float64(int(uberPrice*100)) / 100
	olaPrice = float64(int(olaPrice*100)) / 100

	// Generate appropriate offers based on route
	uberOffer := "10% cashback"
	olaOffer := "Free waiting"

	// Different offers for different routes
	if strings.Contains(strings.ToLower(fromState), "a") {
		uberOffer = "₹100 off next ride"
	}
	if strings.Contains(strings.ToLower(toState), "i") {
		olaOffer = "20% off first ride"
	}

	return []ServiceOffer{
		{ServiceName: "Uber", Price: uberPrice, Offer: uberOffer, Duration: duration},
		{ServiceName: "Ola", Price: olaPrice, Offer: olaOffer, Duration: duration - 30}, // Slightly faster
	}
}

// Generate dynamic restaurant offers for any restaurant in any city
func generateDynamicRestaurantOffers(restaurant, city string) []ServiceOffer {
	// Base price depends on restaurant and city names
	restaurantLen := len(restaurant)
	cityLen := len(city)

	// Base price formula - creates unique but predictable prices
	basePrice := float64(restaurantLen*20 + cityLen*15)
	if basePrice < 150 {
		basePrice = 150 // Minimum price
	}
	if basePrice > 800 {
		basePrice = 800 // Maximum price
	}

	// Adjust for premium restaurants
	premiumRestaurants := map[string]bool{
		"barbeque nation": true, "theobroma": true, "mcdonald's": true,
		"pizza hut": true, "kfc": true, "wow! momo": true,
	}

	if premiumRestaurants[strings.ToLower(restaurant)] {
		basePrice *= 1.3 // Premium pricing
	}

	// Different pricing for different services
	zomatoPrice := basePrice * (1.0 + (rnd.Float64() * 0.1))
	swiggyPrice := basePrice * (0.95 + (rnd.Float64() * 0.1)) // Slightly cheaper on average

	// Round to 2 decimal places
	zomatoPrice = float64(int(zomatoPrice*100)) / 100
	swiggyPrice = float64(int(swiggyPrice*100)) / 100

	// Delivery times and offers
	zomatoDeliveryTime := 25 + rand.Intn(20) // 25-45 minutes
	swiggyDeliveryTime := 20 + rand.Intn(25) // 20-45 minutes

	// Generate appropriate offers
	zomatoOffer := "20% off"
	swiggyOffer := "Free delivery"

	// Different offers based on restaurant
	if strings.Contains(strings.ToLower(restaurant), "p") {
		zomatoOffer = "Buy 1 Get 1"
	}
	if strings.Contains(strings.ToLower(restaurant), "b") {
		swiggyOffer = "₹50 off"
	}

	return []ServiceOffer{
		{ServiceName: "Zomato", Price: zomatoPrice, Offer: zomatoOffer, DeliveryTime: zomatoDeliveryTime},
		{ServiceName: "Swiggy", Price: swiggyPrice, Offer: swiggyOffer, DeliveryTime: swiggyDeliveryTime},
	}
}

// Generate dynamic quick commerce offers for any address in any city
func generateDynamicQuickCommerceOffers(address, city string) []ServiceOffer {
	// Base price depends on address and city
	addressLen := len(address)
	cityLen := len(city)

	// Base price formula
	basePrice := float64(addressLen*3 + cityLen*5)
	if basePrice < 80 {
		basePrice = 80 // Minimum price
	}
	if basePrice > 250 {
		basePrice = 250 // Maximum price
	}

	// Adjust for busy locations
	busyLocations := map[string]bool{
		"railway station": true, "airport": true, "central mall": true,
		"main market": true, "metro station": true,
	}

	if busyLocations[strings.ToLower(address)] {
		basePrice *= 1.15 // Higher pricing for busy areas
	}

	// Different pricing for different services
	zeptoPrice := basePrice * (1.0 + (rnd.Float64() * 0.1))
	blinkitPrice := basePrice * (0.95 + (rnd.Float64() * 0.1)) // Slightly cheaper on average

	// Round to 2 decimal places
	zeptoPrice = float64(int(zeptoPrice*100)) / 100
	blinkitPrice = float64(int(blinkitPrice*100)) / 100

	// Delivery times and offers
	zeptoDeliveryTime := 10 + rand.Intn(10)  // 10-20 minutes
	blinkitDeliveryTime := 8 + rand.Intn(12) // 8-20 minutes

	// Generate appropriate offers
	zeptoOffer := "Free delivery"
	blinkitOffer := "15% off"

	// Different offers based on location
	if strings.Contains(strings.ToLower(address), "station") {
		zeptoOffer = "₹30 cashback"
	}
	if strings.Contains(strings.ToLower(address), "central") {
		blinkitOffer = "Buy 1 Get 1"
	}

	return []ServiceOffer{
		{ServiceName: "Zepto", Price: zeptoPrice, Offer: zeptoOffer, DeliveryTime: zeptoDeliveryTime},
		{ServiceName: "Blinkit", Price: blinkitPrice, Offer: blinkitOffer, DeliveryTime: blinkitDeliveryTime},
	}
}

// Generate dynamic offers for a specific grocery item
func generateDynamicGroceryItemOffers(groceryItem, address string) []ServiceOffer {
	// Base price depends on grocery item properties
	itemLen := len(groceryItem)

	// Generate a base price based on the item name
	var basePrice float64

	// Price categories based on grocery type
	// This mimics real-world pricing where some grocery categories are more expensive
	if strings.Contains(strings.ToLower(groceryItem), "rice") ||
		strings.Contains(strings.ToLower(groceryItem), "flour") ||
		strings.Contains(strings.ToLower(groceryItem), "dal") {
		// Staples
		basePrice = 80.0 + (float64(itemLen) * 2.5)
	} else if strings.Contains(strings.ToLower(groceryItem), "oil") ||
		strings.Contains(strings.ToLower(groceryItem), "ghee") {
		// Cooking oils
		basePrice = 120.0 + (float64(itemLen) * 3.5)
	} else if strings.Contains(strings.ToLower(groceryItem), "milk") ||
		strings.Contains(strings.ToLower(groceryItem), "bread") ||
		strings.Contains(strings.ToLower(groceryItem), "egg") {
		// Daily essentials
		basePrice = 50.0 + (float64(itemLen) * 1.5)
	} else if strings.Contains(strings.ToLower(groceryItem), "fruit") ||
		strings.Contains(strings.ToLower(groceryItem), "vegetable") {
		// Fresh produce
		basePrice = 100.0 + (float64(itemLen) * 2.0)
	} else {
		// Other grocery items
		basePrice = 70.0 + (float64(itemLen) * 2.0)
	}

	// Adjust for premium items
	if strings.Contains(strings.ToLower(groceryItem), "premium") ||
		strings.Contains(strings.ToLower(groceryItem), "organic") {
		basePrice *= 1.3 // Premium pricing
	}

	// Different pricing for different services
	zeptoPrice := basePrice * (1.0 + (rnd.Float64() * 0.1))
	blinkitPrice := basePrice * (0.95 + (rnd.Float64() * 0.1)) // Slightly cheaper on average

	// Round to 2 decimal places
	zeptoPrice = float64(int(zeptoPrice*100)) / 100
	blinkitPrice = float64(int(blinkitPrice*100)) / 100

	// Delivery times
	zeptoDeliveryTime := 10 + rand.Intn(5)  // 10-15 minutes (faster for specific items)
	blinkitDeliveryTime := 8 + rand.Intn(7) // 8-15 minutes

	// Generate appropriate offers based on item type
	zeptoOffer := "Free delivery"
	blinkitOffer := "15% off"

	// Specific offers based on item category
	if strings.Contains(strings.ToLower(groceryItem), "fresh") {
		zeptoOffer = "Farm fresh guarantee"
	} else if strings.Contains(strings.ToLower(groceryItem), "pack") {
		blinkitOffer = "Buy 2 Get 1 free"
	} else if strings.Contains(strings.ToLower(groceryItem), "oil") ||
		strings.Contains(strings.ToLower(groceryItem), "ghee") {
		zeptoOffer = "₹50 off on next order"
	} else if strings.Contains(strings.ToLower(groceryItem), "rice") ||
		strings.Contains(strings.ToLower(groceryItem), "flour") {
		blinkitOffer = "Free kitchen tool"
	}

	return []ServiceOffer{
		{ServiceName: "Zepto", Price: zeptoPrice, Offer: zeptoOffer, DeliveryTime: zeptoDeliveryTime},
		{ServiceName: "Blinkit", Price: blinkitPrice, Offer: blinkitOffer, DeliveryTime: blinkitDeliveryTime},
	}
}

func compareTaxi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fromCountry := strings.ToLower(r.URL.Query().Get("fromCountry"))
	fromState := strings.ToLower(r.URL.Query().Get("fromState"))
	toCountry := strings.ToLower(r.URL.Query().Get("toCountry"))
	toState := strings.ToLower(r.URL.Query().Get("toState"))

	// Create a key to look up in the taxi services map
	key := fmt.Sprintf("%s:%s:%s:%s", fromCountry, fromState, toCountry, toState)

	// Check for existing data or generate dynamic offers
	var offers []ServiceOffer
	if existingOffers, exists := taxiServices[key]; exists {
		offers = existingOffers
	} else {
		offers = generateDynamicTaxiOffers(fromState, toState)
		taxiServices[key] = offers
	}

	json.NewEncoder(w).Encode(offers)
}

// Compare restaurant delivery services
func compareRestaurant(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	country := strings.ToLower(r.URL.Query().Get("country"))
	state := strings.ToLower(r.URL.Query().Get("state"))
	city := strings.ToLower(r.URL.Query().Get("city"))
	restaurant := strings.ToLower(r.URL.Query().Get("restaurant"))

	// Create a key to look up in the restaurant services map
	key := fmt.Sprintf("%s:%s:%s:%s", country, state, city, restaurant)

	// Check for existing data or generate dynamic offers
	var offers []ServiceOffer
	if existingOffers, exists := restaurantServices[key]; exists {
		offers = existingOffers
	} else {
		offers = generateDynamicRestaurantOffers(restaurant, city)
		restaurantServices[key] = offers
	}

	json.NewEncoder(w).Encode(offers)
}

// Compare quick commerce services
func compareQuickCommerce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	country := strings.ToLower(r.URL.Query().Get("country"))
	state := strings.ToLower(r.URL.Query().Get("state"))
	city := strings.ToLower(r.URL.Query().Get("city"))
	address := strings.ToLower(r.URL.Query().Get("address"))
	groceryItem := strings.ToLower(r.URL.Query().Get("groceryItem"))

	var key string

	// Create a key to look up in the quick commerce services map
	if groceryItem != "" {
		// Include the grocery item in the key if specified
		key = fmt.Sprintf("%s:%s:%s:%s:%s", country, state, city, address, groceryItem)
	} else {
		// Otherwise use the standard key
		key = fmt.Sprintf("%s:%s:%s:%s", country, state, city, address)
	}

	// Check for existing data or generate dynamic offers
	var offers []ServiceOffer
	if existingOffers, exists := quickCommerceServices[key]; exists {
		offers = existingOffers
	} else {
		// Use appropriate generator based on whether a grocery item is specified
		if groceryItem != "" {
			offers = generateDynamicGroceryItemOffers(groceryItem, address)
		} else {
			offers = generateDynamicQuickCommerceOffers(address, city)
		}
		quickCommerceServices[key] = offers
	}

	json.NewEncoder(w).Encode(offers)
}
