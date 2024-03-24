<template>
  <div>
    <button @click="sendRequest"></button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { Request } from "@/lib/reqrespConn/conn";
import { IRequest, ProtoHTTP } from "@/lib/reqrespConn/reqresp";
import { TSMap } from "typescript-map";

const req: IRequest = {
  method: "GET",
  url: "/user/12",
  proto: ProtoHTTP,

  header: new TSMap([["Content-Type", "application/json"]]),
  body: "hello",
};

export default defineComponent({
  methods: {
    sendRequest() {
      Request("127.0.0.1", 5050, 5000, req)
        .then((resp) => {
          console.log(resp);
        })
        .catch((e) => {
          console.log(e);
        });
    },
  },
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

nav {
  padding: 30px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
}

nav a.router-link-exact-active {
  color: #42b983;
}
</style>
