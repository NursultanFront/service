<template>
  <div class="ma-6 fill-height justify-start">
    <v-responsive class="align-start text-center fill-height">
      <v-card>
        <v-card-title class="d-flex align-center justify-space-between">
          <div>Users</div>
          <v-btn @click="openNewUser">Add User</v-btn>
        </v-card-title>
        <v-card-text>
          <users-table
            ref="usersTableRef"
            @delete="openDelete($event)"
            @edit="openEdit($event)"
          />
        </v-card-text>
      </v-card>
      <users-edit
        v-model="dialogs.edit.open"
        :edit="dialogs.edit.edit"
        :user="dialogs.edit.item"
        @close="dialogs.edit.open = false"
        @error="failure"
        @success="successEdit"
      />
      <ui-confirmation-dialog
        v-model="dialogs.confirmation.open"
        :title="dialogs.confirmation.title"
        :text="dialogs.confirmation.text"
        :button-text="dialogs.confirmation.buttonText"
        @confirm="sendDeleteUser"
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
    </v-responsive>
  </div>
</template>

<script setup>
import { reactive, ref, nextTick, defineAsyncComponent } from "vue";
import UsersTable from "./UsersTable.vue";
import { useUsersStore } from "@/store/users";

// Dialogs are hidden until the user actually opens one (Add/Edit/Delete),
// so there's no reason to ship and parse their code on the initial page
// load - load them on demand instead.
const UsersEdit = defineAsyncComponent(() => import("./UsersEdit.vue"));
const UiConfirmationDialog = defineAsyncComponent(() =>
  import("../UI/confirmation-dialog.vue")
);
const UiSuccessDialog = defineAsyncComponent(() =>
  import("../UI/success-dialog.vue")
);
const UiFailureDialog = defineAsyncComponent(() =>
  import("../UI/failure-dialog.vue")
);

const usersStore = useUsersStore();
const usersTableRef = ref(null);

const dialogs = reactive({
  confirmation: {
    open: false,
    title: "Delete User",
    text: "Do you want to delete this user?",
    buttonText: "Delete",
    item: {},
  },
  edit: {
    open: false,
    edit: false,
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

function successDelete(userName) {
  dialogs.success.title = "Delete Success";
  dialogs.success.subtitle = `User ${userName} deleted successfully`;
  dialogs.success.open = true;
  usersTableRef.value?.loadItems();
}

function successEdit() {
  const action = dialogs.edit.edit ? "edited" : "added";
  dialogs.success.title = "Success";
  dialogs.success.subtitle = `User ${action} successfully`;
  dialogs.success.open = true;
  dialogs.edit.open = false;
  dialogs.edit.edit = false;
  dialogs.edit.item = {};
  usersTableRef.value?.loadItems();
}

function failure(errors) {
  dialogs.failure.errors = errors;
  dialogs.failure.title = "Something went wrong";
  dialogs.failure.subtitle = "Creating user went wrong";
  dialogs.failure.open = true;
}

function openNewUser() {
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

async function openDelete(item) {
  dialogs.confirmation.item = item;
  await nextTick();
  dialogs.confirmation.open = true;
}

async function sendDeleteUser() {
  const { id, name } = dialogs.confirmation.item;
  try {
    await usersStore.deleteUser(id);
    successDelete(name);
    dialogs.confirmation.item = {};
  } catch (error) {
    const errors = [
      { message: "Deleting user went wrong" },
      { message: `Error Code: ${error.response?.status}` },
      { message: error.response?.data?.error },
    ];
    failure(errors);
  }
}
</script>
