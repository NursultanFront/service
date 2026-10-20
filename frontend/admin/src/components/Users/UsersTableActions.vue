<template>
  <v-btn-toggle
    class="rounded-lg users-actions elevation-2"
    borderless
    background-color="rgba(255, 255, 255)"
  >
    <v-tooltip bottom>
      <template #activator="{ props }">
        <v-btn
          class="text--dark ma-1"
          size="x-small"
          v-bind="props"
          @click.stop.prevent="action('profile')"
        >
          <v-icon icon="fas fa-user"> </v-icon>
        </v-btn>
      </template>
      <span>User Profile</span>
    </v-tooltip>
    <v-tooltip bottom>
      <template #activator="{ props }">
        <v-btn
          class="text--dark ma-1"
          size="x-small"
          v-bind="props"
          @click.stop.prevent="action('edit')"
        >
          <v-icon icon="fas fa-pencil"> </v-icon>
        </v-btn>
      </template>
      <span>Edit User</span>
    </v-tooltip>
    <v-tooltip bottom>
      <template #activator="{ props }">
        <v-btn
          class="ma-1"
          size="x-small"
          v-bind="props"
          @click.stop.prevent="action('delete')"
        >
          <v-icon icon="fas fa-trash"></v-icon>
        </v-btn>
      </template>
      <span> Delete User </span>
    </v-tooltip>
  </v-btn-toggle>
</template>

<script setup>
defineOptions({ name: "UsersDataTableActions" });

defineProps({
  item: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(["profile", "edit", "delete"]);

function action(type) {
  emit(type);
}
</script>
