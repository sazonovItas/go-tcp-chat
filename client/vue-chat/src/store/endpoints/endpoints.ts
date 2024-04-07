import { IMessage } from "../models/message";
import { IToken } from "../models/token";
import { IUser } from "../models/user";

export const signUpEndpoint = "/api/v1/signup";
// response_body: {}

export const signInEndpoint = "/api/v1/signin";
export interface ISignInResponse {
  auth_token: IToken;
  user: IUser;
}

export const signInByTokenEndpoint = "/api/v1/signin/token";
export const chattingEndpoint = "/api/v1/chatting";

export const messagesEndpoint = "/api/v1/messages";
export interface IMessagesRequest {
  auth_token: IToken;
  timestamp: string;
  limit: number;
}
export interface IMessagesResponse {
  messages: Array<IMessage>;
}

export const memberEndpoint = "/api/v1/member";
