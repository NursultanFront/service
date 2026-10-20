// Utilities
import { defineStore } from "pinia";
import api from "@/plugins/axios";

export const useHomesStore = defineStore("homes", () => {
  async function fetchHomes({ userId, page, rows, sort = "" }) {
    const { data } = await api.get(
      `/homes?user_id=${userId}&page=${page}&rows=${rows}${sort}`
    );
    return data;
  }

  async function createHome(form) {
    const { data } = await api.post("/homes", form);
    return data;
  }

  async function updateHome(homeId, form) {
    const { data } = await api.put(`/homes/${homeId}`, form);
    return data;
  }

  async function deleteHome(homeId) {
    await api.delete(`/homes/${homeId}`);
  }

  return { fetchHomes, createHome, updateHome, deleteHome };
});
