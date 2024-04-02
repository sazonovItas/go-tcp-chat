import { IToken } from "./token";

export class User {
  public user: IUser | undefined;
  public token: IToken | undefined;

  constructor(token: IToken, user: IUser) {
    this.token = token;
    this.user = user;
  }
}

export interface IUser {
  id: number;
  name: string;
  login: string;
  color: string;
  password_hash: string;
}
