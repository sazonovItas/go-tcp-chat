import { IToken } from "../models/token";
import { IUser } from "../models/user";

export const signUpEndpoint = "/api/v1/signup";
// response_body: {}

export const signInEndpoint = "/api/v1/signin";
export interface ISignInResponse {
  auth_token: IToken;
  user: IUser;
}
// response_body : {
//    auth_token: {
//      id: string,
//      user_id: number,
//    },
//    user: {
//      id: number,
//      login: string,
//      name: string,
//      color: string,
//      password_hash: string,
//    }
// }

export const signInByTokenEndpoint = "/api/v1/signin/token";
export const chattingEndpoint = "/api/v1/chatting";
