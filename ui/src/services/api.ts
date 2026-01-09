import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081/api';
const AUTH_URL = process.env.NEXT_PUBLIC_AUTH_URL || 'http://localhost:8082/api';

// API Service instance
export const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Auth Service instance
export const authClient = axios.create({
  baseURL: AUTH_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token
const addAuthToken = (config: any) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

// Response interceptor to handle errors
const handleError = (error: any) => {
  if (error.response?.status === 401) {
    // Token expired or invalid
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login';
  }
  return Promise.reject(error);
};

// Add interceptors to both clients
apiClient.interceptors.request.use(addAuthToken, (error) => Promise.reject(error));
apiClient.interceptors.response.use((response) => response, handleError);

authClient.interceptors.request.use(addAuthToken, (error) => Promise.reject(error));
authClient.interceptors.response.use((response) => response, handleError);
