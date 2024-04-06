export interface IMessage {
  guid: string;
  sender_id: number;
  conversation_id: number;
  message_kind: number;
  message: string;
  created_at: string;
  updated_at: string;
}

export const CreateConversationMessage = 0;
export const AddingUserMessage = 1;
export const UserTextMessage = 2;
