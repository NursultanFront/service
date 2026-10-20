<template>
  <v-card
    show-title
    style="height: 100%"
    title-padding="pb-0 pt-4 pl-7 pr-4"
    class="pb-4"
  >
    <v-card-title>
      <v-row no-gutters>
        <v-col cols="6" class="justify-start align-center d-flex">
          Homes
        </v-col>
        <v-col cols="6" class="justify-end align-center d-flex">
          <v-btn outlined @click="openNewHome">Add Home</v-btn>
        </v-col>
      </v-row>
    </v-card-title>
    <v-card-text>
      <v-row no-gutters class="d-flex justify-space-around">
        <v-col cols="12" class="justify-start">
          <user-homes-table
            ref="homesTable"
            :user-id="userId"
            @delete="openDelete($event)"
            @edit="openEdit($event)"
          />
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
  <home-edit
    v-model="dialogs.edit.open"
    :user-id="userId"
    :edit="dialogs.edit.edit"
    :home="dialogs.edit.item"
    @close="dialogs.edit.open = false"
    @error="failure"
    @success="successEdit"
  />
  <ui-confirmation-dialog
    v-model="dialogs.confirmation.open"
    :title="dialogs.confirmation.title"
    :text="dialogs.confirmation.text"
    :button-text="dialogs.confirmation.buttonText"
    @confirm="sendDeleteHome"
    @input="dialogs.confirmation.open = false"
  />
  <ui-success-dialog
    v-model="dialogs.success.open"
    :title="dialogs.success.title"
    :subtitle="dialogs.success.subtitle"
    @close="dialogs.success.open = false"
  />
  <ui-failure-dialog
    v-model="dialogs.failure.open"
    :title="dialogs.failure.title"
    :subtitle="dialogs.failure.subtitle"
    :errors="dialogs.failure.errors"
    @close="dialogs.failure.open = false"
  />
</template>
<script setup>
import { reactive, ref, nextTick, defineAsyncComponent } from "vue";
import UserHomesTable from "../Users/UserHomesTable";
import { useHomesStore } from "@/store/homes";

// Dialogs are hidden until the user actually opens one (Add/Edit/Delete),
// so there's no reason to ship and parse their code on the initial page
// load - load them on demand instead.
const HomeEdit = defineAsyncComponent(() => import("../Users/HomeEdit"));
const UiConfirmationDialog = defineAsyncComponent(() =>
  import("../UI/confirmation-dialog.vue")
);
const UiSuccessDialog = defineAsyncComponent(() =>
  import("../UI/success-dialog.vue")
);
const UiFailureDialog = defineAsyncComponent(() =>
  import("../UI/failure-dialog.vue")
);

defineOptions({ name: "UserHomes" });

defineProps({
  userId: {
    type: String,
    default: "",
    required: true,
  },
});

const homesStore = useHomesStore();
const homesTable = ref(null);

const dialogs = reactive({
  edit: {
    open: false,
    edit: false,
    item: {},
  },
  confirmation: {
    open: false,
    title: "Delete Home",
    text: "Do you want to delete this home?",
    buttonText: "Delete",
    item: {},
  },
  success: {
    open: false,
    title: "",
    subtitle: "",
  },
  failure: {
    open: false,
    title: "",
    subtitle: "",
    errors: [],
  },
});

function openNewHome() {
  dialogs.edit.open = true;
  dialogs.edit.edit = false;
}

async function openEdit(item) {
  dialogs.edit.edit = true;
  await nextTick();
  dialogs.edit.item = item;
  await nextTick();
  dialogs.edit.open = true;
}

async function sendDeleteHome() {
  const { id } = dialogs.confirmation.item;
  try {
    await homesStore.deleteHome(id);
    successDelete();
    dialogs.confirmation.item = {};
  } catch (error) {
    const errors = [
      { message: "Deleting home went wrong" },
      { message: `Error Code: ${error.response?.status}` },
      { message: error.response?.data?.error },
    ];
    failure(errors);
  }
}

async function openDelete(item) {
  dialogs.confirmation.item = item;
  await nextTick();
  dialogs.confirmation.open = true;
}

function successDelete() {
  dialogs.success.title = "Delete Success";
  dialogs.success.subtitle = `Home deleted successfully`;
  dialogs.success.open = true;
  homesTable.value?.loadItems();
}

function successEdit() {
  const action = dialogs.edit.edit ? "edited" : "added";
  dialogs.success.title = "Success";
  dialogs.success.subtitle = `Home ${action} successfully`;
  dialogs.success.open = true;
  dialogs.edit.open = false;
  dialogs.edit.edit = false;
  dialogs.edit.item = {};
  homesTable.value?.loadItems();
}

function failure(errors) {
  dialogs.failure.errors = errors;
  dialogs.failure.title = "Something went wrong";
  dialogs.failure.subtitle = "Creating user went wrong";
  dialogs.failure.open = true;
}
</script>
