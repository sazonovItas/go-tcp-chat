const passwordRegExp = new RegExp("^[A-Za-z][A-Za-z0-9_]{4,20}$");
const loginRegExp = new RegExp("^[A-Za-z][A-Za-z0-9_]{4,20}$");

export function validatePassword(password: string): string | undefined {
  if (passwordRegExp.test(password)) {
    return undefined;
  }

  return "password should be at least 5 characters and digits and not more then 20";
}

export function validateLogin(login: string): string | undefined {
  if (loginRegExp.test(login)) {
    return undefined;
  }

  return "login should have at least 5 characters and digits and not more then 20";
}
