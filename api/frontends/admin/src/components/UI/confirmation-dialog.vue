<template>
  <ui-dialog v-bind="dialogProps" @click:outside="closeDialog">
    <template #body>
      <div class="d-flex flex-column justify-center px-3">
        <div v-if="showCaption" class="text--dark pt-3 text--caption">
          {{ dialogTitle }}
        </div>
        <div class="text-body-1 text--light pt-7 pb-4">
          {{ dialogText }}
        </div>
      </div>
    </template>
    <template #actions>
      <div class="d-flex flex-grow-1 justify-end">
        <v-btn class="align-self-right" text @click="closeDialog">
          Cancel
        </v-btn>

        <v-btn class="align-self-right ml-4" color="primary" @click="confirm">
          {{ dialogButtonText }}
        </v-btn>
      </div>
    </template>
  </ui-dialog>
</template>

<script setup>
import { computed, useAttrs } from "vue";
import UiDialog from "./dialog.vue";

defineOptions({ name: "UiConfirmationDialog" });

const props = defineProps({
  title: {
    type: String,
    default: () => null,
  },
  text: {
    type: String,
    default: () => null,
  },
  buttonText: {
    type: String,
    default: () => null,
  },
  showCaption: {
    type: Boolean,
    default: () => true,
  },
});

const emit = defineEmits(["input", "confirm"]);

const attrs = useAttrs();

const dialogTitle = computed(() => props.title ?? "Confirm");
const dialogText = computed(
  () => props.text ?? "Do you want to confirm this action?"
);
const dialogButtonText = computed(() => props.buttonText ?? "Confirm");
const dialogProps = computed(() => ({
  noTitle: true,
  contentClass: "ui-dialog--confirmation",
  ...attrs,
}));

function closeDialog() {
  emit("input", false);
}

function confirm() {
  emit("confirm");
  closeDialog();
}
</script>
<style lang="scss" scoped>
.text--caption {
  font-size: 20px;
  font-weight: 500;
}
</style>
