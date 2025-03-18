document.addEventListener('DOMContentLoaded', function() {
    // API URL
    const API_BASE_URL = '/api';

    // WebSocket Connection
    let socket = null;
    let wsConnected = false;
    let activeSubscription = null;


    function connectWebSocket() {
        if (socket !== null) {
            return; // Already connected or connecting
        }

        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
        const wsUrl = `${protocol}//localhost:5000/ws`;

        socket = new WebSocket(wsUrl);

        socket.onopen = function() {
            console.log('WebSocket connection established');
            wsConnected = true;

            // If there's a pending subscription, send it now
            if (activeSubscription) {
                subscribeToRealTimeUpdates(activeSubscription);
            }
        };

        socket.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                console.log('Real-time update received:', data);

                // Handle real-time price updates
                if (data && data.category && data.offers) {
                    switch (data.category) {
                        case 'taxi':
                            displayTaxiResults(data.offers, data.route);
                            break;
                        case 'restaurant':
                            displayRestaurantResults(data.offers, data.location.split(', ')[0], data.location);
                            break;
                        case 'quickcommerce':
                            displayQuickCommerceResults(data.offers, data.location.split(', ')[0], data.location);
                            break;
                    }

                    // Add "Updated just now" indication
                    const timestampEl = document.createElement('div');
                    timestampEl.classList.add('update-timestamp');
                    timestampEl.textContent = 'Updated just now';

                    const existingTimestamp = document.querySelector('.update-timestamp');
                    if (existingTimestamp) {
                        existingTimestamp.replaceWith(timestampEl);
                    } else {
                        document.querySelector('.results-header').appendChild(timestampEl);
                    }
                }
            } catch (error) {
                console.error('Error processing WebSocket message:', error);
            }
        };

        socket.onclose = function() {
            console.log('WebSocket connection closed');
            wsConnected = false;
            socket = null;

            // Attempt to reconnect after a delay
            setTimeout(connectWebSocket, 3000);
        };

        socket.onerror = function(error) {
            console.error('WebSocket error:', error);
            wsConnected = false;
        };
    }

    // Subscribe to real-time updates for a specific request
    function subscribeToRealTimeUpdates(request) {
        // Save current subscription
        activeSubscription = request;

        // Send subscription request if WebSocket is connected
        if (wsConnected && socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(request));
            console.log('Subscribed to real-time updates:', request);
        } else {
            // Connection not ready, connect first
            connectWebSocket();
        }
    }

    // Connect to WebSocket on page load
    connectWebSocket();

    // DOM Elements - Common
    const loader = document.getElementById('loader');
    const results = document.getElementById('results');
    const resultsContainer = document.getElementById('results-container');
    const errorMessage = document.getElementById('error-message');
    const errorBackBtn = document.getElementById('error-back');
    const backToFormBtn = document.getElementById('back-to-form');

    // Step Navigation
    const formSteps = document.querySelectorAll('.form-step');
    const categoryCards = document.querySelectorAll('.category-card');
    const backButtons = document.querySelectorAll('.back-btn');

    // Category specific elements
    let selectedCategory = null;

    // Taxi form elements
    const taxiFromCountrySelect = document.getElementById('taxi-from-country');
    const taxiFromStateSelect = document.getElementById('taxi-from-state');
    const taxiToCountrySelect = document.getElementById('taxi-to-country');
    const taxiToStateSelect = document.getElementById('taxi-to-state');
    const compareTaxiBtn = document.getElementById('compare-taxi-btn');

    // Restaurant form elements
    const restaurantCountrySelect = document.getElementById('restaurant-country');
    const restaurantStateSelect = document.getElementById('restaurant-state');
    const restaurantCitySelect = document.getElementById('restaurant-city');
    const restaurantNameSelect = document.getElementById('restaurant-name');
    const compareRestaurantBtn = document.getElementById('compare-restaurant-btn');

    // Quick Commerce form elements
    const quickCommerceCountrySelect = document.getElementById('quickcommerce-country');
    const quickCommerceStateSelect = document.getElementById('quickcommerce-state');
    const quickCommerceCitySelect = document.getElementById('quickcommerce-city');
    const quickCommerceAddressSelect = document.getElementById('quickcommerce-address');
    const compareQuickCommerceBtn = document.getElementById('compare-quickcommerce-btn');

    // Navigation Functions

    // Show a specific form step
    function showStep(stepId) {
        // Hide all steps
        formSteps.forEach(step => {
            step.classList.remove('active');
        });

        // Show the requested step
        const targetStep = document.getElementById(stepId);
        if (targetStep) {
            targetStep.classList.add('active');
        }
    }

    // Category selection
    categoryCards.forEach(card => {
        card.addEventListener('click', () => {
            // Remove selection from all cards
            categoryCards.forEach(c => c.classList.remove('selected'));

            // Select this card
            card.classList.add('selected');

            // Get selected category
            selectedCategory = card.dataset.category;

            // Load initial options for the next step
            loadCountries(selectedCategory);

            // Show appropriate step based on category
            setTimeout(() => {
                showStep(`step-${selectedCategory}`);
            }, 200);
        });
    });

    // Back button navigation
    backButtons.forEach(btn => {
        btn.addEventListener('click', () => {
            const targetStep = btn.dataset.target;
            if (targetStep) {
                showStep(targetStep);
            }
        });
    });

    // Back to form from results
    backToFormBtn.addEventListener('click', () => {
        results.classList.add('hidden');
        showStep('step-category');
    });

    // Back from error message
    errorBackBtn.addEventListener('click', () => {
        errorMessage.classList.add('hidden');

        // Go back to the appropriate form
        if (selectedCategory) {
            showStep(`step-${selectedCategory}`);
        } else {
            showStep('step-category');
        }
    });

    // Helper Functions

    // Show loading indicator
    function showLoading() {
        document.querySelectorAll('.form-step').forEach(step => {
            step.classList.remove('active');
        });
        errorMessage.classList.add('hidden');
        results.classList.add('hidden');
        loader.classList.remove('hidden');
    }

    // Hide loading indicator
    function hideLoading() {
        loader.classList.add('hidden');
    }

    // Show error message
    function showError(message) {
        hideLoading();
        errorMessage.querySelector('p').textContent = message || 'No results found. Please try a different search.';
        errorMessage.classList.remove('hidden');
    }

    // Initialize select dropdowns with options
    function populateSelect(selectElement, options, labelProperty, valueProperty) {
        // Clear existing options except first (placeholder)
        while (selectElement.options.length > 1) {
            selectElement.remove(1);
        }

        // Add new options
        options.forEach(option => {
            const value = valueProperty ? option[valueProperty] : option;
            const label = labelProperty ? option[labelProperty] : option;

            const optionElement = document.createElement('option');
            optionElement.value = value;
            optionElement.textContent = label;
            selectElement.appendChild(optionElement);
        });

        // Enable select
        selectElement.disabled = false;
    }

    // Set default country to "India" for all category forms  
    function setDefaultCountry(selectElement) {
        const options = selectElement.options;
        for (let i = 0; i < options.length; i++) {
            if (options[i].value.toLowerCase() === "india") {
                selectElement.selectedIndex = i;
                selectElement.dispatchEvent(new Event('change')); // Trigger change event
                break;
            }
        }
    }

    // API Data Loading Functions

    // Load countries for all category forms
    async function loadCountries(category) {
        try {
            const response = await fetch(`${API_BASE_URL}/options?category=${category}`);
            const data = await response.json();

            if (data && data.countries && data.countries.length > 0) {
                const countries = data.countries;

                switch(category) {
                    case 'taxi':
                        populateSelect(taxiFromCountrySelect, countries);
                        populateSelect(taxiToCountrySelect, countries);
                        setDefaultCountry(taxiFromCountrySelect);
                        setDefaultCountry(taxiToCountrySelect);
                        break;
                    case 'restaurant':
                        populateSelect(restaurantCountrySelect, countries);
                        setDefaultCountry(restaurantCountrySelect);
                        break;
                    case 'quickcommerce':
                        populateSelect(quickCommerceCountrySelect, countries);
                        setDefaultCountry(quickCommerceCountrySelect);
                        break;
                }
            }
        } catch (error) {
            console.error('Error loading countries:', error);
        }
    }

    // Load states based on country selection
    async function loadStates(category, country, targetElement) {
        try {
            const response = await fetch(`${API_BASE_URL}/options?category=${category}&country=${country}`);
            const data = await response.json();

            if (data && data.states) {
                populateSelect(targetElement, data.states);
            }
        } catch (error) {
            console.error('Error loading states:', error);
        }
    }

    // Load cities based on country and state
    async function loadCities(category, country, state, targetElement) {
        try {
            const response = await fetch(`${API_BASE_URL}/options?category=${category}&country=${country}&state=${state}`);
            const data = await response.json();

            if (data && data.cities) {
                populateSelect(targetElement, data.cities);
            }
        } catch (error) {
            console.error('Error loading cities:', error);
        }
    }

    // Load restaurants or addresses based on location
    async function loadFinalOptions(category, country, state, city, targetElement, address) {
        try {
            let url = `${API_BASE_URL}/options?category=${category}&country=${country}&state=${state}&city=${city}`;

            // Add address parameter if provided (for grocery items)
            if (address) {
                url += `&address=${address}`;
            }

            const response = await fetch(url);
            const data = await response.json();

            if (data) {
                if (category === 'restaurant' && data.restaurants) {
                    populateSelect(targetElement, data.restaurants);
                } else if (category === 'quickcommerce' && data.addresses) {
                    populateSelect(targetElement, data.addresses);

                    // Reset grocery item select when address changes
                    const groceryItemSelect = document.getElementById('quickcommerce-grocery-item');
                    if (groceryItemSelect) {
                        groceryItemSelect.innerHTML = '<option value="" disabled selected>Select Grocery Item</option>';
                        groceryItemSelect.disabled = true;
                        document.querySelector('.grocery-item-wrapper').classList.add('hidden');
                        document.querySelector('.grocery-item-wrapper').classList.remove('active');
                    }
                } else if (category === 'quickcommerce' && data.groceryItems) {
                    // Populate grocery items
                    populateSelect(targetElement, data.groceryItems);
                    targetElement.disabled = false;
                    document.querySelector('.grocery-item-wrapper').classList.remove('hidden');
                    document.querySelector('.grocery-item-wrapper').classList.add('active');
                }
            }
        } catch (error) {
            console.error('Error loading options:', error);
        }
    }

    // Event listeners for select changes

    // Taxi events
    taxiFromCountrySelect.addEventListener('change', () => {
        loadStates('taxi', taxiFromCountrySelect.value, taxiFromStateSelect);
    });

    taxiToCountrySelect.addEventListener('change', () => {
        loadStates('taxi', taxiToCountrySelect.value, taxiToStateSelect);
    });

    // Enable the compare button when all selects are filled
    [taxiFromStateSelect, taxiToStateSelect].forEach(select => {
        select.addEventListener('change', () => {
            if (taxiFromCountrySelect.value && taxiFromStateSelect.value && 
                taxiToCountrySelect.value && taxiToStateSelect.value) {
                compareTaxiBtn.disabled = false;
            }
        });
    });

    // Restaurant events
    restaurantCountrySelect.addEventListener('change', () => {
        loadStates('restaurant', restaurantCountrySelect.value, restaurantStateSelect);
    });

    restaurantStateSelect.addEventListener('change', () => {
        const country = restaurantCountrySelect.value;
        const state = restaurantStateSelect.value;
        if (country && state) {
            loadCities('restaurant', country, state, restaurantCitySelect);
        }
    });

    restaurantCitySelect.addEventListener('change', () => {
        const country = restaurantCountrySelect.value;
        const state = restaurantStateSelect.value;
        const city = restaurantCitySelect.value;

        if (country && state && city) {
            loadFinalOptions('restaurant', country, state, city, restaurantNameSelect);
        }
    });

    restaurantNameSelect.addEventListener('change', () => {
        if (restaurantNameSelect.value) {
            compareRestaurantBtn.disabled = false;
        }
    });

    // Quick Commerce events
    quickCommerceCountrySelect.addEventListener('change', () => {
        loadStates('quickcommerce', quickCommerceCountrySelect.value, quickCommerceStateSelect);
    });

    quickCommerceStateSelect.addEventListener('change', () => {
        const country = quickCommerceCountrySelect.value;
        const state = quickCommerceStateSelect.value;
        if (country && state) {
            loadCities('quickcommerce', country, state, quickCommerceCitySelect);
        }
    });

    quickCommerceCitySelect.addEventListener('change', () => {
        const country = quickCommerceCountrySelect.value;
        const state = quickCommerceStateSelect.value;
        const city = quickCommerceCitySelect.value;

        if (country && state && city) {
            loadFinalOptions('quickcommerce', country, state, city, quickCommerceAddressSelect);
        }
    });

    quickCommerceAddressSelect.addEventListener('change', () => {
        const country = quickCommerceCountrySelect.value;
        const state = quickCommerceStateSelect.value;
        const city = quickCommerceCitySelect.value;
        const address = quickCommerceAddressSelect.value;

        if (address) {
            compareQuickCommerceBtn.disabled = false;

            // Get grocery item select
            const groceryItemSelect = document.getElementById('quickcommerce-grocery-item');

            // Load grocery items for this address
            loadFinalOptions('quickcommerce', country, state, city, groceryItemSelect, address);
        }
    });

    // Grocery item selection
    const groceryItemSelect = document.getElementById('quickcommerce-grocery-item');
    groceryItemSelect.addEventListener('change', () => {
        // Enable compare button (should already be enabled from address selection)
        if (groceryItemSelect.value) {
            compareQuickCommerceBtn.disabled = false;
        }
    });

    // Comparison button event listeners

    // Taxi comparison
    compareTaxiBtn.addEventListener('click', async () => {
        showLoading();

        const fromCountry = taxiFromCountrySelect.value;
        const fromState = taxiFromStateSelect.value;
        const toCountry = taxiToCountrySelect.value;
        const toState = taxiToStateSelect.value;

        try {
            // Initial fetch using regular API
            const queryParams = new URLSearchParams({
                fromCountry: fromCountry,
                fromState: fromState,
                toCountry: toCountry,
                toState: toState
            });

            const response = await fetch(`${API_BASE_URL}/compare/taxi?${queryParams}`);
            const data = await response.json();

            if (response.status >= 400 || data.error) {
                showError(data.error || 'No taxi services found for this route.');
                return;
            }

            if (data && data.length > 0) {
                displayTaxiResults(data, `${fromState} to ${toState}`);

                // Subscribe to real-time updates
                subscribeToRealTimeUpdates({
                    category: 'taxi',
                    fromCountry: fromCountry,
                    fromState: fromState,
                    toCountry: toCountry,
                    toState: toState
                });
            } else {
                showError('No taxi services available for this route.');
            }
        } catch (error) {
            console.error('Error comparing taxi services:', error);
            showError('Error connecting to the server. Please try again.');
        } finally {
            hideLoading();
        }
    });

    // Restaurant comparison
    compareRestaurantBtn.addEventListener('click', async () => {
        showLoading();

        const country = restaurantCountrySelect.value;
        const state = restaurantStateSelect.value;
        const city = restaurantCitySelect.value;
        const restaurant = restaurantNameSelect.value;

        try {
            const queryParams = new URLSearchParams({
                country: country,
                state: state,
                city: city,
                restaurant: restaurant
            });

            const response = await fetch(`${API_BASE_URL}/compare/restaurant?${queryParams}`);
            const data = await response.json();

            if (response.status >= 400 || data.error) {
                showError(data.error || 'No delivery services found for this restaurant.');
                return;
            }

            if (data && data.length > 0) {
                displayRestaurantResults(data, restaurant, `${city}, ${state}`);

                // Subscribe to real-time updates
                subscribeToRealTimeUpdates({
                    category: 'restaurant',
                    country: country,
                    state: state,
                    city: city,
                    restaurant: restaurant
                });
            } else {
                showError('No delivery services available for this restaurant and location.');
            }
        } catch (error) {
            console.error('Error comparing restaurant services:', error);
            showError('Error connecting to the server. Please try again.');
        } finally {
            hideLoading();
        }
    });

    // Quick Commerce comparison
    compareQuickCommerceBtn.addEventListener('click', async () => {
        showLoading();

        const country = quickCommerceCountrySelect.value;
        const state = quickCommerceStateSelect.value;
        const city = quickCommerceCitySelect.value;
        const address = quickCommerceAddressSelect.value;
        const groceryItem = document.getElementById('quickcommerce-grocery-item').value;

        try {
            const queryParams = new URLSearchParams({
                country: country,
                state: state,
                city: city,
                address: address
            });

            // Add grocery item if selected
            if (groceryItem) {
                queryParams.append('groceryItem', groceryItem);
            }

            const response = await fetch(`${API_BASE_URL}/compare/quickcommerce?${queryParams}`);
            const data = await response.json();

            if (response.status >= 400 || data.error) {
                showError(data.error || 'No quick commerce services found for this location.');
                return;
            }

            if (data && data.length > 0) {
                // Display with grocery item if selected
                let displayAddress = address;
                if (groceryItem) {
                    displayAddress = `${groceryItem} to ${address}`;
                }

                displayQuickCommerceResults(data, displayAddress, `${city}, ${state}`);

                // Subscribe to real-time updates
                const subscriptionData = {
                    category: 'quickcommerce',
                    country: country,
                    state: state,
                    city: city,
                    address: address
                };

                // Add grocery item to subscription if selected
                if (groceryItem) {
                    subscriptionData.groceryItem = groceryItem;
                }

                subscribeToRealTimeUpdates(subscriptionData);
            } else {
                showError('No quick commerce services available for this location.');
            }
        } catch (error) {
            console.error('Error comparing quick commerce services:', error);
            showError('Error connecting to the server. Please try again.');
        } finally {
            hideLoading();
        }
    });

    // Results display functions

    // Taxi results
    function displayTaxiResults(offers, route) {
        results.classList.remove('hidden');

        // Check if this is an update to existing results or a new search
        const isUpdate = resultsContainer.querySelector('.route-info') !== null;

        if (!isUpdate) {
            resultsContainer.innerHTML = '';

            // Add route information
            const routeInfo = document.createElement('div');
            routeInfo.className = 'route-info';
            routeInfo.innerHTML = `<h4>Taxi from ${route}</h4>`;
            resultsContainer.appendChild(routeInfo);
        }

        // Sort offers by price
        offers.sort((a, b) => a.Price - b.Price);

        // Find the best deal
        const bestDeal = offers[0];

        // Create or update result cards
        offers.forEach((offer, index) => {
            // Check if card already exists for this service
            let card = resultsContainer.querySelector(`.result-card[data-service="${offer.ServiceName}"]`);
            const isNewCard = !card;

            // Create new card if needed
            if (isNewCard) {
                card = document.createElement('div');
                card.className = 'result-card';
                card.setAttribute('data-service', offer.ServiceName);
            }

            // Apply best deal class
            if (offer.Price === bestDeal.Price) {
                card.classList.add('best-deal');
            } else {
                card.classList.remove('best-deal');
            }

            // Update price with animation if it changed
            const oldPriceElement = card.querySelector('.price');
            const oldPrice = oldPriceElement ? parseFloat(oldPriceElement.textContent.replace('₹', '')) : null;

            // Create content
            card.innerHTML = `
                <h4>${offer.ServiceName}</h4>
                <div class="price${oldPrice !== null && oldPrice !== offer.Price ? ' price-changed' : ''}">₹${offer.Price.toFixed(2)}</div>
                <div class="offer"><i class="fas fa-tag"></i> ${offer.Offer}</div>
                <div class="duration"><i class="fas fa-clock"></i> ${Math.floor(offer.Duration / 60)}h ${offer.Duration % 60}m</div>
            `;

            // Add to container if new
            if (isNewCard) {
                resultsContainer.appendChild(card);
            }
        });

        // Update the results header with real-time indication if this is an update
        if (isUpdate) {
            const timestamp = document.createElement('div');
            timestamp.className = 'update-timestamp';
            timestamp.textContent = 'Updated just now';

            const existingTimestamp = document.querySelector('.update-timestamp');
            if (existingTimestamp) {
                existingTimestamp.replaceWith(timestamp);
            } else {
                document.querySelector('.results-header').appendChild(timestamp);
            }
        }
    }

    // Restaurant results
    function displayRestaurantResults(offers, restaurant, location) {
        results.classList.remove('hidden');

        // Check if this is an update to existing results or a new search
        const isUpdate = resultsContainer.querySelector('.restaurant-info') !== null;

        if (!isUpdate) {
            resultsContainer.innerHTML = '';

            // Add restaurant information
            const restaurantInfo = document.createElement('div');
            restaurantInfo.className = 'restaurant-info';
            restaurantInfo.innerHTML = `<h4>${restaurant} in ${location}</h4>`;
            resultsContainer.appendChild(restaurantInfo);
        }

        // Sort offers by price
        offers.sort((a, b) => a.Price - b.Price);

        // Find the best deal
        const bestDeal = offers[0];

        // Create or update result cards
        offers.forEach((offer, index) => {
            // Check if card already exists for this service
            let card = resultsContainer.querySelector(`.result-card[data-service="${offer.ServiceName}"]`);
            const isNewCard = !card;

            // Create new card if needed
            if (isNewCard) {
                card = document.createElement('div');
                card.className = 'result-card';
                card.setAttribute('data-service', offer.ServiceName);
            }

            // Apply best deal class
            if (offer.Price === bestDeal.Price) {
                card.classList.add('best-deal');
            } else {
                card.classList.remove('best-deal');
            }

            // Update price with animation if it changed
            const oldPriceElement = card.querySelector('.price');
            const oldPrice = oldPriceElement ? parseFloat(oldPriceElement.textContent.replace('₹', '')) : null;

            // Create content
            card.innerHTML = `
                <h4>${offer.ServiceName}</h4>
                <div class="price${oldPrice !== null && oldPrice !== offer.Price ? ' price-changed' : ''}">₹${offer.Price.toFixed(2)}</div>
                <div class="offer"><i class="fas fa-tag"></i> ${offer.Offer}</div>
                <div class="delivery-time"><i class="fas fa-clock"></i> ${offer.DeliveryTime} minutes</div>
            `;

            // Add to container if new
            if (isNewCard) {
                resultsContainer.appendChild(card);
            }
        });

        // Update the results header with real-time indication if this is an update
        if (isUpdate) {
            const timestamp = document.createElement('div');
            timestamp.className = 'update-timestamp';
            timestamp.textContent = 'Updated just now';

            const existingTimestamp = document.querySelector('.update-timestamp');
            if (existingTimestamp) {
                existingTimestamp.replaceWith(timestamp);
            } else {
                document.querySelector('.results-header').appendChild(timestamp);
            }
        }
    }

    // Quick Commerce results
    function displayQuickCommerceResults(offers, address, location) {
        results.classList.remove('hidden');

        // Check if this is an update to existing results or a new search
        const isUpdate = resultsContainer.querySelector('.address-info') !== null;

        if (!isUpdate) {
            resultsContainer.innerHTML = '';

            // Add address information
            const addressInfo = document.createElement('div');
            addressInfo.className = 'address-info';

            // Check if this is a grocery item delivery (format will be "item to address")
            if (address.includes(' to ')) {
                // This is a grocery item delivery
                addressInfo.innerHTML = `<h4>${address}, ${location}</h4>`;
            } else {
                // This is a regular quick commerce delivery
                addressInfo.innerHTML = `<h4>Delivery to ${address}, ${location}</h4>`;
            }

            resultsContainer.appendChild(addressInfo);
        }

        // Sort offers by price
        offers.sort((a, b) => a.Price - b.Price);

        // Find the best deal
        const bestDeal = offers[0];

        // Create or update result cards
        offers.forEach((offer, index) => {
            // Check if card already exists for this service
            let card = resultsContainer.querySelector(`.result-card[data-service="${offer.ServiceName}"]`);
            const isNewCard = !card;

            // Create new card if needed
            if (isNewCard) {
                card = document.createElement('div');
                card.className = 'result-card';
                card.setAttribute('data-service', offer.ServiceName);
            }

            // Apply best deal class
            if (offer.Price === bestDeal.Price) {
                card.classList.add('best-deal');
            } else {
                card.classList.remove('best-deal');
            }

            // Update price with animation if it changed
            const oldPriceElement = card.querySelector('.price');
            const oldPrice = oldPriceElement ? parseFloat(oldPriceElement.textContent.replace('₹', '')) : null;

            // Create content
            card.innerHTML = `
                <h4>${offer.ServiceName}</h4>
                <div class="price${oldPrice !== null && oldPrice !== offer.Price ? ' price-changed' : ''}">₹${offer.Price.toFixed(2)}</div>
                <div class="offer"><i class="fas fa-tag"></i> ${offer.Offer}</div>
                <div class="delivery-time"><i class="fas fa-clock"></i> ${offer.DeliveryTime} minutes</div>
            `;

            // Add to container if new
            if (isNewCard) {
                resultsContainer.appendChild(card);
            }
        });

        // Update the results header with real-time indication if this is an update
        if (isUpdate) {
            const timestamp = document.createElement('div');
            timestamp.className = 'update-timestamp';
            timestamp.textContent = 'Updated just now';

            const existingTimestamp = document.querySelector('.update-timestamp');
            if (existingTimestamp) {
                existingTimestamp.replaceWith(timestamp);
            } else {
                document.querySelector('.results-header').appendChild(timestamp);
            }
        }
    }

    // Initial animation for the form
    document.getElementById('form').style.opacity = '0';
    setTimeout(() => {
        document.getElementById('form').style.opacity = '1';
    }, 300);
});
