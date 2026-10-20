<template>
  <div class="ma-6 fill-height justify-start">
    <v-responsive class="align-start text-center fill-height">
      <v-card>
        <v-card-title class="d-flex align-center justify-space-between">
          <div>
            <v-btn
              icon="fas fa-arrow-left"
              flat
              @click="$router.push({ name: 'Users' })"
            />
            User Profile
          </div>
        </v-card-title>
        <v-card-text class="d-flex pa-0">
          <v-col cols="4">
            <user-details :user="user" />
          </v-col>
          <v-col cols="8">
            <user-homes :user-id="user.id" />
          </v-col>
        </v-card-text>
      </v-card>
    </v-responsive>
  </div>
</template>
<script setup>
import { ref, onMounted, watch } from "vue";
import UserDetails from "../Users/UserDetails.vue";
import UserHomes from "../Users/UserHomes.vue";
import { useUsersStore } from "@/store/users";

const props = defineProps({
  userId: {
    type: String,
    default: "",
  },
});

const usersStore = useUsersStore();

const error = ref(null);
const user = ref({});

async function fetchUser(userId) {
  try {
    user.value = await usersStore.fetchUser(userId);
  } catch (err) {
    console.log("User fetch failed:", err);
    error.value = err;
  }
}

onMounted(() => {
  fetchUser(props.userId);
});

watch(
  () => props.userId,
  (newVal) => fetchUser(newVal)
);
</script>
