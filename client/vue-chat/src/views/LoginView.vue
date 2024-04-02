<template>
  <div class="v-login-container">
    <div class="v-login-header">
      <h1>Sign in/up</h1>
    </div>

    <div class="v-auth-form">
      <input class="v-login-input" type="text" placeholder="login" v-model="login" />
      <input class="v-login-input" type="password" placeholder="password" v-model="password" />

      <div class="v-login-wrapper">
        <div class="v-login-host-port">
          <input class="v-login-input" type="text" placeholder="host" v-model="host" />
          <input class="v-login-input" type="text" placeholder="port" v-model="port" />
        </div>
      </div>
    </div>

    <div class="v-login-wrapper">
      <button class="v-login-btn" @click="signIn">Sign in</button>
      <button class="v-login-btn" @click="signUp">Sign up</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { ResponseToast, NotifySystem } from "../lib/toasts/notifications";
import { IResponse, successResponse } from "../lib/reqresp-conn/reqresp";
import { Request } from "../lib/reqresp-conn/conn";
import { TSMap } from "typescript-map";
import { useStore } from "vuex";
import { validatePassword, validateLogin } from "../lib/utils/validation";
import {
  signUpEndpoint,
  signInEndpoint,
  ISignInResponse,
} from "../store/endpoints/endpoints";
import { User } from "../store/models/user";

const login = ref("");
const password = ref("");
const host = ref("localhost");
const port = ref(5050);

let validLogin: string | undefined;
let validPassword: string | undefined;

export default defineComponent({
  setup() {
    return {
      store: useStore(),
      login: login,
      password: password,
      host: host,
      port: port,
    };
  },
  methods: {
    signUp() {
      if (!this.validateUserData()) {
        return;
      }

      Request(host.value, port.value, this.store.state.requestTimeout, {
        method: "POST",
        url: signUpEndpoint,
        proto: "http",

        header: new TSMap([["Content-Type", "application/json"]]),
        body: JSON.stringify({
          login: login.value.trim(),
          password: password.value.trim(),
        }),
      })
        .then((value: IResponse) => {
          ResponseToast.notify(value.status_code, value.status);
          if (successResponse(value)) {
            this.signIn();
          }
        })
        .catch((error) => {
          NotifySystem.notify("error", error);
        });
    },
    signIn() {
      if (!this.validateUserData()) {
        return;
      }

      Request(host.value, port.value, this.store.state.requestTimeout, {
        method: "POST",
        url: signInEndpoint,
        proto: "http",

        header: new TSMap([["Content-Type", "application/json"]]),
        body: JSON.stringify({
          login: login.value.trimRight(),
          password: password.value.trim(),
        }),
      })
        .then((resp: IResponse) => {
          ResponseToast.notify(resp.status_code, resp.status);

          if (successResponse(resp)) {
            try {
              const value: ISignInResponse = JSON.parse(resp.body);
              this.store.state.user = new User(value.auth_token, value.user);
              NotifySystem.notify("success", `Welcom ${value.user.name}`);
            } catch (e) {
              NotifySystem.notify("error", "error parse json");
            }
          }
        })
        .catch((error) => {
          NotifySystem.notify("error", error);
        });
    },
    validateUserData(): boolean {
      validLogin = validateLogin(login.value.trim());
      if (validLogin !== undefined) {
        NotifySystem.notify("warning", validLogin);
        return false;
      }

      validPassword = validatePassword(password.value);
      if (validPassword !== undefined) {
        NotifySystem.notify("warning", validPassword);
        return false;
      }

      return true;
    },
  },
});
</script>

<style scoped>
.v-login-container {
  min-width: 440px;

  text-align: center;
  display: flex;
  flex-direction: column;
}

.v-auth-form {
  display: flex;
  flex-direction: column;
}

.v-login-btn {
  margin: 7px;
  padding: 7px 14px;
  width: 100px;

  border-radius: 30px;
  border: none;
  background-color: #375fff;

  color: #e6edf3;
  font-family: inherit;
  font-size: 16px;
}

.v-login-btn:hover {
  background-color: #001a83;
}

.v-login-logo {
  width: 200px;

  align-self: center;
}

.v-login-input {
  margin: 7px;
  padding: 7px 14px;
  background-color: #152032;

  color: #e6edf3;
  font-family: inherit;
  font-size: 14px;

  border-radius: 6px;
  border: 1px solid #375fff;
}

.v-login-input:hover {
  border: 1px solid #e6edf3;
}

.v-login-header {
  display: flex;
  flex-direction: column;
  justify-content: center;

  color: #e6edf3;

  font-weight: 100;
}
</style>
