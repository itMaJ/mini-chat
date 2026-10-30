<script setup lang="ts">
import { computed, ref } from "vue";
import JoinForm from "./components/JoinForm.vue";
import ConnectionStatus from "./components/ConnectionStatus.vue";
import MessageList from "./components/MessageList.vue";
import MessageComposer from "./components/MessageComposer.vue";
import { useChat } from "./composables/useChat";

const { state, messages, error, join, send } = useChat();

const content = ref("");

const joining = computed(
  () => state.value === "connecting" || state.value === "connected",
);

const connected = computed(
  () => state.value === "connected",
);

function submitMessage(): void {
  if (send(content.value)) {
    content.value = "";
  }
}
</script>

<template>
  <main>
    <h1>Mini Chat</h1>

    <ConnectionStatus :state="state" />

    <JoinForm
      :disabled="joining"
      @join="join"
    />

    <p v-if="error" class="error" role="alert">
      {{ error }}
    </p>

    <MessageList :messages="messages" />

    <MessageComposer
      v-model="content"
      :disabled="!connected"
      @send="submitMessage"
    />
  </main>
</template>