/**
 * API Client & Auth Handler
 * Handles communication with the Go Backend
 */
const API = {
    // Base URL is injected via global variable or calculated
    baseUrl: window.APP_URL || window.location.origin,

    // Helper to make authenticated requests
    async fetch(url, options = {}) {
        // Ensure URL is complete
        const fullUrl = url.startsWith('http') ? url : `${this.baseUrl}${url}`;
        
        const token = localStorage.getItem('accessToken');
        
        // Get CSRF Token from Meta Tag (Injected by Go Template)
        const csrfTokenMeta = document.querySelector('meta[name="csrf-token"]');
        const csrfToken = csrfTokenMeta ? csrfTokenMeta.getAttribute('content') : '';
        
        const headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'X-CSRF-TOKEN': csrfToken, // Header for API middleware
            ...options.headers
        };

        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = {
            ...options,
            headers
        };

        // Handle CSRF for FormData (if body is not string)
        if (options.body && typeof options.body === 'string') {
            // If it's a string (JSON), we already sent header.
            // If strict CSRF middleware checks body key, append it (usually header is enough for AJAX)
        }

        const response = await fetch(fullUrl, config);
        
        // Handle Token Expiry (401 Unauthorized)
        if (response.status === 401) {
            // Try to refresh token if we have a refresh token
            // For now, simple logout logic
            if (!window.location.pathname.includes('/login') && 
                !window.location.pathname.includes('/register') &&
                !window.location.pathname.includes('/forgot-password')) {
                
                // alert('Session expired. Please login again.'); // Optional UI feedback
                API.logout();
            }
        }

        return response;
    },

    saveTokens(tokens) {
        if (tokens.access && tokens.access.token) {
            localStorage.setItem('accessToken', tokens.access.token);
        }
        if (tokens.refresh && tokens.refresh.token) {
            localStorage.setItem('refreshToken', tokens.refresh.token);
        }
    },

    logout() {
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        window.location.href = `${this.baseUrl}/login`;
    },

    checkAuth() {
        const token = localStorage.getItem('accessToken');
        const path = window.location.pathname;
        
        // Public paths that don't require auth
        const publicPaths = ['/login', '/register', '/forgot-password'];
        const isPublic = publicPaths.some(p => path.includes(p));

        if (!token && !isPublic) {
            window.location.href = `${this.baseUrl}/login`;
        }
        
        // If logged in and trying to access login page, redirect to dashboard
        if (token && isPublic) {
            window.location.href = `${this.baseUrl}/`;
        }
    },
    
    // Helper to parse JWT payload
    getUser() {
        const token = localStorage.getItem('accessToken');
        if(!token) return null;
        try {
            return JSON.parse(atob(token.split('.')[1]));
        } catch (e) {
            return null;
        }
    }
};

// Run Auth Check on Load
document.addEventListener('DOMContentLoaded', () => {
    // Only run auth check if we are not on the Swagger UI
    if (!window.location.pathname.includes('/swagger')) {
        API.checkAuth();
    }
});