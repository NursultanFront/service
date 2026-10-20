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
      <div class="d-flex flex-column justify-center">
        <div class="text-body-1 text--light pt-7 pb-4">
          <v-form v-model="valid" ref="formRef">
            <v-row no-gutters>
              <v-col cols="12" class="px-2">
                <v-text-field
                  v-model="form.type"
                  variant="outlined"
                  label="Type"
                  :rules="[requiredRule]"
                />
              </v-col>
              <v-col cols="12" class="px-2">
                <v-text-field
                  v-model="form.address.address1"
                  variant="outlined"
                  label="Address Line 1"
                  :rules="[requiredRule]"
                />
              </v-col>
              <v-col md="6" sm="12" class="px-2">
                <v-text-field
                  v-model="form.address.address2"
                  variant="outlined"
                  label="Address Line 2"
                />
              </v-col>
              <v-col md="6" sm="12" class="px-2">
                <v-text-field
                  v-model="form.address.zipCode"
                  variant="outlined"
                  label="ZIP Code"
                  :rules="[requiredRule]"
                />
              </v-col>
              <v-col cols="12" class="px-2">
                <v-autocomplete
                  v-model="form.address.country"
                  class="ui-text-field"
                  variant="outlined"
                  :items="countries"
                  item-value="value"
                  item-text="title"
                  clearable
                  :rules="[requiredRule]"
                  label="Country"
                />
              </v-col>
              <v-col md="6" sm="12" class="px-2">
                <v-text-field
                  v-model="form.address.state"
                  variant="outlined"
                  label="State"
                  :rules="[requiredRule]"
                />
              </v-col>
              <v-col md="6" sm="12" class="px-2">
                <v-text-field
                  v-model="form.address.city"
                  variant="outlined"
                  label="City"
                  :rules="[requiredRule]"
                />
              </v-col>
            </v-row>
          </v-form>
        </div>
      </div>
    </template>
    <template #actions>
      <div class="d-flex flex-grow-1 justify-end">
        <v-btn class="align-self-right" text @click="closeDialog">
          Cancel
        </v-btn>

        <v-btn class="align-self-right ml-4" color="primary" @click="editHome">
          {{ dialogButtonText }}
        </v-btn>
      </div>
    </template>
  </ui-dialog>
</template>
<script setup>
import { reactive, ref, computed, onBeforeMount, watch, useAttrs } from "vue";
import Countries from "../Users/Countries.js";
import UiDialog from "../UI/dialog.vue";
import { useHomesStore } from "@/store/homes";

defineOptions({ name: "HomeEdit" });

const props = defineProps({
  userId: {
    type: String,
    default: "",
  },
  home: {
    type: Object,
    default: () => {},
  },
  edit: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["close", "success", "error"]);

const homesStore = useHomesStore();
const attrs = useAttrs();

const valid = ref(false);
const formRef = ref(null);

const defaultForm = () => ({
  type: "",
  address: {
    address1: "",
    address2: "",
    zipCode: "",
    city: "",
    state: "",
    country: "",
  },
});

const form = reactive(defaultForm());

onBeforeMount(() => {
  if (props.edit) {
    Object.assign(form, props.home);
  }
});

watch(
  () => props.home,
  () => {
    if (Object.keys(props.home).length) {
      Object.assign(form, props.home);
    }
  },
  { deep: true }
);

const dialogTitle = computed(() => (props.edit ? "Edit Home" : "Add Home"));
const dialogButtonText = computed(() => (props.edit ? "Edit" : "Add"));
const dialogProps = computed(() => ({
  scrollable: true,
  ...attrs,
}));
const countries = computed(() => Countries);

function requiredRule(v) {
  return !!v || "This field is required";
}

async function editHome() {
  if (formRef.value) {
    await formRef.value.validate();
  }

  if (!valid.value) {
    return;
  }

  const homeId = form.id;

  if (props.edit) {
    delete form.id;
    delete form.userID;
    delete form.dateCreated;
    delete form.dateUpdated;
  } else {
    form.userID = props.userId;
  }

  try {
    if (props.edit) {
      await homesStore.updateHome(homeId, form);
    } else {
      await homesStore.createHome(form);
    }

    emit("success");
    Object.assign(form, defaultForm());
  } catch (error) {
    const homePostData = error.response?.data || {};
    const errors = [
      { message: "Creating home went wrong" },
      { message: `Error Code: ${error.response?.status}` },
      { message: homePostData.error },
    ];
    if (homePostData.fields) {
      for (let i = 0; i < Object.values(homePostData.fields).length; i++) {
        errors.push({
          message: Object.values(homePostData.fields)[i],
        });
      }
    }
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
