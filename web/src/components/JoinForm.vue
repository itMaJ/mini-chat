<script setup lang="ts">
import { ref } from "vue";

const props = defineProps<{
  disabled: boolean;
}>();

const emit = defineEmits<{
  join: [username: string];
}>();

const username = ref("");

function submit(): void {
  if (props.disabled) return;
  emit("join", username.value);
}
</script>

<template>
  <form @submit.prevent="submit">
    <input
      v-model="username"
      aria-label="昵称"
      placeholder="输入昵称"
      autocomplete="off"
      :disabled="props.disabled"
      required
    >
    <button type="submit" :disabled="props.disabled">
      加入聊天
    </button>
  </form>
</template>