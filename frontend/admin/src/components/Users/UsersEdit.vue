<template>
  <ui-dialog
    v-bind="dialogProps"
    @click:outside="closeDialog"
    @close="closeDialog"
  >
    <template #title>
      <div>{{ dialogTitle }}</div>
    </template>
    <template #body>
      <div class="d-flex flex-column justify-center px-3">
        <div class="text-body-1 text--light pt-7 pb-4">
          <v-form v-model="valid" ref="formRef">
            <v-text-field
              v-model="form.name"
              variant="outlined"
              label="Name"
              :rules="[requiredRule]"
            />
            <v-text-field
              v-model="form.email"
              variant="outlined"
              label="Email"
              :rules="[emailRule]"
            />
            <v-select
              v-model="form.roles"
              label="Role"
              multiple
              variant="outlined"
              :items="userRoles"
              disabled
            />
            <v-text-field
              v-model="form.department"
              variant="outlined"
              label="Department"
            />
            <v-text-field
              v-if="!edit"
              v-model="form.password"
              :append-inner-icon="visible ? 'fas fa-eye' : 'fas fa-eye-slash'"
              :type="visible ? 'text' : 'password'"
              placeholder="Enter your password"
              variant="outlined"
              :rules="[requiredRule]"
              @click:append-inner="visible = !visible"
            />
            <v-text-field
              v-if="!edit"
              v-model="form.passwordConfirm"
              :append-inner-icon="
                visibleConfirm ? 'fas fa-eye' : 'fas fa-eye-slash'
              "
              :type="visibleConfirm ? 'text' : 'password'"
              variant="outlined"
              label="Confirm Password"
              :rules="[passwordRule]"
              @click:append-inner="visibleConfirm = !visibleConfirm"
            />
          </v-form>
        </div>
      </div>
    </template>
    <template #actions>
      <div class="d-flex flex-grow-1 justify-end">
        <v-btn class="align-self-right" text @click="closeDialog">
          Cancel
        </v-btn>

        <v-btn class="align-self-right ml-4" color="primary" @click="editUser">
          {{ dialogButtonText }}
        </v-btn>
      </div>
    </template>
  </ui-dialog>
</template>
<script setup>
import { reactive, ref, computed, onBeforeMount, watch, useAttrs } from "vue";
import UiDialog from "../UI/dialog.vue";
import { useUsersStore } from "@/store/users";

defineOptions({ name: "UsersEdit" });

const props = defineProps({
  user: {
    type: Object,
    default: () => {},
  },
  edit: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close", "success", "error"]);

const usersStore = useUsersStore();
const attrs = useAttrs();

const visible = ref(false);
const visibleConfirm = ref(false);
const valid = ref(false);
const formRef = ref(null);

const defaultForm = () => ({
  name: "",
  email: "",
  roles: [],
  department: "",
  password: "",
  passwordConfirm: "",
});

const form = reactive(defaultForm());

onBeforeMount(() => {
  if (props.edit && Object.keys(props.user).length) {
    Object.assign(form, props.user);
  }
});

watch(
  () => props.user,
  () => {
    if (Object.keys(props.user).length) {
      Object.assign(form, props.user);
    }
  },
  { deep: true }
);

const dialogTitle = computed(() => (props.edit ? "Edit User" : "Add User"));
const dialogButtonText = computed(() => (props.edit ? "Edit" : "Add"));
const dialogProps = computed(() => ({
  scrollable: true,
  ...attrs,
}));
const userRoles = computed(() => [
  { title: "Admin", value: "ADMIN" },
  { title: "User", value: "USER" },
]);

function passwordRule(v) {
  if (!v) {
    return false;
  }
  return form.password === form.passwordConfirm || "Passwords don't match";
}

function requiredRule(v) {
  return !!v || "This field is required";
}

function emailRule(v) {
  if (!v) {
    return false;
  }
  if (v.length <= 6 || v.length >= 128) {
    return false;
  }
  const emailRegExp =
    /^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
  return emailRegExp.test(v) || "Invalid Email Format";
}

async function editUser() {
  if (formRef.value) {
    await formRef.value.validate();
  }

  if (!valid.value) {
    return;
  }

  const userId = form.id;

  if (props.edit) {
    delete form.id;
    delete form.dateCreated;
    delete form.dateUpdated;
    delete form.enabled;
  }

  try {
    if (props.edit) {
      await usersStore.updateUser(userId, form);
    } else {
      await usersStore.createUser(form);
    }

    emit("success");
    Object.assign(form, defaultForm());
  } catch (error) {
    const errors = [
      { message: "Creating user went wrong" },
      { message: `Error Code: ${error.response?.status}` },
      { message: error.response?.data?.error },
    ];
    emit("error", errors);
  }
}

function closeDialog() {
  emit("close");
}
</script>
<style lang="scss" scoped>
.text--caption {
  font-size: 20px;
  font-weight: 500;
}
</style>
