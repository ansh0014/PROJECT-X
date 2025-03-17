let socket = null;
let isConnected = false;
let currentSubscription = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;

// Connect to WebSocket server
function connectWebSocket() {

    if (socket) {
        socket.close();
    }

    // Create WebSocket URL using the current location
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = `${protocol}//localhost:5000/ws`;
    
    // Create new WebSocket connection
    socket = new WebSocket(wsUrl);
    
    // Connection opened
    socket.addEventListener("open", (event) => {
        console.log("WebSocket connection established");
        isConnected = true;
        reconnectAttempts = 0;
        
        // If there's a pending subscription, send it
        if (currentSubscription) {
            subscribeToRealTimeUpdates(currentSubscription);
        }
    });
    
    // Connection closed
    socket.addEventListener("close", (event) => {
        console.log("WebSocket connection closed");
        isConnected = false;
        
        // Try to reconnect if not max attempts
        if (reconnectAttempts < maxReconnectAttempts) {
            reconnectAttempts++;
            console.log(`Attempting to reconnect (${reconnectAttempts}/${maxReconnectAttempts})...`);
            setTimeout(connectWebSocket, 2000); // Try to reconnect after 2 seconds
        } else {
            console.log("Failed to reconnect after multiple attempts");
        }
    });
    
    // Connection error
    socket.addEventListener("error", (event) => {
        console.error("WebSocket error:", event);
    });
    
    // Listen for messages
    socket.addEventListener("message", (event) => {
        try {
            const data = JSON.parse(event.data);
            handleRealTimeUpdate(data);
        } catch (error) {
            console.error("Error parsing WebSocket message:", error);
        }
    });
}

// Subscribe to real-time updates for a specific search
function subscribeToRealTimeUpdates(request) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        // Save the subscription for when connection is established
        currentSubscription = request;
        
        // Try to connect if not already connecting
        if (!isConnected && reconnectAttempts < maxReconnectAttempts) {
            connectWebSocket();
        }
        return;
    }
    
    // Send subscription request
    const message = JSON.stringify(request);
    socket.send(message);
    currentSubscription = request;
    
    console.log("Subscribed to real-time updates:", request);
}

// Handle incoming real-time updates
function handleRealTimeUpdate(data) {
    console.log("Received real-time update:", data);
    
    // Display updated results based on category
    switch (data.category) {
        case "taxi":
            displayTaxiResults(data.offers, data.route, true);
            break;
        case "restaurant":
            displayRestaurantResults(data.offers, data.route, data.location, true);
            break;
        case "quickcommerce":
            displayQuickCommerceResults(data.offers, data.route, data.location, true);
            break;
    }
    
    // Add timestamp to show when prices were updated
    updateTimestamp(data.timestamp);
}

// Update the timestamp display
function updateTimestamp(timestamp) {
    const timestampElement = document.getElementById('price-timestamp');
    if (timestampElement) {
        const date = new Date(timestamp * 1000);
        const timeString = date.toLocaleTimeString();
        timestampElement.textContent = `Prices updated at ${timeString}`;
        
        // Show the timestamp with animation
        timestampElement.classList.add('pulse');
        setTimeout(() => {
            timestampElement.classList.remove('pulse');
        }, 1000);
    }
}

// Connect when the page loads
document.addEventListener('DOMContentLoaded', () => {
    connectWebSocket();
});

// Export functions for use in script.js
window.realtime = {
    subscribeToRealTimeUpdates
};