export interface IUser {
  id: number;
  name: string;
  login: string;
  color: string;
  password_hash: string;
}

export interface IPublicUser {
  id: number;
  name: string;
  login: string;
  color: string;
}
