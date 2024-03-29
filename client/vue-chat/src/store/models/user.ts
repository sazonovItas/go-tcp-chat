export class User {
  public login: string | undefined;
  public name: string | undefined;

  private id: number | undefined;
  private sessionToken: string | undefined;

  public User(sessionToken: string, id: number, login: string, name: string) {
    this.sessionToken = sessionToken;

    this.name = name;
  }

  public stringify(): string {
    return JSON.stringify({
      sessionToken: this.sessionToken,

      name: this.name,
    });
  }
}

export interface IUser {
  sessionToken: string;

  id: number;
  login: string;
  name: string;
  password: string;
}
