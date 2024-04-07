import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

import "vue3-toastify/dist/index.css";
import VueChatScroll from "vue-chat-scroll";

const app = createApp(App);

app.use(store);
app.use(router);
app.use(VueChatScroll);
app.mount("#app");
