import { createStore, Store } from "vuex";
import { IUser } from "./models/user";
import { IToken } from "./models/token";
import { IMessage } from "./models/message";
import createPersistedState from "vuex-persistedstate";

export interface IState {
  requestTimeout: number;
  user: IUser | undefined;
  token: IToken | undefined;
  messages: Array<IMessage>;
  host: string;
  port: number;
}

const store: Store<IState> = createStore({
  plugins: [
    createPersistedState({
      getState: (key) => {
        const obj = localStorage.getItem(key);
        if (obj) {
          return JSON.parse(obj);
        }
        return null;
      },
      setState: (key, state) =>
        localStorage.setItem(key, JSON.stringify(state)),
    }),
  ],
  state() {
    return {
      requestTimeout: 2000,
      retryTimeout: 5000,
      user: undefined,
      token: undefined,
      messages: undefined,
      host: "",
      port: 0,
    };
  },
  getters: {},
  mutations: {
    setChatAppData(state, payload) {
      state.user = payload.user;
      state.token = payload.token;
      state.host = payload.host;
      state.port = payload.port;
      state.messages = new Array<IMessage>();
    },
    appendMessage(state, payload) {
      state.messages?.push(payload);
    },
    unshiftMessages(state, payload) {
      state.messages?.unshift(...payload);
    },
    updateMessages(state, payload) {
      state.messages = payload.messages;
    },
    logOut(state) {
      state.messages = new Array<IMessage>();
      state.token = undefined;
      state.user = undefined;
    },
  },
  actions: {},
  modules: {},
});

export default store;
