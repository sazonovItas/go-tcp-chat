<template>
  <div>
    <button v-on:click="hiServer" v-bind:disabled="disable_req"></button>
    <h1>{{ count }}</h1>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { IRequest } from "../model/modules/reqrespConn/reqresp";
import { TSMap } from "typescript-map";
import Request from "../model/modules/reqrespConn/conn";

const req: IRequest = {
  method: "GET",
  url: "/user/1020",
  proto: "http",

  header: new TSMap<string, string | number>([
    ["Content-Type", "application/json"],
  ]),
  body: "Request sender",
};

export default defineComponent({
  data() {
    return {
      user: undefined,
      disable_req: false,
      count: 0,
    };
  },
  methods: {
    hiServer() {
      this.disable_req = true;
      this.count++;

      Request("127.0.0.1", 5050, 5000, req).then((value) => {
        this.disable_req = false;
        console.log(value, value.body.length);
      });
    },
  },
});
</script>

<style>
html {
  background-color: #1e1817;
}
</style>
