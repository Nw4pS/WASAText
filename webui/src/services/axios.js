import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
	headers: {
    'Content-Type': 'application/json',
  }
});

// prima di inviare qualsiasi richiesta http preleva il token salvato al login nel localStorage e lo inserisce come Bearer
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default instance;
