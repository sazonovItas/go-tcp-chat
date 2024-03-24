export class User {
  public name: string | undefined;

  private sessionToken: string | undefined;

  public User(sessionToken: string, name: string) {
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
  name: string;
  password: string;
}
