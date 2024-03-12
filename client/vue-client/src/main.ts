import { createApp } from "vue";
import App from "./view/App.vue";
import WSSocket from "./model/modules/socket/wssocket";
import { IRequest, IResponse } from "./model/modules/socket/reqresp";
import { TSMap } from "typescript-map";

const req: IRequest = {
  method: "GET",
  url: "/user/1020",
  proto: "http",

  header: new TSMap<string, string>([["Content-Type", "text/json"]]),
  body: "hello",
};

const wssocket = new WSSocket("127.0.0.1", 5050, undefined, (data: Buffer) => {
  try {
    const resp: IResponse = JSON.parse(data.toString());
    console.log(resp);
  } catch (e) {
    console.log(e);
  }
});
wssocket.setTimeout(600, () => {
  wssocket.connectSocket();
  try {
    wssocket.socketSend(JSON.stringify(req), 0x1, true);
  } catch (e) {
    console.error(e);
  }
});

setTimeout(() => {
  wssocket.closeSocket();
}, 10000);

createApp(App).mount("#app");
