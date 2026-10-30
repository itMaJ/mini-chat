export interface ChatRequest {
  type: "chat";
  content: string;
}

export interface ChatMessage {
  type: "chat";
  username: string;
  content: string;
  sent_at: string;
}

// 本地列表使用的编号，方便 Vue 跟踪每条消息。
export interface ChatListItem extends ChatMessage {
  localId: number;
}

export type ConnectionState =
  | "idle"
  | "connecting"
  | "connected"
  | "disconnected";