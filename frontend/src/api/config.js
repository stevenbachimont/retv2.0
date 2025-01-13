export const API_BASE_URL = 'http://localhost:8080'; // Ajustez selon votre configuration backend

export const apiClient = {
    async get(endpoint) {
        const response = await fetch(`${API_BASE_URL}${endpoint}`);
        return response.json();
    },

    async post(endpoint, data) {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        return response.json();
    },

    // Ajoutez d'autres méthodes selon vos besoins (PUT, DELETE, etc.)
}; 