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
    @update:options="loadItems"
  >
    <template v-slot:[`item.address.country`]="{ item }">
      <div>{{ getCountry(item.address.country) }}</div>
    </template>
    <template #[`item.actions`]="{ item }">
      <users-home-table-actions
        @delete="$emit('delete', item)"
        @edit="$emit('edit', item)"
        :item="item"
      />
    </template>
  </data-table-server>
</template>
<script setup>
import { reactive, ref, computed, watch } from "vue";
import DataTableServer from "../DataTable/DataTableServer.vue";
import { UserHomesTableHeaders } from "../Users/Users.js";
import UsersHomeTableActions from "../Users/UsersHomeTableActions.vue";
import SortQuery from "../DataTable/SortQuery";
import Countries from "../Users/Countries.js";
import { useHomesStore } from "@/store/homes";

defineOptions({ name: "UserHomesTable" });

const props = defineProps({
  userId: {
    type: String,
    default: "",
    required: true,
  },
});

defineEmits(["delete", "edit"]);

const homesStore = useHomesStore();

const tableOptions = reactive({
  page: 1,
  itemsPerPage: 3,
  sortBy: [],
});
const error = ref({});
const users = ref([]);
const loading = ref(false);
const serverItemsLength = ref(0);
const usersItemsPerPageOptions = [
  { title: "1", value: 1 },
  { title: "2", value: 2 },
  { title: "3", value: 3 },
];

const countries = computed(() => Countries);
const headers = computed(() => UserHomesTableHeaders);

function getCountry(code) {
  return countries.value.filter((e) => e.value === code)[0].title;
}

function sortQuery(s) {
  return SortQuery(s);
}

async function loadItems() {
  if (props.userId === "") {
    return;
  }

  const { page, itemsPerPage, sortBy } = tableOptions;

  const sort = sortQuery(sortBy);

  try {
    const data = await homesStore.fetchHomes({
      userId: props.userId,
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

watch(
  () => props.userId,
  () => loadItems()
);

defineExpose({ loadItems });
</script>
