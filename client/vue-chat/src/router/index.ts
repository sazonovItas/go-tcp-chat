import { createRouter, createWebHashHistory, RouteRecordRaw } from "vue-router";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "login",
    component: () => import("@/views/LoginView.vue"),
  },
  {
    path: "/chat",
    name: "chat",
    component: () => import("@/views/ChatView.vue"),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
