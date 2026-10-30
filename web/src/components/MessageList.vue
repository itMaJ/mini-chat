<script setup lang="ts">
import { ref, watch } from "vue";
import type { ChatListItem } from "../types/message";

const props = defineProps<{
  messages: ChatListItem[];
}>();

const list = ref<HTMLUListElement | null>(null);

function formatTime(value: string): string {
  return new Date(value).toLocaleTimeString();
}

// Vue 更新消息列表后，再滚动到底部。
watch(
  () => props.messages.length,
  () => {
    if (list.value) {
      list.value.scrollTop = list.value.scrollHeight;
    }
  },
  { flush: "post" },
);
</script>

<template>
  <ul
    ref="list"
    class="messages"
    aria-label="聊天消息"
    aria-live="polite"
  >
    <li v-for="message in props.messages" :key="message.localId">
      <div class="message-meta">
        {{ message.username }} · {{ formatTime(message.sent_at) }}
      </div>
      <div class="message-content">{{ message.content }}</div>
    </li>
  </ul>
</template>