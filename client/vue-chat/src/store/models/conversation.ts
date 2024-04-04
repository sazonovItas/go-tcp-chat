export interface IConversation {
  id: number;
  title: string;
  color: string;
  creator_id: number;
  conversation_kind: number;
}

export const ConversationP2PKind = 0;
export const ConversationGroupKind = 1;
