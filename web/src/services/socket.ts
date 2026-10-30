import type {ChatMessage, ChatRequest, ConnectionState,} from "../types/message";

interface SocketHandlers {
  onState: (state: ConnectionState) => void;
  onMessage: (message: ChatMessage) => void;
  onError: (message: string) => void;
}

function isChatMessage(value: unknown): value is ChatMessage {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  return (
    "type" in value && value.type === "chat" &&
    "username" in value && typeof value.username === "string" &&
    "content" in value && typeof value.content === "string" &&
    "sent_at" in value && typeof value.sent_at === "string" && !Number.isNaN(Date.parse(value.sent_at))
  );
}

// Chat 使用的 WebSocket
export class ChatSocket {
  private socket: WebSocket | null = null;
  private events: AbortController | null = null;
  private readonly handlers: SocketHandlers;

  // 构造器，传入提前写好的 SocketHandler 控制状态改变、发送消息、抛出错误的行为。
  constructor(handlers: SocketHandlers) {
    this.handlers = handlers;
  }

  // WebSocket 类的 Connect 功能
  connect(username: string): void {
    if (this.socket !== null) {
      throw new Error("已有连接，请等待当前连接关闭");
    }

    const url = new URL("/ws", window.location.href);
    url.protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    url.searchParams.set("username", username);

    const connection = new WebSocket(url);
    const events = new AbortController();
    const options = { signal: events.signal };

    this.socket = connection;
    this.events = events;
    this.handlers.onState("connecting");

    connection.addEventListener("open", () => {
        this.handlers.onState("connected");
    }, options);

    connection.addEventListener("message", (event) => {
        if (typeof event.data !== "string") return;

        let value: unknown;

        try {
          value = JSON.parse(event.data);
        } catch {
          this.handlers.onError("消息解析失败");
          return;
        }

        if (!isChatMessage(value)) {
          this.handlers.onError("收到格式不正确的消息");
          return;
        }

        // 发送消息
        this.handlers.onMessage(value);
    }, options);

    connection.addEventListener("error", () => {
        this.handlers.onError("连接发生错误，请检查 Go 服务是否启动");
    }, options);

    connection.addEventListener("close", () => {
        events.abort();
        this.socket = null;
        this.events = null;
        this.handlers.onState("disconnected");
    }, options);
  }

  send(content: string): void {
    if (this.socket === null || this.socket.readyState !== WebSocket.OPEN) {
      throw new Error("请先加入聊天");
    }

    const text = content.trim();

    if (!text) {
      throw new Error("消息不能为空");
    }

    const request: ChatRequest = {
      type: "chat",
      content: text,
    };

    const payload = JSON.stringify(request);

    if (new TextEncoder().encode(payload).length > 4096) {
      throw new Error("消息不能超过 4KB");
    }

    this.socket.send(payload);
  }

  disconnect(): void {
    this.events?.abort();
    this.socket?.close();

    this.events = null;
    this.socket = null;

    this.handlers.onState("disconnected");
  }
}