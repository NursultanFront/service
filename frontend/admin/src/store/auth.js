// Utilities
import { ref } from "vue";
import { defineStore } from "pinia";
import axios from "axios";

const AUTH_API = import.meta.env.VITE_AUTH_API;
const AUTH_KID = import.meta.env.VITE_AUTH_KID;

// The real JWT now lives only in an httpOnly cookie set by the auth service -
// JS can never read it (that's the point, it closes off token theft via
// XSS). All this store keeps is a non-sensitive "am I logged in" flag so the
// router guard can decide whether to show the login page.
export const useAuthStore = defineStore("auth", () => {
  const isAuthenticated = ref(localStorage.getItem("is_authenticated") === "true");

  async function login(email, password) {
    try {
      await axios.get(`${AUTH_API}/auth/token/${AUTH_KID}`, {
        withCredentials: true,
        auth: { username: email, password },
      });
    } catch (error) {
      const message =
        error.response?.data?.message || "Invalid email or password";
      throw new Error(message, { cause: error });
    }

    isAuthenticated.value = true;
    localStorage.setItem("is_authenticated", "true");
  }

  // Verifies against the server whether the httpOnly cookie is still valid -
  // the localStorage flag alone is just a leftover memory of a past login,
  // not proof the token still works (it could have expired or been cleared).
  async function checkAuth() {
    try {
      await axios.get(`${AUTH_API}/auth/authenticate`, {
        withCredentials: true,
      });
      isAuthenticated.value = true;
      localStorage.setItem("is_authenticated", "true");
    } catch {
      isAuthenticated.value = false;
      localStorage.removeItem("is_authenticated");
    }
  }

  async function logout() {
    isAuthenticated.value = false;
    localStorage.removeItem("is_authenticated");

    try {
      await axios.post(`${AUTH_API}/auth/logout`, {}, { withCredentials: true });
    } catch {
      // best-effort: the cookie will also just expire on its own
    }
  }

  return { isAuthenticated, login, logout, checkAuth };
});
