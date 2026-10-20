<template>
  <v-app>
    <v-main class="d-flex align-center justify-center login-page">
      <v-card class="pa-6" width="400" elevation="4">
        <v-card-subtitle class="text-center pb-4">Admin Sign In</v-card-subtitle>
        <v-card-text>
          <v-form ref="form" v-model="valid" @submit.prevent="submit">
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              variant="outlined"
              :rules="[requiredRule]"
              autofocus
            />
            <v-text-field
              v-model="password"
              label="Password"
              :type="visible ? 'text' : 'password'"
              :append-inner-icon="visible ? 'fas fa-eye' : 'fas fa-eye-slash'"
              variant="outlined"
              :rules="[requiredRule]"
              @click:append-inner="visible = !visible"
            />
            <v-alert
              v-if="error"
              type="error"
              density="compact"
              class="mb-4"
            >
              {{ error }}
            </v-alert>
            <v-btn block color="primary" type="submit" :loading="loading">
              Sign In
            </v-btn>
          </v-form>
        </v-card-text>
      </v-card>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAuthStore } from "@/store/auth";

defineOptions({ name: "Login" });

const email = ref("admin@example.com");
const password = ref("gophers");
const visible = ref(false);
const valid = ref(false);
const loading = ref(false);
const error = ref("");
const form = ref(null);

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

function requiredRule(v) {
  return !!v || "This field is required";
}

async function submit() {
  await form.value.validate();
  if (!valid.value) {
    return;
  }

  error.value = "";
  loading.value = true;

  try {
    await authStore.login(email.value, password.value);

    const redirect = route.query.redirect || { name: "Users" };
    router.push(redirect);
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background-color: #f5f5f5;
}
</style>
