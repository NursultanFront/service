// Utilities
import axios from "axios";

// One shared client for the sales API. withCredentials is the axios
// equivalent of fetch's `credentials: "include"` - it makes the browser
// attach the httpOnly auth cookie to every request automatically.
const api = axios.create({
  baseURL: import.meta.env.VITE_SERVICE_API,
  withCredentials: true,
});

export default api;
