// Utilities
import { defineStore } from "pinia";
import api from "@/plugins/axios";

export const useUsersStore = defineStore("users", () => {
  async function fetchUsers({ page, rows, sort = "" }) {
    const { data } = await api.get(`/users?page=${page}&rows=${rows}${sort}`);
    return data;
  }

  async function fetchUser(userId) {
    const { data } = await api.get(`/users/${userId}`);
    return data;
  }

  async function createUser(form) {
    const { data } = await api.post("/users", form);
    return data;
  }

  async function updateUser(userId, form) {
    const { data } = await api.put(`/users/${userId}`, form);
    return data;
  }

  async function deleteUser(userId) {
    await api.delete(`/users/${userId}`);
  }

  return { fetchUsers, fetchUser, createUser, updateUser, deleteUser };
});
