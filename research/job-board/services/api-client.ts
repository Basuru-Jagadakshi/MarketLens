import axios from "axios";
import { ENV } from "@/config/env";

export const apiClient = axios.create({
  baseURL: ENV.NEXT_PUBLIC_API_BASE_URL,
  timeout: ENV.TIMEOUT_MS,
  headers: {
    "Content-Type": "application/json",
  },
});

// Observability & Security Interceptor
apiClient.interceptors.request.use(
  (config) => {
    // Standard injection point for backend Auth tokens (e.g., Bearer tokens)
    return config;
  },
  (error) => {
    console.error(`[API Request Error]:`, error);
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Graceful error abstraction framework
    const errorMessage = error.response?.data?.message || "An unexpected system anomaly occurred";
    console.error(`[API Response Exception] [${error.response?.status}]:`, errorMessage);
    return Promise.reject(new Error(errorMessage));
  }
);