<script setup lang="ts">
import { ref, watch } from "vue";

const props = defineProps<{
  disabled: boolean;
}>();

const content = defineModel<string>({ required: true });

const emit = defineEmits<{
  send: [];
}>();

const input = ref<HTMLInputElement | null>(null);

function submit(): void {
  if (props.disabled) return;

  emit("send");
  input.value?.focus();
}

watch(
  () => props.disabled,
  (disabled) => {
    if (!disabled) {
      input.value?.focus();
    }
  },
  { flush: "post" },
);
</script>

<template>
  <form @submit.prevent="submit">
    <input
      ref="input"
      v-model="content"
      aria-label="消息内容"
      placeholder="输入消息，按回车发送"
      autocomplete="off"
      :disabled="props.disabled"
    >
    <button type="submit" :disabled="props.disabled">
      发送
    </button>
  </form>
</template>