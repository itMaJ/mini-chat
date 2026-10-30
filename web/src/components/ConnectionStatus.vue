<script setup lang="ts">
import type { ConnectionState } from "../types/message";

const props = defineProps<{
  state: ConnectionState;
}>();

const labels: Record<ConnectionState, string> = {
  idle: "请填写昵称并加入聊天",
  connecting: "正在连接……",
  connected: "已连接",
  disconnected: "连接已断开，可以重新加入",
};
</script>

<template>
  <p role="status">{{ labels[props.state] }}</p>
</template>