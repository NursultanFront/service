<template>
  <v-navigation-drawer expand-on-hover rail>
    <v-list>
      <v-list-item
        prepend-avatar="https://www.ardanlabs.com/images/ardanlabs-logo.svg"
      ></v-list-item>
    </v-list>

    <v-divider></v-divider>

    <v-list density="compact" nav>
      <v-list-item
        v-for="menu in availableMenus"
        :key="menu.name"
        class="bar__items"
        :prepend-icon="menu.icon"
        :title="menu.title"
        :to="{ name: menu.name }"
        variant="text"
      ></v-list-item>
    </v-list>

    <template #append>
      <v-list density="compact" nav>
        <v-list-item
          class="bar__items"
          prepend-icon="fas fa-sign-out-alt"
          title="Logout"
          variant="text"
          @click="logout"
        ></v-list-item>
      </v-list>
    </template>
  </v-navigation-drawer>
</template>

<script setup>
import { useRouter } from "vue-router";
import { useAuthStore } from "@/store/auth";

const availableMenus = [
  { title: "Users", name: "Users", icon: "fas fa-users" },
];

const authStore = useAuthStore();
const router = useRouter();

async function logout() {
  await authStore.logout();
  router.push({ name: "Login" });
}
</script>

<style>
.bar__items {
  .v-list-item-title {
    color: #272727;
  }
  .v-icon {
    color: #89a1b0;
  }
}
</style>
