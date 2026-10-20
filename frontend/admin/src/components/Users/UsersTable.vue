<template>
  <data-table-server
    v-model:items-per-page="tableOptions.itemsPerPage"
    v-model:sort-by="tableOptions.sortBy"
    v-model:page="tableOptions.page"
    :headers="headers"
    :items-length="serverItemsLength"
    :items="users"
    :loading="loading"
    :items-per-page-options="usersItemsPerPageOptions"
    :show-select="false"
    is-first-column-fixed
    hide-title
    has-actions
    class="elevation-1"
    @click:row="goToClientsProfile"
    @update:options="loadItems"
  >
    <template v-slot:[`item.user_id`]="{ item }">
      <div>{{ item.id }}</div>
    </template>
    <template v-slot:[`item.roles`]="{ item }">
      <div>{{ item.roles.join(", ") }}</div>
    </template>
    <template v-slot:[`item.dateCreated`]="{ item }">
      <div>{{ item.dateCreated.substring(0, 10) }}</div>
    </template>
    <template v-slot:[`item.dateUpdated`]="{ item }">
      <div>{{ item.dateUpdated.substring(0, 10) }}</div>
    </template>
    <template #[`item.actions`]="{ item }">
      <users-table-actions
        @delete="$emit('delete', item)"
        @edit="$emit('edit', item)"
        @profile="goToClientsProfile(item.id)"
        :item="item"
      />
    </template>
  </data-table-server>
</template>
<script setup>
import { reactive, ref, computed } from "vue";
import { useRouter } from "vue-router";
import DataTableServer from "../DataTable/DataTableServer.vue";
import { UsersTableHeaders } from "../Users/Users.js";
import UsersTableActions from "../Users/UsersTableActions.vue";
import SortQuery from "../DataTable/SortQuery";
import { useUsersStore } from "@/store/users";

defineOptions({ name: "UsersTable" });

defineEmits(["delete", "edit"]);

const router = useRouter();
const usersStore = useUsersStore();

const tableOptions = reactive({
  page: 1,
  itemsPerPage: 5,
  sortBy: [],
});
const error = ref({});
const users = ref([]);
const loading = ref(false);
const serverItemsLength = ref(0);
const usersItemsPerPageOptions = [
  { title: "5", value: 5 },
  { title: "10", value: "10" },
  { title: "20", value: "20" },
];

const headers = computed(() => UsersTableHeaders);

function sortQuery(s) {
  return SortQuery(s);
}

function goToClientsProfile(id, e) {
  let userId = id;
  if (typeof userId !== "string") {
    userId = e.item.id;
  }
  router.push({
    name: "UserProfile",
    params: { id: userId },
  });
}

async function loadItems() {
  const { page, itemsPerPage, sortBy } = tableOptions;

  const sort = sortQuery(sortBy);

  try {
    const data = await usersStore.fetchUsers({
      page,
      rows: itemsPerPage,
      sort,
    });

    serverItemsLength.value = data.total;
    users.value = data.items;
  } catch (err) {
    console.log("fetchedCall failed:", err);
    error.value = err;
  }
}

defineExpose({ loadItems });
</script>
