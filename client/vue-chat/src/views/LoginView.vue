<template>
  <div class="v-login-container">
    <div class="v-login-header">
      <h1>Sign in/up</h1>
    </div>

    <div class="v-auth-form">
      <input
        class="v-login-input"
        type="text"
        placeholder="login"
        v-model="login"
      />
      <input
        class="v-login-input"
        type="password"
        placeholder="password"
        v-model="password"
      />

      <div class="v-login-wrapper">
        <div class="v-login-host-port">
          <input
            class="v-login-input"
            type="text"
            placeholder="host"
            v-model="host"
          />
          <input
            class="v-login-input"
            type="text"
            placeholder="port"
            v-model="port"
          />
        </div>
      </div>
    </div>

    <div class="v-login-wrapper">
      <button class="v-login-btn">Sign in</button>
      <button class="v-login-btn">Sign up</button>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from "vue";
import { toast } from "vue3-toastify";

const login = ref("");
const password = ref("");
const host = ref("localhost");
const port = ref(5050);

export default defineComponent({
  setup() {
    return {
      login,
      password,
      host,
      port,
    };
  },
  methods: {
    testToast() {
      toast.success(login.value, { theme: "dark" });
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
