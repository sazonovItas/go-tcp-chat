import WSSocket from "@/lib/socket/wssocket";
import { createStore, Store } from "vuex";
import { User } from "./models/user";

export interface IState {
  requestTimeout: number;
  user: User | undefined;
  conn: WSSocket | undefined;
}

const store: Store<IState> = createStore({
  state() {
    return {
      requestTimeout: 2000,
      user: undefined,
      conn: undefined,
    };
  },
  getters: {},
  mutations: {},
  actions: {},
  modules: {},
});

export default store;
