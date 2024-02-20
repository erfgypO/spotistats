import axios from "axios";

const httpClient = axios.create({
  baseURL: '/backend'
});

httpClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

httpClient.interceptors.response.use((response) => {
  return response;
}, (error) => {
  if (error.response.status === 401) {
    localStorage.removeItem('token')
  }
  return Promise.reject(error);
});

export default httpClient;
