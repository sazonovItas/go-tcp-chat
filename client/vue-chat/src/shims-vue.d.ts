/* eslint-disable */
declare module "*.vue" {
  import type { DefineComponent, ComponentCustomProperties } from "vue";
  import { Store } from "vuex";
  const component: DefineComponent<{}, {}, any>;
  export default component;

  interface State {
    requestTimeout: number;
    user: User | undefined;
    conn: WSSocket | undefined;
  }

  interface ComponentCustomProperties {
    store: Store<State>;
  }
}
