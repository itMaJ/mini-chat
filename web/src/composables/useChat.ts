
import { onUnmounted, ref } from "vue";
import { ChatSocket } from "../services/socket";
import type { ChatListItem, ConnectionState } from "../types/message";

export function useChat() {
    const state = ref<ConnectionState>("idle");
    const messages = ref<ChatListItem[]>([]);
    const error = ref("")

    let nextId = 1;

    const client = new ChatSocket({
        onState(nextState) {
            state.value = nextState
        },
        onMessage(message) {
            messages.value.push({
                ...message,
                localId: nextId++
            })
        },
        onError(message) {
            error.value = message
        },
    });

    // 对 catch 到的错误，封装为错误字符串
    function reportError(reason: unknown): void {
        error.value = reason instanceof Error ? reason.message : "操作失败";
    }

    function join(input: string): void {
        const username = input.trim()
        const length = Array.from(username).length

        if (length < 1 || length > 20) {
            error.value = "昵称需要包含 1～20 个字符";
            return;
        }

        error.value = ""
        try {
            client.connect(username)
        } catch(reason) {
            reportError(reason)
        }
    }

    // 发送消息，返回是否发送成功
    function send(message: string): boolean {
        error.value = ""

        try {
            client.send(message)
        } catch(reason) {
            reportError(reason)
            return false
        }
        return true
    }

    onUnmounted(() => {
        client.disconnect();
    });

    return {
        state,
        messages,
        error,
        join,
        send
    }
}
